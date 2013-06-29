package scmgo

import (
	"strings"
)

func RuneIndex(str []rune, r rune) (int) {
	for i, c := range str {
		if c == r { return i }
	}
	return -1
}

func RuneIndexAny(str []rune, s string) (int) {
	for i, c := range str {
		if strings.IndexRune(s, c) != -1 {
			return i
		}
	}
	return -1
}
