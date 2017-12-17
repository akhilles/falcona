package falcona

type Move struct {
	move  uint32
	score int
}

type MoveList struct {
	moves [128]Move
	count int
}

type TargetFn func(sq int, occ uint64) uint64

const (
	isCapture   = 0x00F00000
	isPromotion = 0x0F000000
	isCastle    = 0x10000000
	isEnpassant = 0x20000000
	isPawnstart = 0x40000000
)

func (pos *Position) generateMoves() *MoveList {
	ml := &MoveList{[128]Move{}, 0}
	for piece := WhiteKnight ^ pos.side; piece <= WhiteKing; piece += 2 {
		pos.pieceMoves(ml, piece)
	}
	pos.castleMoves(ml)
	if pos.side == White {
		pos.whitePawnMoves(ml)
	} else {
		pos.blackPawnMoves(ml)
	}
	return ml
}

func from(move uint32) uint8 {
	return uint8(move) & 0xFF
}

func to(move uint32) uint8 {
	return uint8(move>>8) & 0xFF
}

func moved(move uint32) uint8 {
	return uint8(move>>16) & 0xF
}

func captured(move uint32) uint8 {
	return uint8(move>>20) & 0xF
}

func promoted(move uint32) uint8 {
	return uint8(move>>24) & 0xF
}

var targetFns = [5]TargetFn{
	knightTargets, bishopTargets, rookTargets, queenTargets, kingTargets,
}

func newMove(from, to int, piece, captured, promoted uint8, flags uint32) uint32 {
	return uint32(from) | (uint32(to) << 8) | (uint32(piece) << 16) | (uint32(captured) << 20) | (uint32(promoted) << 24) | flags
}

func (ml *MoveList) addQuietMove(from, to int, piece, promoted uint8, flags uint32) {
	move := newMove(from, to, piece, 0xF, promoted, flags)
	ml.moves[ml.count] = Move{move, 3}
	ml.count++
}

func (ml *MoveList) addCaptureMove(from, to int, piece, captured, promoted uint8, flags uint32) {
	move := newMove(from, to, piece, captured, promoted, flags)
	ml.moves[ml.count] = Move{move, 3}
	ml.count++
}

func (pos *Position) pieceMoves(ml *MoveList, piece uint8) {
	pieces := pos.pieces[piece]
	var from, to int
	for pieces != 0 {
		pieces, from = pop(pieces)
		targets := targetFns[(piece>>1)-1](from, pos.occ)
		captures := targets & pos.colors[pos.side^1]
		targets = targets & ^(pos.occ)
		for captures != 0 {
			captures, to = pop(captures)
			captured := pos.findPiece(uint(to))
			ml.addCaptureMove(from, to, piece, captured, 0xF, 0)
		}
		for targets != 0 {
			targets, to = pop(targets)
			ml.addQuietMove(from, to, piece, 0xF, 0)
		}
	}
}

func (pos *Position) castleMoves(ml *MoveList) {
	from := pos.kings[pos.side]
	if ((pos.castles & castleKingside[pos.side]) != 0) && (pos.occ&castleKingsideEmpty[pos.side] == 0) {
		ml.addQuietMove(int(from), int(from+2), king(pos.side), 0xF, isCastle)
	}
	if ((pos.castles & castleQueenside[pos.side]) != 0) && (pos.occ&castleQueensideEmpty[pos.side] == 0) {
		ml.addQuietMove(int(from), int(from-3), king(pos.side), 0xF, isCastle)
	}
}

func (pos *Position) whitePawnMoves(ml *MoveList) {
	var from, to int

	forward1 := (pos.pieces[WhitePawn] << 8) & ^pos.occ
	forward2 := ((forward1 & maskRank[R3]) << 8) & ^pos.occ
	rightCapture := (pos.pieces[WhitePawn] << 9) & (pos.colors[pos.side^1] | (1 << pos.enpassant)) & ^maskFile[FA]
	leftCapture := (pos.pieces[WhitePawn] << 7) & (pos.colors[pos.side^1] | (1 << pos.enpassant)) & ^maskFile[FH]

	for forward2 != 0 {
		forward2, to = pop(forward2)
		from = to - 16
		ml.addQuietMove(from, to, WhitePawn, 0xF, isPawnstart)
	}

	for forward1 != 0 {
		forward1, to = pop(forward1)
		from = to - 8
		if to >= A8 {
			ml.addQuietMove(from, to, WhitePawn, WhiteQueen, 0)
			ml.addQuietMove(from, to, WhitePawn, WhiteKnight, 0)
			ml.addQuietMove(from, to, WhitePawn, WhiteRook, 0)
			ml.addQuietMove(from, to, WhitePawn, WhiteBishop, 0)
		} else {
			ml.addQuietMove(from, to, WhitePawn, 0xF, 0)
		}
	}

	for rightCapture != 0 {
		rightCapture, to = pop(rightCapture)
		from = to - 9
		captured := pos.findPiece(uint(to))
		if to >= A8 || to <= H1 {
			ml.addCaptureMove(from, to, WhitePawn, captured, WhiteQueen, 0)
			ml.addCaptureMove(from, to, WhitePawn, captured, WhiteKnight, 0)
			ml.addCaptureMove(from, to, WhitePawn, captured, WhiteRook, 0)
			ml.addCaptureMove(from, to, WhitePawn, captured, WhiteBishop, 0)
		} else {
			if to == int(pos.enpassant) {
				ml.addCaptureMove(from, to, WhitePawn, BlackPawn, 0xF, isEnpassant)
			} else {
				ml.addCaptureMove(from, to, WhitePawn, captured, 0xF, 0)
			}
		}
	}

	for leftCapture != 0 {
		leftCapture, to = pop(leftCapture)
		from = to - 7
		captured := pos.findPiece(uint(to))
		if to >= A8 || to <= H1 {
			ml.addCaptureMove(from, to, WhitePawn, captured, WhiteQueen, 0)
			ml.addCaptureMove(from, to, WhitePawn, captured, WhiteKnight, 0)
			ml.addCaptureMove(from, to, WhitePawn, captured, WhiteRook, 0)
			ml.addCaptureMove(from, to, WhitePawn, captured, WhiteBishop, 0)
		} else {
			if to == int(pos.enpassant) {
				ml.addCaptureMove(from, to, WhitePawn, BlackPawn, 0xF, isEnpassant)
			} else {
				ml.addCaptureMove(from, to, WhitePawn, captured, 0xF, 0)
			}
		}
	}
}

func (pos *Position) blackPawnMoves(ml *MoveList) {
	var from, to int

	forward1 := (pos.pieces[BlackPawn] >> 8) & ^pos.occ
	forward2 := ((forward1 & maskRank[R6]) >> 8) & ^pos.occ
	rightCapture := (pos.pieces[BlackPawn] >> 7) & (pos.colors[pos.side^1] | (1 << pos.enpassant)) & ^maskFile[FA]
	leftCapture := (pos.pieces[BlackPawn] >> 9) & (pos.colors[pos.side^1] | (1 << pos.enpassant)) & ^maskFile[FH]

	for forward2 != 0 {
		forward2, to = pop(forward2)
		from = to + 16
		ml.addQuietMove(from, to, BlackPawn, 0xF, isPawnstart)
	}

	for forward1 != 0 {
		forward1, to = pop(forward1)
		from = to + 8
		if to >= A8 {
			ml.addQuietMove(from, to, BlackPawn, BlackQueen, 0)
			ml.addQuietMove(from, to, BlackPawn, BlackKnight, 0)
			ml.addQuietMove(from, to, BlackPawn, BlackRook, 0)
			ml.addQuietMove(from, to, BlackPawn, BlackBishop, 0)
		} else {
			ml.addQuietMove(from, to, BlackPawn, 0xF, 0)
		}
	}

	for rightCapture != 0 {
		rightCapture, to = pop(rightCapture)
		from = to + 7
		captured := pos.findPiece(uint(to))
		if to >= A8 || to <= H1 {
			ml.addCaptureMove(from, to, BlackPawn, captured, BlackQueen, 0)
			ml.addCaptureMove(from, to, BlackPawn, captured, BlackKnight, 0)
			ml.addCaptureMove(from, to, BlackPawn, captured, BlackRook, 0)
			ml.addCaptureMove(from, to, BlackPawn, captured, BlackBishop, 0)
		} else {
			if to == int(pos.enpassant) {
				ml.addCaptureMove(from, to, BlackPawn, WhitePawn, 0xF, isEnpassant)
			} else {
				ml.addCaptureMove(from, to, BlackPawn, captured, 0xF, 0)
			}
		}
	}

	for leftCapture != 0 {
		leftCapture, to = pop(leftCapture)
		from = to + 9
		captured := pos.findPiece(uint(to))
		if to >= A8 || to <= H1 {
			ml.addCaptureMove(from, to, BlackPawn, captured, BlackQueen, 0)
			ml.addCaptureMove(from, to, BlackPawn, captured, BlackKnight, 0)
			ml.addCaptureMove(from, to, BlackPawn, captured, BlackRook, 0)
			ml.addCaptureMove(from, to, BlackPawn, captured, BlackBishop, 0)
		} else {
			if to == int(pos.enpassant) {
				ml.addCaptureMove(from, to, BlackPawn, WhitePawn, 0xF, isEnpassant)
			} else {
				ml.addCaptureMove(from, to, BlackPawn, captured, 0xF, 0)
			}
		}
	}
}

func rookTargets(sq int, occ uint64) uint64 {
	occ = rookMagic[sq].mask & occ
	index := (occ * rookMagic[sq].magic) >> 52
	return rookMagicMoves[sq][index]
}

func bishopTargets(sq int, occ uint64) uint64 {
	occ = bishopMagic[sq].mask & occ
	index := (occ * bishopMagic[sq].magic) >> 55
	return bishopMagicMoves[sq][index]
}

func queenTargets(sq int, occ uint64) uint64 {
	return rookTargets(sq, occ) | bishopTargets(sq, occ)
}

func knightTargets(sq int, occ uint64) uint64 {
	return knightMoves[sq]
}

func kingTargets(sq int, occ uint64) uint64 {
	return kingMoves[sq]
}

func (pos *Position) attackedBy(sq int, color uint8) bool {
	if pawnAttacks[color^1][sq]&pos.pieces[pawn(color)] != 0 {
		return true
	}
	for piece := WhiteKnight ^ color; piece <= WhiteKing; piece += 2 {
		targets := targetFns[(piece>>1)-1](sq, pos.occ)
		if (targets & pos.pieces[piece]) != 0 {
			return true
		}
	}
	return false
}
