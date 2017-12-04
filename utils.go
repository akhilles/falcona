package falcona

func abs(val int) int {
	if val < 0 {
		return -val
	}
	return val
}

func toRowCol(square int) (int, int) {
	return square >> 3, square & 7
}

func toSquare(row, col int) int {
	return (row << 3) + col
}

func pawn(color uint8) uint8 {
	return WhitePawn ^ color
}

func knight(color uint8) uint8 {
	return WhiteKnight ^ color
}

func bishop(color uint8) uint8 {
	return WhiteBishop ^ color
}

func rook(color uint8) uint8 {
	return WhiteRook ^ color
}

func queen(color uint8) uint8 {
	return WhiteQueen ^ color
}

func king(color uint8) uint8 {
	return WhiteKing ^ color
}

func color(piece uint8) uint8 {
	return (^piece) & 1
}

func isWhite(piece uint8) bool {
	return (piece & 1) == 1
}

func isBlack(piece uint8) bool {
	return (piece & 1) == 0
}

func isPawn(piece uint8) bool {
	return piece|1 == WhitePawn
}

func isKnight(piece uint8) bool {
	return piece|1 == WhiteKnight
}

func isBishop(piece uint8) bool {
	return piece|1 == WhiteBishop
}

func isRook(piece uint8) bool {
	return piece|1 == WhiteRook
}

func isQueen(piece uint8) bool {
	return piece|1 == WhiteQueen
}

func isKing(piece uint8) bool {
	return piece|1 == WhiteKing
}
