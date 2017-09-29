package falcona

func popcount(bb uint64) int {
	bb -= (bb >> 1) & 0x5555555555555555
	bb = ((bb >> 2) & 0x3333333333333333) + (bb & 0x3333333333333333)
	bb = ((bb >> 4) + bb) & 0x0F0F0F0F0F0F0F0F
	return int((bb * 0x0101010101010101) >> 56)
}
