package request

import (
	"strings"
	"unicode"
)

func isUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func checkPath(p string) bool {
	if !strings.HasPrefix(p, "/") {
		return false
	}
	pSlice := strings.Split(p, " ")
	if len(pSlice) > 1 {
		return false
	}
	return true
}
