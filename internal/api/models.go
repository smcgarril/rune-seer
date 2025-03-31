package api

import "fmt"

var utf8Mask = []string{
	"0xxxxxxx",
	"10xxxxxx",
	"110xxxxx",
	"1110xxxx",
	"11110xxx",
}

type RuneInfo struct {
	Char      string
	RuneIndex int
	RuneBytes []RuneByte
	CodePoint rune
}

type RuneByte struct {
	ByteIndex       int
	ByteInRuneIndex int
	Byte            byte
	Binary          string
	Utf8Mask        string
	Utf8Remainder   string
}

func processStringInput(input string) []RuneInfo {
	var runes []RuneInfo
	n := 0
	for i, c := range input {
		r := RuneInfo{
			Char:      string(c),
			RuneIndex: i,
			CodePoint: c}
		utf8Bytes := []byte(string(c))
		for j, b := range utf8Bytes {
			r.RuneBytes = append(r.RuneBytes, RuneByte{
				ByteIndex:       n,
				ByteInRuneIndex: j,
				Byte:            b,
				Binary:          fmt.Sprintf("%08b", b),
			})
			n++
		}
		runes = append(runes, r)
	}
	return runes
}

func processRune(c rune) RuneInfo {
	var mask string
	var remainder string

	r := RuneInfo{
		Char:      string(c),
		CodePoint: c,
	}

	utf8Bytes := []byte(string(c))
	utf8ByteBinary := fmt.Sprintf("%08b", utf8Bytes[0])

	n := len(utf8Bytes)
	if n > 1 {
		mask = utf8Mask[n]
		remainder = string(utf8ByteBinary[n+1:])
	} else {
		mask = utf8Mask[0]
		remainder = string(utf8ByteBinary[1:])
	}

	for j, b := range utf8Bytes {
		utf8ByteBinary = fmt.Sprintf("%08b", b)
		if j > 0 {
			mask = utf8Mask[1]
			remainder = string(utf8ByteBinary[2:])
		}
		r.RuneBytes = append(r.RuneBytes, RuneByte{
			ByteInRuneIndex: j,
			Byte:            b,
			Binary:          utf8ByteBinary,
			Utf8Mask:        mask,
			Utf8Remainder:   remainder,
		})
	}
	return r
}
