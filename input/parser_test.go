package input

import (
	"conways/life"
	"reflect"
	"strings"
	"testing"
)

func TestParseGridReader(t *testing.T) {
	want := life.Grid{{1, 0}, {0, 1}}
	got, err := ParseGrid(strings.NewReader("2 2\n#.\n.#\n"))
	if err != nil {
		t.Fatalf("ParseGrid returned error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ParseGrid = %v; want %v", got, want)
	}
}

func TestReadGridBadInputs(t *testing.T) {
	tests := []struct {
		name string
		file string
	}{
		{"empty file", "../testdata/bad_empty.txt"},
		{"invalid header", "../testdata/bad_header.txt"},
		{"zero size", "../testdata/bad_zero_size.txt"},
		{"negative size", "../testdata/bad_negative.txt"},
		{"size too large", "../testdata/bad_too_big.txt"},
		{"row width mismatch", "../testdata/bad_row_width.txt"},
		{"invalid character", "../testdata/bad_char.txt"},
		{"missing rows", "../testdata/bad_missing_rows.txt"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			grid, err := ReadGrid(tc.file)
			if err == nil {
				t.Fatalf("readGrid(%q) returned nil error and grid %v", tc.file, grid)
			}
		})
	}
}
