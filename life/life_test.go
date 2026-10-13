package life_test

import (
	"conways/input"
	"conways/life"
	"reflect"
	"testing"
)

func TestReadGridValidScenarios(t *testing.T) {
	tests := []struct {
		name   string
		file   string
		width  int
		height int
	}{
		{"edge 1x1 alive", "../testdata/edge_1x1_alive.txt", 1, 1},
		{"edge 1xN row", "../testdata/edge_1xN_row.txt", 5, 1},
		{"edge Nx1 column", "../testdata/edge_Nx1_col.txt", 1, 5},
		{"edge all alive", "../testdata/edge_all_alive.txt", 6, 6},
		{"edge all dead", "../testdata/edge_all_dead.txt", 6, 6},
		{"performance 100x100", "../testdata/perf_100x100.txt", 100, 100},
		{"performance 1000x1000", "../testdata/perf_1000x1000.txt", 1000, 1000},
		{"T1 block", "../testdata/t1_block.txt", 4, 4},
		{"T2 blinker", "../testdata/t2_blinker.txt", 5, 5},
		{"T3 glider", "../testdata/t3_glider.txt", 6, 6},
		{"T4 toroidal edge", "../testdata/t4_wrap.txt", 5, 5},
		{"T5 dead edge", "../testdata/t5_edge_dead.txt", 5, 5},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			grid, err := input.ReadGrid(tc.file)
			if err != nil {
				t.Fatalf("ReadGrid(%q) returned error: %v", tc.file, err)
			}
			if len(grid) != tc.height || len(grid[0]) != tc.width {
				t.Fatalf("ReadGrid(%q) dimensions = %dx%d; want %dx%d", tc.file, len(grid[0]), len(grid), tc.width, tc.height)
			}
		})
	}
}

func TestStepScenarios(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		want        string
		generations int
		wrap        bool
	}{
		{"edge 1x1 alive", "../testdata/edge_1x1_alive.txt", "../testdata/edge_1x1_alive.gen1.want", 1, false},
		{"T1 block generation 1", "../testdata/t1_block.txt", "../testdata/t1_block.gen1.want", 1, false},
		{"T1 block generation 7", "../testdata/t1_block.txt", "../testdata/t1_block.gen7.want", 7, false},
		{"T2 blinker generation 1", "../testdata/t2_blinker.txt", "../testdata/t2_blinker.gen1.want", 1, false},
		{"T2 blinker generation 2", "../testdata/t2_blinker.txt", "../testdata/t2_blinker.gen2.want", 2, false},
		{"T3 glider generation 4", "../testdata/t3_glider.txt", "../testdata/t3_glider.gen4.want", 4, false},
		{"T4 toroidal edge generation 1", "../testdata/t4_wrap.txt", "../testdata/t4_wrap.gen1.want", 1, true},
		{"T5 dead edge generation 1", "../testdata/t5_edge_dead.txt", "../testdata/t5_edge_dead.gen1.want", 1, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			grid, err := input.ReadGrid(tc.input)
			if err != nil {
				t.Fatalf("ReadGrid(%q) returned error: %v", tc.input, err)
			}
			want, err := input.ReadGrid(tc.want)
			if err != nil {
				t.Fatalf("ReadGrid(%q) returned error: %v", tc.want, err)
			}

			for generation := 0; generation < tc.generations; generation++ {
				grid, err = life.Step(grid, tc.wrap)
				if err != nil {
					t.Fatalf("Step(%q), generation %d, returned error: %v", tc.input, generation+1, err)
				}
			}
			if !reflect.DeepEqual(grid, want) {
				t.Fatalf("Step(%q) = %v; want %v", tc.input, grid, want)
			}
		})
	}
}

func TestPerformanceScenarios(t *testing.T) {
	tests := []struct {
		name string
		file string
	}{
		{"100x100", "../testdata/perf_100x100.txt"},
		{"1000x1000", "../testdata/perf_1000x1000.txt"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			grid, err := input.ReadGrid(tc.file)
			if err != nil {
				t.Fatalf("ReadGrid(%q) returned error: %v", tc.file, err)
			}
			if _, err := life.Step(grid, false); err != nil {
				t.Fatalf("Step(%q) returned error: %v", tc.file, err)
			}
		})
	}
}
