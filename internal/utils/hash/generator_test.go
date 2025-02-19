package hash

import "testing"

func BenchmarkGenerator(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Generator(10000)
	}
}
