package main

import (
	"fmt"

	"github.com/akhilles/falcona"
)

func main() {
	var vary uint64
	var index int
	vary = 3434
	fmt.Println(vary)
	vary, index = falcona.Pop(vary)
	fmt.Println(vary, index)
}
