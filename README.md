# Conway's Game of Life

A Go implementation of Conway's Game of Life on a finite rectangular grid.
The simulator supports ordinary dead borders and optional toroidal wrapping.

## Rules

Each generation is computed from the previous generation, so updates happen
simultaneously.

| Current cell | Live neighbors | Next state |
| --- | ---: | --- |
| Alive | 0 or 1 | Dead |
| Alive | 2 or 3 | Alive |
| Alive | 4 or more | Dead |
| Dead | Exactly 3 | Alive |
| Dead | Otherwise | Dead |

By default, cells outside the grid are dead. With `--wrap`, the grid becomes
toroidal: the top and bottom rows and the left and right columns are adjacent.

## Running the simulator

Run the application with:

```sh
go run ./cmd/life [options] [file]
```

If no file is provided, the default input is `testdata/t1_block.txt`.

Options:

```text
--generations N   Number of generations to simulate (default: 20)
--wrap            Wrap edges toroidally
```

For example:

```sh
go run ./cmd/life --generations 4 testdata/t3_glider.txt
go run ./cmd/life --wrap --generations 1 testdata/t4_wrap.txt
```

The terminal view displays each generation as it runs. The simulator also
stops early when the same state has appeared three times.

Invalid input or invalid command-line values are written to standard error and
return a nonzero exit status.

## Input format

The first line contains the width and height. Each following line contains one
row. `#` represents a live cell and `.` represents a dead cell.

```text
5 5
.....
..#..
..#..
..#..
.....
```

The parser requires:

- Width and height between 1 and 10,000.
- Exactly the declared number of rows.
- Every row to match the declared width.
- Only `#` and `.` as cell characters.

## Project structure

```text
life/              Grid representation and simulation rules
input/             Grid-file parsing
cmd/life/          Command-line application
cmd/gen-testdata/  Test-data generator
testdata/          Example grids and expected generations
```

The `life` package contains the reusable simulation engine. Its main operation
is:

```go
next, err := life.Step(current, wrap)
```

The input grid is not modified. Each step allocates a new grid, so one
generation uses O(width × height) time and O(width × height) additional memory.

## Testing

Run the complete test suite and static checks with:

```sh
go test ./...
go vet ./...
```

The tests cover malformed input, edge-sized grids, dead-border behavior,
toroidal behavior, standard patterns, and large performance fixtures.

## Generating test data

The test-data generator can recreate the fixtures in `testdata/`:

```sh
go run ./cmd/gen-testdata
```

Additional options are available for the output directory, random seed, and a
10,000 × 10,000 fixture:

```sh
go run ./cmd/gen-testdata -dir /tmp/life-data -seed 7
go run ./cmd/gen-testdata -big
```
