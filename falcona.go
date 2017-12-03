package falcona

import (
	"fmt"
)

type Magic struct {
	mask  uint64
	magic uint64
}

func Print() {
	fmt.Println("HELLO WORLD")
	InitMasks()
	InitMoves()
	InitKeys()
}
