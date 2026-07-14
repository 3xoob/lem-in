# lem-in

A Go implementation of the 42-school `lem-in` project: given a colony of
rooms connected by tunnels, find the quickest way to move a given number of
ants from the `start` room to the `end` room, moving only one ant per room
per turn (except through `start` and `end`).

## Overview

The program reads a text description of an "ant farm" (rooms, tunnels, a
start room, an end room, and a number of ants), computes the maximum set of
room-disjoint paths between `start` and `end`, distributes the ants across
those paths, and prints the turn-by-turn movement of every ant until they
have all reached `end`.

## Features

- Parses the `.txt` ant-farm format, including room coordinates, the
  `##start` / `##end` markers, and `-` separated tunnel links.
- Validates the input and reports specific `ERROR: invalid data format, ...`
  messages for: non-positive/non-numeric ant counts, duplicate room names,
  malformed room or link lines, self-links, duplicate links, links to
  unknown rooms, and missing start/end rooms.
- Computes the maximum number of room-disjoint paths from `start` to `end`
  using an Edmonds-Karp (BFS-based) max-flow algorithm, then decomposes the
  resulting flow into individual paths.
- Assigns ants to paths greedily (minimizing each ant's projected finish
  turn) and simulates their movement turn by turn.
- Echoes the original input file content before printing the simulated
  movements.

## Technologies

- Go (module targets Go 1.23; developed/tested with the Go toolchain)
- Go's standard `testing` package for unit and end-to-end tests

## Project Structure

```
main.go                       # Entry point: reads the file, echoes it, and drives the pipeline
Functions/
  Struct.go                   # Room, Path, and AntFarm data structures
  ParseFile.go                # Input parsing and validation (parseAntFarm)
  BFS.go                      # Room-disjoint path computation (Edmonds-Karp max-flow)
  AntsMovement.go             # Ant-to-path assignment and turn-by-turn simulation
  lemin_test.go                # Unit and end-to-end tests
examples/                     # Sample ant-farm input files (valid and invalid)
test.txt                      # An additional sample ant-farm input file
go.mod                        # Go module definition (module lemin, go 1.23)
```

## Requirements

- Go 1.23 or later

## Installation

```bash
git clone https://github.com/3xoob/lem-in
cd lem-in
go build -o lem-in .
```

## Usage

Run the built binary with the path to an ant-farm file as its only argument:

```bash
./lem-in examples/example00.txt
```

The program requires exactly one argument (the file path); otherwise it
prints `ERROR: invalid number of arguments`. If the file cannot be opened,
it prints `ERROR: failed to open file.`.

### Input format

```
<number of ants>
##start
<room name> <x> <y>
...
##end
<room name> <x> <y>
...
<room name>-<room name>
...
```

- The first non-empty line is the number of ants (must be a positive
  integer).
- `##start` and `##end` mark the room definition immediately following them.
- Room lines have the form `name x y`.
- Link lines have the form `room1-room2`.

### Output format

The program first echoes the raw input file, then prints one line per turn.
Each line lists every ant that moved that turn as space-separated
`Lant-room` tokens (e.g. `L1-2 L2-3`), until all ants have reached `end`.

## Testing

Unit and end-to-end tests live in `Functions/lemin_test.go` and cover input
parsing (valid and invalid cases), room-disjoint path computation, and
full simulations against the files in `examples/` (including the two
`badexample*.txt` files, which are expected to be rejected by the parser).

Run them with:

```bash
go test ./...
```

## Example

Given `examples/example00.txt`:

```
4
##start
0 0 3
2 2 5
3 4 0
##end
1 8 3
0-2
2-3
3-1
```

Running `./lem-in examples/example00.txt` prints the input followed by the
turn-by-turn moves, e.g.:

```
L1-2
L1-3 L2-2
L1-1 L2-3 L3-2
L2-1 L3-3 L4-2
L3-1 L4-3
L4-1
```

## Learning Objectives

This project is part of the 42-school curriculum. Based on the code, it
exercises:

- Parsing and validating a custom, line-oriented text format with strict
  error reporting.
- Graph modeling of rooms and tunnels.
- Applying a max-flow algorithm (Edmonds-Karp) to find the maximum number of
  room-disjoint paths between two nodes.
- Turn-based simulation and scheduling of multiple agents (ants) across
  several concurrent paths under a one-ant-per-room constraint.
- Writing Go unit and end-to-end tests, including invariant checks on
  simulated output (no shared rooms, no tunnel used twice per turn, every
  ant reaches the destination).

## Limitations

- If no path exists between `start` and `end`, the program prints
  `ERROR: no valid paths found` and exits without producing any moves.
- The parser enforces a specific format: it rejects non-positive/invalid
  ant counts, duplicate room names, duplicate links, self-links, links to
  undefined rooms, malformed room/link lines, and missing start or end
  rooms.

## License

This repository includes a `LICENSE` / `COPYRIGHT.md` file. The source code
is made publicly available for portfolio and viewing purposes only; no
permission is granted to copy, modify, distribute, or reuse it without
prior written permission from the copyright holder.
