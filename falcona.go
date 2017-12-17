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
	board.initStandard()
	board.pos[0].print()

}
