package falcona

import (
	"bytes"
	"fmt"
	"math/bits"
)

func slowcount(bb uint64) int {
	bb -= (bb >> 1) & 0x5555555555555555
	bb = ((bb >> 2) & 0x3333333333333333) + (bb & 0x3333333333333333)
	bb = ((bb >> 4) + bb) & 0x0F0F0F0F0F0F0F0F
	return int((bb * 0x0101010101010101) >> 56)
}

func count(bb uint64) int {
	return bits.OnesCount64(bb)
}

func scanforward(bb uint64) int {
	return bits.LeadingZeros64(bb)
}

func scanreverse(bb uint64) int {
	return bits.TrailingZeros64(bb)
}

func pop(bb uint64) (uint64, int) {
	return bb & (bb - 1), scanreverse(bb)
}

func set(bb uint64, index int) uint64 {
	return bb | (1 << uint(index))
}

func sample(bb uint64, index int) (mask uint64) {
	count := count(bb)
	for i := 0; i < count; i++ {
		popped := (bb & (bb - 1)) ^ bb
		bb &= bb - 1
		if (1<<uint(i))&index != 0 {
			mask |= popped
		}
	}
	return mask
}

func print(bb uint64) {
	buffer := bytes.NewBufferString("")
	//buffer.WriteString(fmt.Sprintf("0x%016X\n", bb))
	for row := 7; row >= 0; row-- {
		buffer.WriteByte('1' + byte(row))
		for col := 0; col <= 7; col++ {
			offset := row<<3 + col
			buffer.WriteByte(' ')
			if bb&(1<<uint(offset)) != 0 {
				buffer.WriteString("\u2022") // Set
			} else {
				buffer.WriteString("\u22C5") // Clear
			}
		}
		buffer.WriteByte('\n')
	}
	buffer.WriteString("  a b c d e f g h  \n")
	fmt.Println(buffer.String())
}
