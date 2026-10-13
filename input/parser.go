package input

import (
	"bufio"
	"conways/life"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func ReadGrid(path string) (life.Grid, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open grid file: %w", err)
	}
	defer file.Close()

	return ParseGrid(file, path)
}

// ParseGrid reads and validates a grid from r. The source name is included in
// errors that refer to the grid itself, which keeps file and uploaded input
// errors useful to callers.
func ParseGrid(r io.Reader, sourceNames ...string) (life.Grid, error) {
	source := "grid"
	if len(sourceNames) > 0 {
		source = sourceNames[0]
	}
	scanner := bufio.NewScanner(r)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("read grid header: %w", err)
		}
		return nil, fmt.Errorf("grid file is empty")
	}

	header := strings.Fields(scanner.Text())
	if len(header) != 2 {
		return nil, fmt.Errorf("grid header must contain width and height")
	}

	width, err := strconv.Atoi(header[0])
	if err != nil || width < 1 || width > 10000 {
		return nil, fmt.Errorf("invalid width %q: must be between 1 and 10000", header[0])
	}

	height, err := strconv.Atoi(header[1])
	if err != nil || height < 1 || height > 10000 {
		return nil, fmt.Errorf("invalid height %q: must be between 1 and 10000", header[1])
	}

	grid, err := life.NewGrid(height, width)
	if err != nil {
		return nil, fmt.Errorf("new grid from %q: %w", source, err)
	}

	for row := range height {
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return nil, fmt.Errorf("read row %d: %w", row+1, err)
			}
			return nil, fmt.Errorf("missing row %d", row+1)
		}

		line := strings.TrimSuffix(scanner.Text(), "\r")
		if len(line) != width {
			return nil, fmt.Errorf("row %d has width %d; expected %d", row+1, len(line), width)
		}

		for column := range width {
			switch line[column] {
			case '#':
				grid[row][column] = 1
			case '.':
				grid[row][column] = 0
			default:
				return nil, fmt.Errorf("invalid character %q at row %d, column %d", line[column], row+1, column+1)
			}
		}
	}

	if scanner.Scan() {
		return nil, fmt.Errorf("grid contains more than %d rows", height)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read grid: %w", err)
	}

	return grid, nil
}
