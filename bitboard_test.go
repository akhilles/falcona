package falcona

import "testing"

func BenchmarkSlowcount(b *testing.B) {
	var bb uint64
	bb = 0

	for n := 0; n < b.N; n++ {
		bb++
		Slowcount(bb)
	}
}

func BenchmarkFastcount(b *testing.B) {
	var bb uint64
	bb = 0
	for n := 0; n < b.N; n++ {
		bb++
		Fastcount(bb)
	}
}
