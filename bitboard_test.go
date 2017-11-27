package falcona

import (
	"fmt"
	"testing"
)

func BenchmarkPop(b *testing.B) {
	var bb uint64
	var ind int
	var news uint64
	bb = 0
	for n := 0; n < b.N; n++ {
		bb++
		ind, news = Pop(bb)
	}
	if ind == 0 {
		if news == 0 {
			fmt.Println("WEHUJ")
		}
	}
}

func BenchmarkPopSlow(b *testing.B) {
	var bb uint64
	bb = 0
	for n := 0; n < b.N; n++ {
		bb++
		_ = PopSlow(&bb)
	}
}
