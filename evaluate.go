package falcona

func (board *Board) evaluatePieceType(piece uint8) int {
	pos := &board.pos[board.ply]
	score := 0
	var sq int
	bb := pos.pieces[piece]
	for bb != 0 {
		bb, sq = pop(bb)
		if color(piece) == White {
			score += (pieceValue[piece] + pieceTables[piece>>1][sq])
		} else {
			score -= (pieceValue[piece] + pieceTables[piece>>1][mirror[sq]])
		}
	}
	return score
}

func (board *Board) evaluate() int {
	score := 0
	for p := 0; p < 12; p++ {
		score += board.evaluatePieceType(uint8(p))
	}
	if board.pos[board.ply].side == White {
		return score
	}
	return -score
}

var pieceValue = [12]int{
	100, 100, 408, 408, 418, 418, 635, 635, 1260, 1260, 0, 0,
}

var pieceTables = [6][64]int{
	{
		0, 0, 0, 0, 0, 0, 0, 0,
		-6, 8, -4, -2, -2, -4, 8, -6,
		-7, -7, -5, -3, -3, -5, -7, -7,
		-7, 0, -1, 9, 9, -1, 0, -7,
		-13, -7, 8, 16, 16, 8, -7, -13,
		-13, -4, 10, 12, 12, 10, -4, -13,
		-10, 1, 4, 2, 2, 4, 1, -10,
		0, 0, 0, 0, 0, 0, 0, 0,
	}, {
		-98, -33, -21, -15, -15, -21, -33, -98,
		-31, -9, 3, 7, 7, 3, -9, -31,
		-6, 19, 28, 36, 36, 28, 19, -6,
		-13, 8, 19, 25, 25, 19, 8, -13,
		-13, 9, 22, 24, 24, 22, 9, -13,
		-36, -11, 0, 5, 5, 0, -11, -36,
		-42, -22, -11, -5, -5, -11, -22, -42,
		-72, -48, -40, -37, -37, -40, -48, -72,
	}, {
		-23, -11, -15, -20, -20, -15, -11, -23,
		-17, 4, -2, -6, -6, -2, 4, -17,
		-14, 3, 1, -4, -4, 1, 3, -14,
		-11, 7, 3, -1, -1, 3, 7, -11,
		-11, 9, 6, 0, 0, 6, 9, -11,
		-10, 9, 6, 1, 1, 6, 9, -10,
		-15, 5, 1, -5, -5, 1, 5, -15,
		-27, -12, -18, -22, -22, -18, -12, -27,
	}, {
		-12, -8, -6, -3, -3, -6, -8, -12,
		-6, 2, 4, 6, 6, 4, 2, -6,
		-11, -4, 0, 1, 1, 0, -4, -11,
		-11, -4, 0, 1, 1, 0, -4, -11,
		-11, -3, -1, 1, 1, -1, -3, -11,
		-11, -5, -2, 1, 1, -2, -5, -11,
		-11, -4, -2, 0, 0, -2, -4, -11,
		-13, -8, -8, -5, -5, -8, -8, -13,
	}, {
		-1, -2, -1, 0, 0, -1, -2, -1,
		-1, 4, 4, 3, 3, 4, 4, -1,
		-1, 3, 4, 5, 5, 4, 3, -1,
		-2, 5, 4, 4, 4, 4, 5, -2,
		-1, 4, 5, 4, 4, 5, 4, -1,
		-1, 3, 5, 5, 5, 5, 3, -1,
		-2, 3, 5, 4, 4, 5, 3, -2,
		0, -2, -2, -1, -1, -2, -2, 0,
	}, {
		47, 61, 39, 16, 16, 39, 61, 47,
		59, 80, 47, 24, 24, 47, 80, 59,
		74, 95, 57, 35, 35, 57, 95, 74,
		89, 104, 72, 47, 47, 72, 104, 89,
		103, 107, 88, 69, 69, 88, 107, 103,
		114, 137, 102, 69, 69, 102, 137, 114,
		146, 166, 133, 104, 104, 133, 166, 146,
		147, 174, 148, 111, 111, 148, 174, 147,
	},
}

var mirror = [64]int{
	56, 57, 58, 59, 60, 61, 62, 63,
	48, 49, 50, 51, 52, 53, 54, 55,
	40, 41, 42, 43, 44, 45, 46, 47,
	32, 33, 34, 35, 36, 37, 38, 39,
	24, 25, 26, 27, 28, 29, 30, 31,
	16, 17, 18, 19, 20, 21, 22, 23,
	8, 9, 10, 11, 12, 13, 14, 15,
	0, 1, 2, 3, 4, 5, 6, 7,
}
