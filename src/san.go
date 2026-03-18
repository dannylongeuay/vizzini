package main

import (
	"fmt"
	"strings"
)

// ToSAN converts a move to Standard Algebraic Notation.
func (b *Board) ToSAN(move Move) string {
	var mu MoveUnpacked
	move.Unpack(&mu)

	// Castling
	if mu.IsCastling() {
		if mu.moveKind == KING_CASTLE {
			return "O-O"
		}
		return "O-O-O"
	}

	var sb strings.Builder

	isPawn := mu.originSquare == WHITE_PAWN || mu.originSquare == BLACK_PAWN

	// Piece letter
	var pieceChar string
	if !isPawn {
		switch mu.originSquare {
		case WHITE_KNIGHT, BLACK_KNIGHT:
			pieceChar = "N"
		case WHITE_BISHOP, BLACK_BISHOP:
			pieceChar = "B"
		case WHITE_ROOK, BLACK_ROOK:
			pieceChar = "R"
		case WHITE_QUEEN, BLACK_QUEEN:
			pieceChar = "Q"
		case WHITE_KING, BLACK_KING:
			pieceChar = "K"
		}
		sb.WriteString(pieceChar)
	}

	isCapture := mu.IsCapture()

	// Disambiguation for non-pawn pieces
	if !isPawn {
		legal := b.LegalMoves()
		var ambiguous []Move
		for _, m := range legal {
			if m == move {
				continue
			}
			var am MoveUnpacked
			m.Unpack(&am)
			if am.originSquare == mu.originSquare && am.dstCoord == mu.dstCoord {
				ambiguous = append(ambiguous, m)
			}
		}
		if len(ambiguous) > 0 {
			originFile := int(mu.originCoord) % 8
			originRank := int(mu.originCoord) / 8
			sameFile := false
			sameRank := false
			for _, am := range ambiguous {
				var amu MoveUnpacked
				am.Unpack(&amu)
				if int(amu.originCoord)%8 == originFile {
					sameFile = true
				}
				if int(amu.originCoord)/8 == originRank {
					sameRank = true
				}
			}
			if !sameFile {
				sb.WriteByte(byte('a' + originFile))
			} else if !sameRank {
				sb.WriteByte(byte('1' + originRank))
			} else {
				sb.WriteByte(byte('a' + originFile))
				sb.WriteByte(byte('1' + originRank))
			}
		}
	}

	// Pawn capture: always include source file
	if isPawn && isCapture {
		originFile := int(mu.originCoord) % 8
		sb.WriteByte(byte('a' + originFile))
	}

	if isCapture {
		sb.WriteByte('x')
	}

	// Destination square (lowercase)
	sb.WriteString(stringCoordLower(mu.dstCoord))

	// Promotion suffix
	switch mu.moveKind {
	case QUEEN_PROMOTION, QUEEN_PROMOTION_CAPTURE:
		sb.WriteString("=Q")
	case ROOK_PROMOTION, ROOK_PROMOTION_CAPTURE:
		sb.WriteString("=R")
	case BISHOP_PROMOTION, BISHOP_PROMOTION_CAPTURE:
		sb.WriteString("=B")
	case KNIGHT_PROMOTION, KNIGHT_PROMOTION_CAPTURE:
		sb.WriteString("=N")
	}

	// Check / checkmate suffix
	if err := b.MakeMove(move); err == nil {
		if b.InCheck() {
			legal := b.LegalMoves()
			if len(legal) == 0 {
				sb.WriteByte('#')
			} else {
				sb.WriteByte('+')
			}
		}
		b.UndoMove()
	}

	return sb.String()
}

// ParseSAN parses a move in Standard Algebraic Notation.
func (b *Board) ParseSAN(san string) (Move, error) {
	s := san

	// Strip trailing check/mate indicators
	s = strings.TrimRight(s, "+#")

	// Castling
	if s == "O-O-O" {
		legal := b.LegalMoves()
		for _, m := range legal {
			var mu MoveUnpacked
			m.Unpack(&mu)
			if mu.moveKind == QUEEN_CASTLE {
				return m, nil
			}
		}
		return 0, fmt.Errorf("invalid_move: %v", san)
	}
	if s == "O-O" {
		legal := b.LegalMoves()
		for _, m := range legal {
			var mu MoveUnpacked
			m.Unpack(&mu)
			if mu.moveKind == KING_CASTLE {
				return m, nil
			}
		}
		return 0, fmt.Errorf("invalid_move: %v", san)
	}

	// Parse promotion
	var promotionKinds [2]MoveKind // [quiet, capture]
	promotionKinds[0] = 255        // sentinel = no promotion
	promotionKinds[1] = 255
	if idx := strings.Index(s, "="); idx >= 0 {
		promChar := s[idx+1:]
		s = s[:idx]
		switch promChar {
		case "Q":
			promotionKinds[0] = QUEEN_PROMOTION
			promotionKinds[1] = QUEEN_PROMOTION_CAPTURE
		case "R":
			promotionKinds[0] = ROOK_PROMOTION
			promotionKinds[1] = ROOK_PROMOTION_CAPTURE
		case "B":
			promotionKinds[0] = BISHOP_PROMOTION
			promotionKinds[1] = BISHOP_PROMOTION_CAPTURE
		case "N":
			promotionKinds[0] = KNIGHT_PROMOTION
			promotionKinds[1] = KNIGHT_PROMOTION_CAPTURE
		default:
			return 0, fmt.Errorf("invalid_move: %v", san)
		}
	}

	// Remove 'x' (capture indicator)
	s = strings.ReplaceAll(s, "x", "")

	if len(s) < 2 {
		return 0, fmt.Errorf("invalid_move: %v", san)
	}

	// Extract destination (last 2 chars)
	dstStr := s[len(s)-2:]
	dstCoord, err := StringToCoord(dstStr)
	if err != nil {
		return 0, fmt.Errorf("invalid_move: %v", san)
	}

	// Determine piece type and disambiguation
	rest := s[:len(s)-2]

	var pieceSquares []Square // possible origin squares for this piece type
	var disambigFile int = -1
	var disambigRank int = -1

	isPawn := true
	if len(rest) > 0 && rest[0] >= 'A' && rest[0] <= 'Z' {
		isPawn = false
		pieceChar := rest[0]
		rest = rest[1:]
		switch b.sideToMove {
		case WHITE:
			switch pieceChar {
			case 'N':
				pieceSquares = []Square{WHITE_KNIGHT}
			case 'B':
				pieceSquares = []Square{WHITE_BISHOP}
			case 'R':
				pieceSquares = []Square{WHITE_ROOK}
			case 'Q':
				pieceSquares = []Square{WHITE_QUEEN}
			case 'K':
				pieceSquares = []Square{WHITE_KING}
			default:
				return 0, fmt.Errorf("invalid_move: %v", san)
			}
		case BLACK:
			switch pieceChar {
			case 'N':
				pieceSquares = []Square{BLACK_KNIGHT}
			case 'B':
				pieceSquares = []Square{BLACK_BISHOP}
			case 'R':
				pieceSquares = []Square{BLACK_ROOK}
			case 'Q':
				pieceSquares = []Square{BLACK_QUEEN}
			case 'K':
				pieceSquares = []Square{BLACK_KING}
			default:
				return 0, fmt.Errorf("invalid_move: %v", san)
			}
		}
	} else {
		// Pawn
		if b.sideToMove == WHITE {
			pieceSquares = []Square{WHITE_PAWN}
		} else {
			pieceSquares = []Square{BLACK_PAWN}
		}
	}

	// Parse disambiguation from rest
	for _, ch := range rest {
		if ch >= 'a' && ch <= 'h' {
			disambigFile = int(ch - 'a')
		} else if ch >= '1' && ch <= '8' {
			disambigRank = int(ch - '1')
		}
	}

	// Filter legal moves
	legal := b.LegalMoves()
	var matches []Move
	for _, m := range legal {
		var mu MoveUnpacked
		m.Unpack(&mu)

		// Piece type match
		matched := false
		for _, ps := range pieceSquares {
			if mu.originSquare == ps {
				matched = true
				break
			}
		}
		if !matched {
			continue
		}

		// Destination match
		if mu.dstCoord != dstCoord {
			continue
		}

		// Pawn source file (for captures and disambiguation)
		if isPawn && disambigFile >= 0 {
			if int(mu.originCoord)%8 != disambigFile {
				continue
			}
		}

		// Disambiguation file
		if !isPawn && disambigFile >= 0 {
			if int(mu.originCoord)%8 != disambigFile {
				continue
			}
		}

		// Disambiguation rank
		if disambigRank >= 0 {
			if int(mu.originCoord)/8 != disambigRank {
				continue
			}
		}

		// Promotion match
		if promotionKinds[0] != 255 {
			if mu.moveKind != promotionKinds[0] && mu.moveKind != promotionKinds[1] {
				continue
			}
		} else {
			// No promotion requested: exclude promotion moves
			switch mu.moveKind {
			case KNIGHT_PROMOTION, BISHOP_PROMOTION, ROOK_PROMOTION, QUEEN_PROMOTION,
				KNIGHT_PROMOTION_CAPTURE, BISHOP_PROMOTION_CAPTURE, ROOK_PROMOTION_CAPTURE, QUEEN_PROMOTION_CAPTURE:
				continue
			}
		}

		matches = append(matches, m)
	}

	if len(matches) == 0 {
		return 0, fmt.Errorf("invalid_move: %v", san)
	}
	if len(matches) > 1 {
		return 0, fmt.Errorf("ambiguous_move: %v", san)
	}
	return matches[0], nil
}

// ParseMove auto-detects UCI or SAN notation and parses accordingly.
func (b *Board) ParseMove(s string) (Move, error) {
	// Detect UCI: 4-5 chars, first 4 are [a-h][1-8][a-h][1-8], optional 5th is [qrbn]
	rs := []rune(s)
	if len(rs) >= 4 && len(rs) <= 5 {
		isUCI := rs[0] >= 'a' && rs[0] <= 'h' &&
			rs[1] >= '1' && rs[1] <= '8' &&
			rs[2] >= 'a' && rs[2] <= 'h' &&
			rs[3] >= '1' && rs[3] <= '8'
		if len(rs) == 5 {
			isUCI = isUCI && (rs[4] == 'q' || rs[4] == 'r' || rs[4] == 'b' || rs[4] == 'n')
		}
		if isUCI {
			return b.ParseUCIMove(s)
		}
	}
	return b.ParseSAN(s)
}
