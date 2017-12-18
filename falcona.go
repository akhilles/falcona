package falcona

type Magic struct {
	mask  uint64
	magic uint64
}

func Print() {
	InitMasks()
	InitMoves()
	InitKeys()

	//perftRoot("n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1", 6)
	//perftRoot("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 3)

	uci()
}
