package life

import (
	"errors"
)

var dir = [][]int{
	{-1, -1},
	{-1, 0},
	{-1, 1},
	{0, -1},
	{0, 1},
	{1, -1},
	{1, 0},
	{1, 1},
}

func neighbors(grid [][]int, r, c int, wrap bool) int {

	count := 0

	for _, d := range dir {
		row := r + d[0]
		column := c + d[1]

		if wrap {
			row = (row + len(grid)) % len(grid)
			column = (column + len(grid[0])) % len(grid[0])
		} else if row < 0 || row >= len(grid) || column < 0 || column >= len(grid[0]) {
			continue
		}

		if grid[row][column] == 1 {
			count++
		}
	}

	return count
}

var ErrBadInput = errors.New("bad input") // sentinel error

func Step(grid [][]int, wrap bool) ([][]int, error) {

	if len(grid) == 0 || len(grid[0]) == 0 {
		return grid, ErrBadInput
	}

	result, err := NewGrid(len(grid), len(grid[0]))
	if err != nil {
		return grid, err
	}

	for y, rows := range grid {
		for x, _ := range rows {
			n := neighbors(grid, y, x, wrap)
			if grid[y][x] == 1 {
				if n < 2 {
					result[y][x] = 0
				}
				if n >= 4 {
					result[y][x] = 0
				}
				if n == 2 || n == 3 {
					result[y][x] = 1
				}
			} else {
				if n == 3 {
					result[y][x] = 1
				}
			}
		}
	}

	return result, nil
}

type Grid [][]int

func NewGrid(h, w int) (Grid, error) {

	if h < 0 || w < 0 {
		return nil, ErrBadInput
	}

	g := make(Grid, h)
	for r := range g {
		g[r] = make([]int, w)
	}

	return g, nil
}
