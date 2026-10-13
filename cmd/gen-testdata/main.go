package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

var dir string

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func write(name, content string) {
	path := filepath.Join(dir, name)
	must(os.WriteFile(path, []byte(content), 0o644))
	fmt.Println("wrote", path)
}

// grid builds a valid input from rows of text (W H header derived from rows).
func grid(rows ...string) string {
	return fmt.Sprintf("%d %d\n%s\n", len(rows[0]), len(rows), strings.Join(rows, "\n"))
}

// fill builds a w x h grid where every cell has the same state.
func fill(w, h int, alive bool) string {
	ch := "."
	if alive {
		ch = "#"
	}
	rows := make([]string, h)
	for i := range rows {
		rows[i] = strings.Repeat(ch, w)
	}
	return grid(rows...)
}

// randomGrid streams a random grid to disk so huge sizes don't need huge strings.
func randomGrid(name string, w, h int, density float64, rng *rand.Rand) {
	path := filepath.Join(dir, name)
	f, err := os.Create(path)
	must(err)
	defer f.Close()
	bw := bufio.NewWriterSize(f, 1<<20)
	writeRandom(bw, w, h, density, rng)
	must(bw.Flush())
	fmt.Println("wrote", path)
}

func writeRandom(out io.Writer, w, h int, density float64, rng *rand.Rand) {
	fmt.Fprintf(out, "%d %d\n", w, h)
	row := make([]byte, w+1)
	row[w] = '\n'
	for r := 0; r < h; r++ {
		for c := 0; c < w; c++ {
			if rng.Float64() < density {
				row[c] = '#'
			} else {
				row[c] = '.'
			}
		}
		out.Write(row)
	}
}

func main() {
	flag.StringVar(&dir, "dir", "testdata", "output directory")
	big := flag.Bool("big", false, "also generate 10000x10000 (about 100 MB)")
	seed := flag.Int64("seed", 42, "random seed (same seed, same files)")
	flag.Parse()
	must(os.MkdirAll(dir, 0o755))

	// ---- the five acceptance tests (T1 to T5) with exact expected output
	block := grid("....", ".##.", ".##.", "....")
	write("t1_block.txt", block)
	write("t1_block.gen1.want", block)
	write("t1_block.gen7.want", block)

	write("t2_blinker.txt", grid(".....", "..#..", "..#..", "..#..", "....."))
	write("t2_blinker.gen1.want", grid(".....", ".....", ".###.", ".....", "....."))
	write("t2_blinker.gen2.want", grid(".....", "..#..", "..#..", "..#..", "....."))

	write("t3_glider.txt", grid(".#....", "..#...", "###...", "......", "......", "......"))
	write("t3_glider.gen4.want", grid("......", "..#...", "...#..", ".###..", "......", "......"))

	// T4 and T5 share one input: a blinker split across the left/right edge.
	edge := grid(".....", ".....", "##..#", ".....", ".....")
	write("t4_wrap.txt", edge)                                                    // run with --wrap
	write("t4_wrap.gen1.want", grid(".....", "#....", "#....", "#....", ".....")) // wrapped result
	write("t5_edge_dead.txt", edge)                                               // run without --wrap
	write("t5_edge_dead.gen1.want", fill(5, 5, false))                            // everything dies

	// ---- edge cases
	write("edge_1x1_alive.txt", "1 1\n#\n")
	write("edge_1x1_alive.gen1.want", "1 1\n.\n") // lone cell dies
	write("edge_all_dead.txt", fill(6, 6, false))
	write("edge_all_alive.txt", fill(6, 6, true))
	write("edge_1xN_row.txt", grid("#####"))
	write("edge_Nx1_col.txt", "1 5\n#\n#\n#\n#\n#\n")

	// ---- malformed inputs: your CLI should exit 1 with a clear message
	write("bad_row_width.txt", "3 2\n..#\n.#\n")     // row shorter than header
	write("bad_char.txt", "3 2\n..#\n.x.\n")         // character not # or .
	write("bad_header.txt", "three two\n..#\n")      // header not numbers
	write("bad_zero_size.txt", "0 0\n")              // below the 1..10000 range
	write("bad_missing_rows.txt", "3 3\n...\n...\n") // fewer rows than declared
	write("bad_too_big.txt", "10001 5\n.....\n")     // above the range
	write("bad_empty.txt", "")                       // empty file
	write("bad_negative.txt", "-3 2\n...\n...\n")    // negative size

	// ---- performance inputs (spec: random 30% density)
	rng := rand.New(rand.NewSource(*seed))
	randomGrid("perf_100x100.txt", 100, 100, 0.30, rng)
	randomGrid("perf_1000x1000.txt", 1000, 1000, 0.30, rng)
	if *big {
		randomGrid("perf_10000x10000.txt", 10000, 10000, 0.30, rng)
	}
}
