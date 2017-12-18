package falcona

import "fmt"

func perftRoot(pos string, depth int) {
	board := Board{}
	board.initFEN(pos)
	moves := 0

	ml := board.pos[board.ply].generateMoves()
	for i := 0; i < ml.count; i++ {
		if !board.makeMove(ml.moves[i].move) {
			continue
		}
		temp := board.perft(depth - 1)
		fmt.Print("moves: ", temp, " - ")
		ml.moves[i].print()
		moves += temp
		board.unmakeMove()
	}

	fmt.Println("DEPTH:", depth, " - MOVES:", moves)
}

func (board *Board) perft(depth int) int {
	moves := 0

	ml := board.pos[board.ply].generateMoves()

	for i := 0; i < ml.count; i++ {
		if !board.makeMove(ml.moves[i].move) {
			continue
		}

		board.isRepetition()
		if depth == 1 {
			moves++
			board.unmakeMove()
		} else {
			temp := board.perft(depth - 1)
			board.unmakeMove()
			moves += temp
		}
	}

	return moves
}
