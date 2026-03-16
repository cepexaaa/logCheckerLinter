package rules

import (
	"go/token"
	"unicode"
)

func CheckLowercase(msg string, pos token.Pos) (ok bool, fix string, errMsg string) {
	if len(msg) == 0 {
		return true, "", ""
	}
	r := []rune(msg)[0]
	if unicode.IsLetter(r) && unicode.IsUpper(r) {
		fix = string(unicode.ToLower(r)) + msg[len(string(r)):]
		return false, fix, "the log message must begin with a lowercase letter"
	}
	if !unicode.IsLetter(r) {
		return false, "", "the log message must begin with a lowercase letter"
	}
	return true, "", ""
}
