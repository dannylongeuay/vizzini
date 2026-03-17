# Vizzini

> *"You fool! You fell victim to one of the classic blunders!"*

A chess engine written in Go. No search, no evaluation — just a pure, fully legal move generator with an interactive CLI. Think of it as the foundation upon which a truly *inconceivable* chess engine could be built.

## Features

- **Bitboard board representation** — 64-bit integers for fast, efficient piece tracking
- **Full legal move generation** — all piece types, castling, en passant, promotions, the works
- **Magic bitboards** — precomputed sliding piece attacks for bishops, rooks, and queens
- **Zobrist hashing** — incremental position hashing for rapid identification
- **UCI move notation parsing** — speak the universal chess language
- **Move/undo history** — full state restoration, because sometimes you need a do-over
- **Interactive CLI** — Unicode board display so you can admire your position
- **Perft testing framework** — correctness validation against known positions
- **Zero external dependencies** — *"I built this with nothing but my wits and a Go compiler"*

## Getting Started

### Prerequisites

- **Go 1.25+**
- Optionally: [Nix](https://nixos.org/) + [just](https://github.com/casey/just) for the full dev experience

### Clone

```sh
git clone https://github.com/dannylongeuay/vizzini.git
cd vizzini
```

### Build & Run

Using `just`:

```sh
just run
```

Or with raw Go:

```sh
go build -o bin/vizzini ./src/...
./bin/vizzini
```

## Usage

Vizzini launches an interactive CLI where you can play moves in UCI notation. The board is displayed with Unicode chess pieces after each move.

```
_________________________
|♜ |♞ |♝ |♛ |♚ |♝ |♞ |♜ | 8
_________________________
|♙ |♙ |♙ |♙ |♙ |♙ |♙ |♙ | 7
_________________________
|  |  |  |  |  |  |  |  | 6
_________________________
|  |  |  |  |  |  |  |  | 5
_________________________
|  |  |  |  |  |  |  |  | 4
_________________________
|  |  |  |  |  |  |  |  | 3
_________________________
|♟ |♟ |♟ |♟ |♟ |♟ |♟ |♟ | 2
_________________________
|♜ |♞ |♝ |♛ |♚ |♝ |♞ |♜ | 1
_________________________
 A  B  C  D  E  F  G  H

Submit move: e2e4
```

Moves use standard UCI format: `<from><to>[promotion]`

- `e2e4` — pawn to e4
- `g1f3` — knight to f3
- `e7e8q` — pawn promotes to queen

## Development

### Nix Dev Shell

```sh
nix develop
```

This drops you into a shell with Go, golangci-lint, goimports, and just pre-configured.

### Task Runner

All tasks are available through `just`:

| Command | Alias | Description |
| --- | --- | --- |
| `just build` | `just b` | Build binary to `bin/vizzini` |
| `just run` | `just r` | Build and run |
| `just test` | `just t` | Run short tests |
| `just test-long` | `just tl` | Run all tests (including long Perft suites) |
| `just lint` | `just l` | Lint with golangci-lint |
| `just format` | `just f` | Format with goimports |

## Architecture

All source lives in `src/`:

| File | Purpose |
| --- | --- |
| `main.go` | Entry point and interactive game loop |
| `types.go` | Core type definitions — `Square`, `Color`, `Move`, `Bitboard`, `Hash`, etc. |
| `const.go` | Game constants and enumerations |
| `bitboard.go` | Bitboard operations (set/clear/pop bits, LS1B) |
| `bitboard_tables.go` | Precomputed attack tables and magic bitboard initialization |
| `board.go` | Board state, FEN parsing, Unicode board display |
| `move.go` | Move encoding (bit-packed `uint32`), make/undo move logic |
| `move_gen.go` | Legal move generation for all piece types |
| `board_hash.go` | Zobrist hashing — incremental hash updates |
| `util.go` | UCI parsing, coordinate helpers, utility functions |

## Testing

Run the short test suite:

```sh
just test
```

Run the full suite including Perft validation against known positions.

```sh
just test-long
```

Perft tests walk the move tree to a given depth and compare total node counts against established results — the gold standard for verifying move generation correctness.

