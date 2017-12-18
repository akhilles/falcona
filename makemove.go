package falcona

import "os"
import "fmt"

func (pos *Position) flipPiece(sq uint8, piece uint8) {
	pos.pieces[piece] ^= (1 << sq)
	pos.colors[color(piece)] ^= (1 << sq)
	pos.poskey ^= keyPieces[int(piece)*64+int(sq)]
}

func (pos *Position) movePiece(from, to uint8, piece uint8) {
	pos.pieces[piece] ^= ((1 << from) | (1 << to))
	pos.colors[color(piece)] ^= ((1 << from) | (1 << to))
	pos.poskey ^= keyPieces[int(piece)*64+int(from)] ^ keyPieces[int(piece)*64+int(to)]
}

func (board *Board) makeMove(move uint32) bool {
	from := from(move)
	to := to(move)
	moved := moved(move)
	captured := captured(move)
	promoted := promoted(move)
	ca := (move & isCastle) != 0
	ps := (move & isPawnstart) != 0
	ep := (move & isEnpassant) != 0

	board.ply++
	board.pos[board.ply] = board.pos[board.ply-1]
	posold := &board.pos[board.ply-1]
	pos := &board.pos[board.ply]

	pos.fiftymove++
	if pos.enpassant != 64 {
		pos.poskey ^= keyEnpassant[pos.enpassant&7]
		pos.enpassant = 64
	}

	pos.movePiece(from, to, moved)

	if promoted != 0xF {
		pos.fiftymove = 0
		pos.flipPiece(to, moved)
		pos.flipPiece(to, promoted)
	}

	if isKing(moved) {
		pos.kings[pos.side] = to
	}

	if captured != 0xF {
		pos.fiftymove = 0
		pos.flipPiece(to, captured)
		if ep {
			pos.flipPiece(to, captured)
			if pos.side == White {
				pos.flipPiece(to-8, captured)
			} else {
				pos.flipPiece(to+8, captured)
			}
		}
	} else if ca {
		switch to {
		case G1:
			pos.movePiece(H1, F1, WhiteRook)
		case C1:
			pos.movePiece(A1, D1, WhiteRook)
		case G8:
			pos.movePiece(H8, F8, BlackRook)
		case C8:
			pos.movePiece(A8, D8, BlackRook)
		}
	} else if ps {
		pos.fiftymove = 0
		pos.enpassant = (from + to) / 2
		pos.poskey ^= keyEnpassant[pos.enpassant&7]
	} else if isPawn(moved) {
		pos.fiftymove = 0
	}

	pos.occ = pos.colors[White] | pos.colors[Black]
	pos.castles &= castleRights[from] & castleRights[to]
	pos.poskey ^= keyCastleSimple[posold.castles] ^ keyCastleSimple[pos.castles]
	pos.poskey ^= keyColor[White]
	pos.side ^= 1

	if pos.attackedBy(int(pos.kings[pos.side^1]), pos.side) {
		board.ply--
		return false
	}

	if pos.poskey != pos.getPoskey() {
		fmt.Println(moveToString(move), move)
		pos.print()
		os.Exit(3)
	}

	return true
}

func (board *Board) unmakeMove() {
	board.ply--
}
