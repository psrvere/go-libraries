package writingfile

import (
	"crypto/rand"
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func writeToFile(size int) {
	s, unit := convertSize(size)
	path := fmt.Sprintf("./files/file_%v%s.txt", s, unit)

	b := make([]byte, size)
	_, err := rand.Read(b)
	check(err)

	f, err := os.Create(path)
	check(err)
	defer f.Close()

	_, err = f.Write(b)
	check(err)
}

func convertSize(size int) (int, string) {
	switch {
	case size < 1000:
		return size, "bytes"
	case size >= 1000 && size < 1000_1000:
		return size / 1000, "kB"
	case size >= 1000_000 && size < 1000_000_000:
		fmt.Println(size, size/1000_000)
		return size / 1000_000, "MB"
	case size >= 1000_000_000:
		return size / 1000_000_000, "GB"
	default:
		panic("wront input")
	}
}
