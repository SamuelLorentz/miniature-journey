# Conway's Game of Life: Tech Challenge

Build a simulator for Conway's Game of Life on a finite grid, in the language of your choice, in about 3 to 4 hours. Correctness comes first, then design, tests and performance.

## The rules

The world is a 2D grid of cells, each alive or dead. All cells update at once: every next-generation state is computed from the current generation only. Each cell has 8 neighbors (horizontal, vertical and diagonal).

| Current state | Live neighbors | Next state |
| --- | --- | --- |
| Alive | 0 or 1 | Dead (underpopulation) |
| Alive | 2 or 3 | Alive (survival) |
| Alive | 4 or more | Dead (overpopulation) |
| Dead | exactly 3 | Alive (reproduction) |
| Dead | anything else | Dead |

**Edges.** Your engine must support two modes, selectable at runtime:

- **Dead border:** cells outside the grid count as dead.
- **Toroidal:** the grid wraps, so the top row neighbors the bottom row and the left column neighbors the right column.

## Requirements

Work through the levels in order. Level 1 and 2 are mandatory; levels 3 and 4 separate strong submissions from good ones.

| Level | Scope | Must deliver |
| --- | --- | --- |
| 1. Core engine | Pure function `step(grid) -> grid` | Correct rules, both edge modes, no mutation of the input grid |
| 2. CLI | Load a grid, run N generations, print the result | Parse the input format below, exit code 1 with a clear message on bad input |
| 3. Tests | Automated suite | Unit tests for each rule, plus the patterns in the test cases section |
| 4. Performance | Large and sparse grids | Meet the targets in the constraints section and explain your data structure |
| Stretch | Pick any | Pattern detection, infinite grid, RLE import, or a live terminal or web renderer |

Stretch goals in detail:

- **Stability detection:** stop early when the grid becomes a still life or a repeating cycle, and report the period.
- **Infinite grid:** no fixed bounds; store only live cells (for example, a set of coordinates).
- **RLE import:** read patterns in the standard Run Length Encoded format.

## Input and output format

The grid is plain text. The first line holds the width and height; each following line is one row, with `#` for alive and `.` for dead.

```text
5 5
.....
..#..
..#..
..#..
.....
```

The CLI takes the file, a generation count and an optional wrap flag. Output uses the same format, so it can be piped back in.

```text
life <file> --generations <N> [--wrap]
```

Rules for parsing: width and height are between 1 and 10,000; every row must match the declared width; any other character is an error. Print only the final generation unless `--all` is given, which prints every generation separated by a blank line.

## Test cases

Your suite must include at least these five cases. The expected outputs below are exact; use them as acceptance tests.

| # | Pattern | Mode | Generations | Expected result |
| --- | --- | --- | --- | --- |
| T1 | Block (still life) | Dead border | any | Unchanged |
| T2 | Blinker (oscillator) | Dead border | 1 | Rotated 90 degrees; period 2 |
| T3 | Glider | Dead border | 4 | Moved one cell down and one cell right |
| T4 | Blinker across the left and right edges | Toroidal | 1 | Vertical blinker in column 0 |
| T5 | Same input as T4 | Dead border | 1 | Empty grid |

**T1: block, 4x4.** Input and every later generation:

```text
4 4
....
.##.
.##.
....
```

**T2: blinker, 5x5.** Input, then generation 1:

```text
5 5        5 5
.....      .....
..#..      .....
..#..      .###.
..#..      .....
.....      .....
```

**T3: glider, 6x6.** Input, then generation 4:

```text
6 6         6 6
.#....      ......
..#...      ..#...
###...      ...#..
......      .###..
......      ......
......      ......
```

**T4 and T5: edge behavior, 5x5.** Input (cells wrap from column 4 to column 0):

```text
5 5
.....
.....
##..#
.....
.....
```

With `--wrap`, generation 1 has live cells at rows 1, 2 and 3 of column 0 only. Without it, generation 1 is entirely empty.

Also add tests of your own for: a 1x1 grid, an all-alive grid, a malformed input file, and a random grid compared against a slow reference implementation.

## Constraints and performance

Targets are measured on a typical laptop with a random 30% density starting grid. Compiled languages use the first column; interpreted languages (Python, Ruby, JavaScript without typed arrays) use the second.

| Workload | Compiled | Interpreted |
| --- | --- | --- |
| 100 x 100, 1,000 generations | under 1 s | under 10 s |
| 1,000 x 1,000, 100 generations | under 10 s | under 90 s |

Other constraints:

- Standard library only for the engine; test frameworks and argument parsers are allowed.
- Hold at most two grids in memory at once (current and next) in the dense engine.
- No global mutable state; the same input must always give the same output.
- Document the time and memory complexity of `step` in your README, and name one optimization you chose not to make.

## Evaluation rubric

Submissions are scored out of 100. A submission that fails any of T1 to T5 cannot score above 60.

| Area | Weight | What reviewers look for |
| --- | --- | --- |
| Correctness | 35 | All rules right, both edge modes, T1 to T5 pass, simultaneous update (no in-place bugs) |
| Code design | 20 | Engine separate from I/O, small functions, clear names, no duplicated neighbor logic |
| Testing | 20 | Rule-level unit tests, edge cases, a reference-implementation comparison |
| Performance | 15 | Meets the targets, and the data structure choice is justified |
| Communication | 10 | README with run instructions, trade-offs and what you would do with more time |

Stretch goals earn up to 10 bonus points, capped at a total of 100.

## Submission and review

Send a git repository (or a zip) that runs with one documented command and passes its tests with another. Time-box yourself to 3 or 4 hours; unfinished work with honest notes beats a rushed stretch goal.

The README must cover:

1. How to build, run and test.
2. The grid representation you chose and why.
3. Time and memory complexity of one generation.
4. Known limitations and what you would do next.

The review is a 45-minute conversation. Expect to walk through your code, then extend it live, for example: add a new edge mode, change the rules to HighLife (a dead cell is born with 3 or 6 neighbors), or explain how you would parallelize `step`.
