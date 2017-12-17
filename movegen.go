package falcona

type Move struct {
	move  uint32
	score int
}

type MoveList struct {
	moves [128]Move
	count int
}

const (
	isCapture   = 0x00F00000
	isPromotion = 0x0F000000
	isCastle    = 0x10000000
	isEnpassant = 0x20000000
	isPawnstart = 0x40000000
)

func newMove(from, to int, piece, captured, promoted uint8, flags uint32) uint32 {
	return uint32(from) | (uint32(to) << 8) | (uint32(piece) << 16) | (uint32(captured) << 20) | (uint32(promoted) << 24) | (flags << 28)
}

func (ml *MoveList) addQuietMove(from, to int, piece, promoted uint8, flags uint32) {
	move := newMove(from, to, piece, 0, promoted, flags)
	ml.moves[ml.count] = Move{move, 3}
	ml.count++
}

func (ml *MoveList) addCaptureMove(from, to int, piece, captured, promoted uint8, flags uint32) {
	move := newMove(from, to, piece, captured, promoted, flags)
	ml.moves[ml.count] = Move{move, 3}
	ml.count++
}

func (pos *Position) generateMoves() *MoveList {
	ml := MoveList{[128]Move{}, 0}

}

func (pos *Position) generator(ml *MoveList, piece uint8, targets uint32) {

}

func (pos *Position) generateRookMoves(ml *MoveList) {
	rooks := pos.pieces[rook(pos.side)]
	var from, to int
	for rooks != 0 {
		rooks, from = pop(rooks)
		targets := pos.rookTargets(from)
		captures := targets & pos.colors[pos.side^1]
		targets = targets & ^(pos.occ)
		for captures != 0 {
			captures, to = pop(targets)
			captured := pos.findPiece(uint(to))
			ml.addCaptureMove(from, to, rook(pos.side), captured, 0, 0)
		}
		for targets != 0 {
			targets, to = pop(targets)
			ml.addQuietMove(from, to, rook(pos.side), 0, 0)
		}
	}
}

func (pos *Position) rookTargets(sq int) uint64 {
	occ := rookMagic[sq].mask & pos.occ
	index := (occ * rookMagic[sq].magic) >> 52
	return rookMagicMoves[sq][index]
}

func (pos *Position) bishopTargets(sq int) uint64 {
	occ := bishopMagic[sq].mask & pos.occ
	index := (occ * bishopMagic[sq].magic) >> 55
	return bishopMagicMoves[sq][index]
}

func (pos *Position) attackedBy(sq int, color uint8) bool {
	if knightMoves[sq]&pos.pieces[knight(color)] != 0 {
		return true
	}
	if pawnAttacks[color^1][sq]&pos.pieces[pawn(color)] != 0 {
		return true
	}
	if pos.bishopTargets(sq)&(pos.pieces[bishop(color)]|pos.pieces[queen(color)]) != 0 {
		return true
	}
	if pos.rookTargets(sq)&(pos.pieces[rook(color)]|pos.pieces[queen(color)]) != 0 {
		return true
	}
	if kingMoves[sq]&pos.pieces[king(color)] != 0 {
		return true
	}
	return false
}
