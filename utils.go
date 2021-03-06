package falcona

import (
	"fmt"
	"strconv"
)

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

func (m *Move) print() {
	pieceChars := []string{"p", "P", "n", "N", "b", "B", "r", "R", "q", "Q", "k", "K"}

	fmt.Print(pieceChars[moved(m.move)]+": ", from(m.move), " -> ", to(m.move))

	if captured(m.move) != 0xF {
		fmt.Print(" (captured " + pieceChars[captured(m.move)] + ")")
	}
	if promoted(m.move) != 0xF {
		fmt.Print(" (promoted " + pieceChars[promoted(m.move)] + ")")
	}
	if (m.move & isCastle) != 0 {
		fmt.Print(" CASTLE")
	}
	if (m.move & isEnpassant) != 0 {
		fmt.Print(" EP")
	}
	if (m.move & isPawnstart) != 0 {
		fmt.Print(" PS")
	}
	fmt.Println()
}

func (ml *MoveList) print() {
	for i := 0; i < ml.count; i++ {
		fmt.Print(i+1, " - ")
		ml.moves[i].print()
	}
}

func (pos *Position) print() {
	pieceChars := []string{"p", "P", "n", "N", "b", "B", "r", "R", "q", "Q", "k", "K", "", "", "", "\u22C5"}
	for row := 7; row >= 0; row-- {
		fmt.Println()
		fmt.Print(1 + row)
		for col := 0; col <= 7; col++ {
			piece := pos.findPiece(uint(toSquare(row, col)))
			fmt.Print(" " + pieceChars[piece])
		}
	}
	fmt.Println()
	fmt.Println("  a b c d e f g h")
	fmt.Println()

	if pos.side == White {
		fmt.Println("side:    WHITE")
	} else {
		fmt.Println("side:    BLACK")
	}
	fmt.Print("castles: ")
	if pos.castles&1 != 0 {
		fmt.Print("K")
	}
	if pos.castles&2 != 0 {
		fmt.Print("Q")
	}
	if pos.castles&4 != 0 {
		fmt.Print("k")
	}
	if pos.castles&8 != 0 {
		fmt.Print("q")
	}
	fmt.Println()
	fmt.Println("ep:      " + strconv.Itoa(int(pos.enpassant)))
	fmt.Printf("poskey:  %x\n\n", pos.poskey)
}

func (pos *Position) moveFromString(move string) uint32 {
	from := toSquare(int(move[1]-'1'), int(move[0]-'a'))
	to := toSquare(int(move[3]-'1'), int(move[2]-'a'))
	moved := pos.findPiece(uint(from))
	captured := pos.findPiece(uint(to))

	if isKing(moved) && abs(from-to) == 2 {
		return newMove(from, to, moved, 0xF, 0xF, isCastle)
	}

	if isPawn(moved) {
		if abs(from-to) == 16 {
			return newMove(from, to, moved, 0xF, 0xF, isPawnstart)
		} else if abs(from-to) != 8 && captured == 0xF {
			return newMove(from, to, moved, moved^1, 0xF, isEnpassant)
		} else if len(move) > 4 {
			switch move[4] {
			case 'q', 'Q':
				return newMove(from, to, moved, captured, queen(pos.side), 0)
			case 'r', 'R':
				return newMove(from, to, moved, captured, rook(pos.side), 0)
			case 'b', 'B':
				return newMove(from, to, moved, captured, bishop(pos.side), 0)
			case 'n', 'N':
				return newMove(from, to, moved, captured, knight(pos.side), 0)
			}
		}
	}

	return newMove(from, to, moved, captured, 0xF, 0)
}
