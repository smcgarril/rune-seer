package api

import (
	"testing"
)

func TestProcessStringInput(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output []RuneInfo
	}{
		{
			name:  "Basic ASCII",
			input: "Hello",
			output: []RuneInfo{
				{"H", 0, []RuneByte{{0, 0, 'H', "01001000", "110xxxxx", "00000000"}}, 'H'},
				{"e", 1, []RuneByte{{1, 0, 'e', "01100101", "110xxxxx", "00000000"}}, 'e'},
				{"l", 2, []RuneByte{{2, 0, 'l', "01101100", "110xxxxx", "00000000"}}, 'l'},
				{"l", 3, []RuneByte{{3, 0, 'l', "01101100", "110xxxxx", "00000000"}}, 'l'},
				{"o", 4, []RuneByte{{4, 0, 'o', "01101111", "110xxxxx", "00000000"}}, 'o'},
			},
		},
		{
			name:  "Mixed ASCII and Unicode",
			input: "a😊b",
			output: []RuneInfo{
				{"a", 0, []RuneByte{{0, 0, 'a', "01100001", "110xxxxx", "00000000"}}, 'a'},
				{"😊", 1, []RuneByte{
					{1, 0, 0xF0, "11110000", "11110xxx", "00000000"},
					{2, 1, 0x9F, "10011111", "10xxxxxx", "00000000"},
					{3, 2, 0x98, "10011000", "10xxxxxx", "00000000"},
					{4, 3, 0x8A, "10001010", "10xxxxxx", "00000000"},
				}, '😊'},
				{"b", 2, []RuneByte{{5, 0, 'b', "01100010", "110xxxxx", "00000000"}}, 'b'},
			},
		},
		{
			name:   "Empty String",
			input:  "",
			output: []RuneInfo{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := processStringInput(tc.input)
			if len(result) != len(tc.output) {
				t.Fatalf("Expected length %d, got %d", len(tc.output), len(result))
			}

			for i, expected := range tc.output {
				got := result[i]

				if expected.Char != got.Char {
					t.Fatalf("Expected char %s, got %s", expected.Char, got.Char)
				}

				if len(expected.RuneBytes) != len(got.RuneBytes) {
					t.Fatalf("Expected %d bytes, got %d", len(expected.RuneBytes), len(got.RuneBytes))
				}

				for j, expectedByte := range expected.RuneBytes {
					gotByte := got.RuneBytes[j]
					if expectedByte.ByteIndex != gotByte.ByteIndex || expectedByte.ByteInRuneIndex != gotByte.ByteInRuneIndex ||
						expectedByte.Byte != gotByte.Byte || expectedByte.Binary != gotByte.Binary {
						t.Fatalf("Mismatch in RuneByte at index %d: expected %v, got %v", j, expectedByte, gotByte)
					}
				}

				if expected.CodePoint != got.CodePoint {
					t.Fatalf("Expected code point %U, got %U", expected.CodePoint, got.CodePoint)
				}
			}
		})
	}
}
