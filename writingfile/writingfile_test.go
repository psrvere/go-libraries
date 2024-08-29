package writingfile

import (
	"fmt"
	"testing"
)

func TestWriteToFile(t *testing.T) {
	// writeToFile(100) // 100 bytes
	// writeToFile(1000_00) // 100 kb
	writeToFile(1000_000_00) // 100 MB
	// writeToFile(1000_000_000) // 1 GB
	// writeToFile(1000_000_000_0) // 10 GB
}

func BenchmarkWriteToFile(b *testing.B) {
	cases := []int{
		100,
		1000_00,
		1000_000_00,
		1000_000_000,
		1000_000_000_0,
	}

	for _, size := range cases {
		s, unit := convertSize(size)
		b.Run(fmt.Sprintf("Size: %d%s", s, unit), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				writeToFile(size)
			}
		})
	}
}
