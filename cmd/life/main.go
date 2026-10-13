package main

import (
	"conways/input"
	"conways/life"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	defaultPath    = "testdata/t1_block.txt"
	maxGenerations = 20
	maxRepetitions = 3
)

func gridAsStr(g life.Grid) string {
	var sb strings.Builder
	for _, row := range g {
		for _, alive := range row {
			if alive == 1 {
				sb.WriteByte('#')
			} else {
				sb.WriteByte('.')
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	wrap := flag.Bool("wrap", false, "wrap grid edges toroidally")
	generations := flag.Int("generations", maxGenerations, "number of generations to run")
	flag.Parse()

	path := defaultPath
	if flag.NArg() > 0 {
		path = flag.Arg(0)
	}
	if *generations < 0 {
		fmt.Fprintln(os.Stderr, "generations must not be negative")
		os.Exit(1)
	}

	g, err := input.ReadGrid(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	seen := map[string]int{gridAsStr(g): 1}
	for i := 0; i < *generations; i++ {
		fmt.Print("\033[H\033[2J")
		fmt.Printf("Generation %d\n\n", i)
		fmt.Print(gridAsStr(g))

		g, err = life.Step(g, *wrap)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		state := gridAsStr(g)
		seen[state]++
		if seen[state] >= maxRepetitions {
			break
		}
		time.Sleep(300 * time.Millisecond)
	}
}
