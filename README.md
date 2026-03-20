# Vizzini

> *"You fool! You fell victim to one of the classic blunders!"*

A UCI-compliant chess engine written in Go — complete with search, evaluation, 
and time management. Play against it in the interactive CLI or plug it into 
your favorite chess GUI.

## Features

> *"Never go in against a Sicilian when death is on the line!"*

- **Bitboard board representation** — 64-bit integers for fast, efficient piece tracking
- **Full legal move generation** — all piece types, castling, en passant, promotions, the works
- **Magic bitboards** — precomputed sliding piece attacks for bishops, rooks, and queens
- **Zobrist hashing** — incremental position hashing for rapid identification
- **UCI move notation parsing** — speak the universal chess language
- **Move/undo history** — full state restoration, because sometimes you need a do-over
- **Interactive CLI** — Unicode board display so you can admire your position
- **Perft testing framework** — correctness validation against known positions
- **UCI protocol** — full UCI compliance for chess GUI integration
- **Negamax search with alpha-beta pruning** — iterative deepening, aspiration windows, null move pruning, late move reductions, and check extensions
- **Position evaluation** — piece-square tables, material scoring, evaluation noise, and temperature-based move selection
- **Move ordering** — PV move, MVV-LVA, killer moves, and history heuristic
- **Quiescence search** — tactical stability at search boundaries
- **Draw detection** — repetition detection and 50-move rule
- **Time management** — configurable depth, node limits, and clock-aware search
- **Zero external dependencies** — *"I built this with nothing but my wits and a Go compiler"*

## Search

> *"INCONCEIVABLE!"*

> *"You keep using that word. I do not think it means what you think it means."*

Vizzini uses **negamax with alpha-beta pruning**, enhanced with several techniques to search deeper and faster:

- **Iterative deepening** — searches progressively deeper, allowing time-controlled play and seeding move ordering from shallower iterations
- **Aspiration windows** — from depth 4+, narrows the search window to ±50cp around the previous score; re-searches with a full window on fail
- **Null move pruning** — skips a turn to test if the position is already winning; reduction R=2 (R=3 at depth ≥6), disabled in check or with only pawns
- **Late move reductions (LMR)** — quiet moves searched after the first 4 get reduced depth; re-searched at full depth if they beat alpha
- **Check extensions** — extends depth by 1 when in check so tactical lines aren't cut short
- **Quiescence search** — continues searching captures and check evasions at leaf nodes until the position is tactically quiet
- **PV table** — million-entry hash table storing the best move per position for move ordering and PV line reconstruction
- **Killer moves** — two quiet refutation moves remembered per depth level
- **History heuristic** — tracks quiet moves that improve alpha, scaled into move ordering priority
- **Draw detection** — returns a draw score on repetition or when the 50-move clock reaches 100

## Evaluation

> *"You're trying to trick me into giving away something. It won't work."*

Position evaluation combines **material scoring** with **piece-square tables** — each piece type has a 64-entry bonus/penalty table rewarding good placement (centralized knights, rooks on the seventh rank, safe kings, etc.).

Two features add variety to play:

- **Evaluation noise** — uniform random noise of ±5 centipawns on leaf evaluations prevents perfectly deterministic play from identical positions
- **Temperature-based move selection** — during the first 10 moves, a softmax distribution (temperature 0.2) selects probabilistically among root moves within 75cp of the best, producing more varied openings; disabled when a forced mate is found

## Getting Started

> *"Let me explain... No, there is too much. Let me sum up."*

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

Vizzini supports two modes, detected automatically on startup:

### UCI Mode

Type `uci` at the prompt to enter UCI mode. This is the standard interface for chess GUIs like Arena, CuteChess, or Banksia — just point the GUI at the Vizzini binary.

### Player vs Engine

> *"I can't compete with you physically, and you're no match for my brains."*

Type anything else (or just press Enter) to play an interactive game against Vizzini. You submit moves in UCI notation and the engine responds with its own.

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

## API

> *"I challenge you to a battle of wits."*

Vizzini exposes a JSON API for chess analysis. Start the server:

```sh
./bin/vizzini serve
# or
just serve
```

The server listens on port **8080**. All POST endpoints accept and return `application/json`.

### Endpoints

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/health` | Health check |
| `POST` | `/validmoves` | Legal moves for a FEN position |
| `POST` | `/submitmove` | Submit a player move, get updated position |
| `POST` | `/bestmove` | Get the engine's best move (without applying it) |
| `POST` | `/submitbestmove` | Get the engine's best move and apply it |

### `GET /health`

Returns `{"status": "ok"}`.

### `POST /validmoves`

Get all legal moves for a position.

```json
// Request
{"fen": "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"}

// Response
{
  "fen": "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
  "side_to_move": "black",
  "status": "active",
  "move_count": 20,
  "moves": [
    {"uci": "a7a6", "san": "a6", "from": "a7", "to": "a6", "capture": false, "castling": false, "check": false, "promotion": null}
  ]
}
```

### `POST /submitmove`

Submit a move in UCI notation and get the resulting position with its legal moves.

```json
// Request
{"fen": "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", "move": "e2e4"}

// Response
{
  "uci": "e2e4", "san": "e4", "from": "e2", "to": "e4",
  "fen": "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
  "status": "active", "side_to_move": "black",
  "move_count": 20, "moves": [...]
}
```

### `POST /bestmove`

Ask the engine for its best move without applying it. Optionally control search depth or timeout.

```json
// Request
{"fen": "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1", "depth": 6}

// Response
{
  "fen": "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
  "depth": 6, "uci": "e7e5", "san": "e5",
  "from": "e7", "to": "e5",
  "score": 0, "nodes": 45832, "source": "search"
}
```

| Field | Type | Description |
| --- | --- | --- |
| `depth` | int | Max search depth (capped at 20) |
| `timeout_ms` | int | Search time limit in milliseconds (default: 5 000 ms) |

### `POST /submitbestmove`

Same as `/bestmove`, but also applies the move and returns the resulting position with legal moves — combines `/bestmove` + `/submitmove` in one call.

```json
// Request
{"fen": "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"}

// Response
{
  "uci": "e7e5", "san": "e5", "from": "e7", "to": "e5",
  "fen": "rnbqkbnr/pppp1ppp/8/8/4Pp2/8/PPPP1PPP/RNBQKBNR w KQkq ...",
  "status": "active", "side_to_move": "white",
  "move_count": 29, "moves": [...],
  "depth": 5, "score": 0, "nodes": 12345, "source": "search"
}
```

### Errors

All errors return a consistent shape:

```json
{"error": "invalid_fen", "message": "fen field is required"}
```

Error codes: `invalid_request`, `missing_fen`, `invalid_fen`, `missing_move`, `invalid_move`, `no_move`, `not_found`.

### Configuration

| Variable | Effect | Default |
| --- | --- | --- |
| `CORS_PERMISSIVE` | Set to any value to allow all origins (`*`) | Only `https://chess.cyberdan.dev` |

Body size is capped at 1 MB. CORS preflight (`OPTIONS`) is handled automatically.

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
| `just serve` | `just s` | Build and start the HTTP API server |

## Architecture

All source lives in `src/`:

| File | Purpose |
| --- | --- |
| `main.go` | Entry point, UCI mode, and interactive player-vs-engine mode |
| `types.go` | Core type definitions — `Square`, `Color`, `Move`, `Bitboard`, `Hash`, etc. |
| `const.go` | Game constants and enumerations |
| `bitboard.go` | Bitboard operations (set/clear/pop bits, LS1B) |
| `bitboard_tables.go` | Precomputed attack tables and magic bitboard initialization |
| `board.go` | Board state, FEN parsing, Unicode board display |
| `move.go` | Move encoding (bit-packed `uint32`), make/undo move logic |
| `move_gen.go` | Legal move generation for all piece types |
| `board_hash.go` | Zobrist hashing — incremental hash updates |
| `search.go` | Negamax search — alpha-beta, iterative deepening, aspiration windows, NMP, LMR, PV table, move ordering |
| `uci.go` | UCI protocol — command parsing and engine communication |
| `evaluate.go` | Position evaluation — piece-square tables, material scoring |
| `server.go` | HTTP server, API handlers, CORS and body-limit middleware |
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

