package falcona

// Precomputed values
var (
	kingMoves        [64]uint64
	knightMoves      [64]uint64
	pawnAttacks      [2][64]uint64
	rookMagicMoves   [64][4096]uint64
	bishopMagicMoves [64][512]uint64

	maskRay   [64][64]uint64
	maskBlock [64][64]uint64

	keyCastleSimple [16]uint64
)

func slidingMoves(square int, blockers uint64, dirs [4]int, stop [4]uint64) uint64 {
	var moves uint64
	for index := 0; index < 4; index++ {
		for sq := square + dirs[index]; sq >= A1 && sq <= H8 && maskSquare[sq]&stop[index] == 0; sq += dirs[index] {
			moves = set(moves, sq)
			if maskSquare[sq]&blockers != 0 {
				break
			}
		}
	}
	return moves
}

func InitKeys() {
	for castleCombo := uint8(0); castleCombo < 16; castleCombo++ {
		if castleCombo&castleKingside[White] != 0 {
			keyCastleSimple[castleCombo] ^= keyCastle[0]
		}
		if castleCombo&castleQueenside[White] != 0 {
			keyCastleSimple[castleCombo] ^= keyCastle[1]
		}
		if castleCombo&castleKingside[Black] != 0 {
			keyCastleSimple[castleCombo] ^= keyCastle[2]
		}
		if castleCombo&castleQueenside[Black] != 0 {
			keyCastleSimple[castleCombo] ^= keyCastle[3]
		}
		//fmt.Printf("%016X\n", keyCastleSimple[castleCombo])
	}
}

func InitMasks() {
	for i := A1; i <= H8; i++ {
		r, c := toRowCol(i)

		for j := A1; j <= H8; j++ {
			r1, c1 := toRowCol(j)
			var blockZone uint64
			if r == r1 {
				maskRay[i][j] = maskRank[r]
				for k := c + 1; k <= c1; k++ {
					blockZone |= maskFile[k]
				}
				for k := c1; k < c; k++ {
					blockZone |= maskFile[k]
				}
			} else if c == c1 {
				maskRay[i][j] = maskFile[c]
				for k := r + 1; k <= r1; k++ {
					blockZone |= maskRank[k]
				}
				for k := r1; k < r; k++ {
					blockZone |= maskRank[k]
				}
			} else if r1+c == r+c1 {
				maskRay[i][j] = maskDiagRight[7-r+c]
				for k := c + 1; k <= c1; k++ {
					blockZone |= maskFile[k]
				}
				for k := c1; k < c; k++ {
					blockZone |= maskFile[k]
				}
			} else if r+c == r1+c1 {
				maskRay[i][j] = maskDiagLeft[r+c]
				for k := c + 1; k <= c1; k++ {
					blockZone |= maskFile[k]
				}
				for k := c1; k < c; k++ {
					blockZone |= maskFile[k]
				}
			}
			maskBlock[i][j] = maskRay[i][j] & blockZone
			//fmt.Printf("%2d %2d %016X %016X\n", i, j, maskRay[i][j], maskBlock[i][j])
		}
	}
}

func InitMoves() {
	for i := A1; i <= H8; i++ {
		r, c := toRowCol(i)

		toTrim := [4]uint64{
			maskFile[FA], maskFile[FH], maskRank[R1], maskRank[R8],
		}
		// Rook moves
		rookBlockers := rookMagic[i].mask
		for _, maskElement := range toTrim {
			if count(rookBlockers&maskElement) <= 1 {
				rookBlockers &= ^maskElement
			}
		}
		numCombinations := 1 << uint(count(rookBlockers))
		for j := 0; j < numCombinations; j++ {
			mask := sample(rookBlockers, j)
			index := (mask * rookMagic[i].magic) >> 52
			rookMagicMoves[i][index] = slidingMoves(i, mask, [4]int{1, -1, 8, -8}, [4]uint64{maskFile[FA], maskFile[FH], 0, 0})
		}

		// Bishop moves
		bishopBlockers := bishopMagic[i].mask
		for _, maskElement := range toTrim {
			bishopBlockers &= ^maskElement
		}
		numCombinations = 1 << uint(count(bishopBlockers))
		for j := 0; j < numCombinations; j++ {
			mask := sample(bishopBlockers, j)
			index := (mask * bishopMagic[i].magic) >> 55
			bishopMagicMoves[i][index] = slidingMoves(i, mask, [4]int{9, 7, -7, -9}, [4]uint64{maskFile[FA], maskFile[FH], maskFile[FA], maskFile[FH]})
		}

		// Pawn attacks
		if r > R1 && r < R8 {
			if c > FA {
				pawnAttacks[White][i] = set(pawnAttacks[White][i], i+7)
				pawnAttacks[Black][i] = set(pawnAttacks[Black][i], i-9)
			}
			if c < FH {
				pawnAttacks[White][i] = set(pawnAttacks[White][i], i+9)
				pawnAttacks[Black][i] = set(pawnAttacks[Black][i], i-7)
			}
		}
		//fmt.Printf("%2d %016X %016X\n", i, pawnAttacks[White][i], pawnAttacks[Black][i])

		// Knight and king moves
		for j := A1; j <= H8; j++ {
			r1, c1 := toRowCol(j)

			if i == j {
				continue
			}
			if (abs(r-r1) == 1 && abs(c-c1) == 2) || (abs(r-r1) == 2 && abs(c-c1) == 1) {
				knightMoves[i] = set(knightMoves[i], j)
			}
			if abs(r-r1) < 2 && abs(c-c1) < 2 {
				kingMoves[i] = set(kingMoves[i], j)
			}
		}
		//print(knightMoves[i])
	}
}
