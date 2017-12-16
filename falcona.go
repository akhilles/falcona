package falcona

type Magic struct {
	mask  uint64
	magic uint64
}

func Print() {
	InitMasks()
	InitMoves()
	InitKeys()

	board := Board{}
	board.initFEN("rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2")
	board.pos[0].print()
}
