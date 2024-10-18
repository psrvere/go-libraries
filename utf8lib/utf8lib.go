package utf8lib

import (
	"fmt"
	"unicode/utf8"
)

func MaxRune() {
	fmt.Printf("%q\n", utf8.MaxRune)
}

func RuneError() {
	fmt.Printf("%q\n", utf8.RuneError)

	// for empty bytes, DecodeLastRune sends utf8.RuneError
	r, size := utf8.DecodeLastRune([]byte{})
	if r == utf8.RuneError {
		fmt.Printf("rune error found, size: %v\n", size)
	}
}

func AppendRune() {
	fmt.Println(string(utf8.AppendRune(nil, 0x10000)))
	fmt.Println(string(utf8.AppendRune([]byte("initial"), 0x10000)))
}

func DecodeLastRune() {
	b := []byte("Hello, 世界")

	for len(b) > 0 {
		r, size := utf8.DecodeLastRune(b)
		fmt.Printf("%c\t%v\n", r, size)
		b = b[:len(b)-size]
	}
}

func DecodeLastRuneInString() {
	str := "Hello, 世界"

	for len(str) > 0 {
		r, size := utf8.DecodeLastRuneInString(str)
		fmt.Printf("%c\t%v\n", r, size)

		str = str[:len(str)-size]
	}
}

func DecodeRune() {
	b := []byte("Hello, 世界")

	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		fmt.Printf("%c\t%v\n", r, size)
		b = b[size:]
	}
}

func DecodeRuneInString() {
	str := "Hello, 世界"

	for len(str) > 0 {
		r, size := utf8.DecodeRuneInString(str)
		fmt.Printf("%c\t%v\n", r, size)

		str = str[:len(str)-size]
	}
}

func EncodeRune() {
	r := '世'
	buf := make([]byte, 3)
	n := utf8.EncodeRune(buf, r)
	fmt.Printf("buf: %v, n: %v, bufstring: %v\n", buf, n, string(buf))
}

func OutOfRangeRune() {
	// out of range
	// less than zero
	// greater than 0x10FFF
	runes := []rune{-1, 0x110000, utf8.RuneError}
	for i, c := range runes {
		buf := make([]byte, 3)
		size := utf8.EncodeRune(buf, c)
		// Note %[2]s means format second arg as string
		fmt.Printf("%d: %d %[2]s %d\n", i, buf, size)
	}
}

// utf8.RuneCount returns count of rune
func RuneCount() {
	buf := []byte("Hello, 世界")
	fmt.Println("bytes =", len(buf))
	fmt.Println("runes =", utf8.RuneCount(buf))
}

// utf8.RuneLen returns count of bytes
func RuneLen() {
	buf := "Hello, 世界"
	count := 0
	for _, r := range buf {
		count += utf8.RuneLen(r)
	}
	fmt.Printf("total bytes: %v\n", count)
}
