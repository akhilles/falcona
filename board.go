package falcona

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type Position struct {
	pieces [12]uint64
	colors [2]uint64
	kings  [2]uint8

	side      uint8
	enpassant uint8
	castles   uint8
	fiftymove uint8

	poskey uint64
}

type Board struct {
	pos [MaxMoves]Position

	ply    int
	hisply int
}

func (board *Board) initFEN(fen string) {
	board.hisply = 0
	board.ply = 0

	board.pos[board.hisply] = Position{}
	pos := &board.pos[board.hisply]

	substrings := strings.Split(fen, " ")

	sq := A8
	for _, char := range substrings[0] {
		var piece uint8 = 12
		switch char {
		case 'P':
			piece = WhitePawn
		case 'p':
			piece = BlackPawn
		case 'N':
			piece = WhiteKnight
		case 'n':
			piece = BlackKnight
		case 'B':
			piece = WhiteBishop
		case 'b':
			piece = BlackBishop
		case 'R':
			piece = WhiteRook
		case 'r':
			piece = BlackRook
		case 'Q':
			piece = WhiteQueen
		case 'q':
			piece = BlackQueen
		case 'K':
			piece = WhiteKing
			pos.kings[White] = uint8(sq)
		case 'k':
			piece = BlackKing
			pos.kings[Black] = uint8(sq)
		case '/':
			sq -= 16
		case '1', '2', '3', '4', '5', '6', '7', '8':
			sq += int(char - '0')
		}
		if piece < 12 {
			pos.pieces[piece] = set(pos.pieces[piece], sq)
			pos.colors[color(piece)] = set(pos.colors[color(piece)], sq)
			sq++
		}
	}

	pos.side = White
	if substrings[1] == "b" {
		pos.side = Black
	}

	for _, char := range substrings[2] {
		switch char {
		case 'K':
			pos.castles |= castleKingside[White]
		case 'Q':
			pos.castles |= castleQueenside[White]
		case 'k':
			pos.castles |= castleKingside[Black]
		case 'q':
			pos.castles |= castleQueenside[Black]
		}
	}

	pos.enpassant = 64
	if substrings[3] != "-" {
		r := int(substrings[3][1] - '1')
		c := int(substrings[3][0] - 'a')
		pos.enpassant = uint8(toSquare(r, c))
	}

	if n, err := strconv.Atoi(substrings[4]); err == nil {
		pos.fiftymove = uint8(n)
	}

	pos.poskey = getPoskey(pos)
}

func (pos *Position) findPiece(sq uint) uint8 {
	var squareMask uint64 = 1 << sq
	for piece, bb := range pos.pieces {
		if bb&squareMask != 0 {
			return uint8(piece)
		}
	}
	return 12
}

func (pos *Position) print() {
	pieceChars := []string{"p", "P", "k", "K", "b", "B", "r", "R", "q", "Q", "k", "K", "\u22C5"}
	buffer := bytes.NewBufferString("")
	for row := 7; row >= 0; row-- {
		buffer.WriteByte('1' + byte(row))
		for col := 0; col <= 7; col++ {
			piece := pos.findPiece(uint(toSquare(row, col)))
			buffer.WriteByte(' ')
			buffer.WriteString(pieceChars[piece])
		}
		buffer.WriteByte('\n')
	}
	buffer.WriteString("  a b c d e f g h  \n")

	if pos.side == White {
		buffer.WriteString("WHITE to move\n")
	} else {
		buffer.WriteString("BLACK to move\n")
	}

	if pos.castles&1 != 0 {
		buffer.WriteByte('K')
	}
	if pos.castles&2 != 0 {
		buffer.WriteByte('Q')
	}
	if pos.castles&4 != 0 {
		buffer.WriteByte('k')
	}
	if pos.castles&8 != 0 {
		buffer.WriteByte('q')
	}
	buffer.WriteString(", ep:" + strconv.Itoa(int(pos.enpassant)) + "\n")

	fmt.Println(buffer.String())
}