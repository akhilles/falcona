package falcona

import "math/bits"

func Slowcount(bb uint64) int {
	bb -= (bb >> 1) & 0x5555555555555555
	bb = ((bb >> 2) & 0x3333333333333333) + (bb & 0x3333333333333333)
	bb = ((bb >> 4) + bb) & 0x0F0F0F0F0F0F0F0F
	return int((bb * 0x0101010101010101) >> 56)
}

func Count(bb uint64) int {
	return bits.OnesCount64(bb)
}

func Scanforward(bb uint64) int {
	return bits.LeadingZeros64(bb)
}

func Scanreverse(bb uint64) int {
	return bits.TrailingZeros64(bb)
}

func pop(bb uint64) (uint64, int) {
	return bb & (bb - 1), Scanforward(bb)
}

func set(bb uint64, index int) uint64 {
	return bb | (1 << uint(index))
}

func sample(bb uint64, index int) (mask uint64) {
	count := Count(bb)
	for i := 0; i < count; i++ {
		popped := (bb & (bb - 1)) ^ bb
		bb &= bb - 1
		if (1<<uint(i))&index != 0 {
			mask |= popped
		}
	}
	return mask
}
