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
	CodePoint string
}

type RuneByte struct {
	ByteIndex     int
	RuneByteIndex int
	Byte          byte
	Binary        string
	Utf8Mask      string
	Utf8Remainder string
}

func ProcessInput(input string) []RuneInfo {
	var runes []RuneInfo
	n := 0
	for i, r := range input {
		rune := RuneInfo{
			Char:      string(r),
			RuneIndex: i,
			CodePoint: fmt.Sprintf("%d", r)}
		utf8Bytes := []byte(string(r))
		for j, b := range utf8Bytes {
			rune.RuneBytes = append(rune.RuneBytes, RuneByte{
				ByteIndex:     n,
				RuneByteIndex: j,
				Byte:          b,
				Binary:        fmt.Sprintf("%08b", b),
				Utf8Mask:      utf8Mask[j],
			})
			n++
		}
		runes = append(runes, rune)
	}
	return runes
}

func ProcessRune(r rune) RuneInfo {
	var mask string
	var remainder string

	rune := RuneInfo{
		Char:      string(r),
		CodePoint: fmt.Sprintf("%d", r),
	}

	utf8Bytes := []byte(string(r))
	utf8ByteBinary := fmt.Sprintf("%08b", utf8Bytes[0])

	n := len(utf8Bytes)
	if n > 1 {
		mask = utf8Mask[n]
		remainder = string(utf8ByteBinary[n+1:])
	} else {
		mask = utf8Mask[0]
		remainder = string(utf8ByteBinary[1:])
	}

	fmt.Println(remainder)

	for j, b := range utf8Bytes {
		utf8ByteBinary = fmt.Sprintf("%08b", b)
		if j > 0 {
			mask = utf8Mask[1]
			remainder = string(utf8ByteBinary[2:])
		}
		rune.RuneBytes = append(rune.RuneBytes, RuneByte{
			RuneByteIndex: j,
			Byte:          b,
			Binary:        utf8ByteBinary,
			Utf8Mask:      mask,
			Utf8Remainder: remainder,
		})
	}
	return rune
}
