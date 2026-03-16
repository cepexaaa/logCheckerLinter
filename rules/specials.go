package rules

import "regexp"

var repeatedPunctuation = regexp.MustCompile(`[!?.]{2,}`)

func CheckNoSpecials(msg string) (ok bool, errMsg string) {
	if repeatedPunctuation.MatchString(msg) {
		return false, "log message must not contain repeated punctuation marks"
	}
	// Emojis and other non-ASCII characters are detected by the CheckEnglish rule
	return true, ""
}
