package rules

func CheckEnglish(msg string) (ok bool, errMsg string) {
	for _, r := range msg {
		if r > 127 {
			return false, "log message must be in English (ASCII only)"
		}
	}
	return true, ""
}
