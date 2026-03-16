package rules

import (
	"regexp"
	"strings"
)

type SensitiveChecker struct {
	patterns      []*regexp.Regexp
	keywords      []string
	checkVarNames bool
}

func NewSensitiveChecker(keywords []string) *SensitiveChecker {
	var patterns []*regexp.Regexp

	for _, kw := range keywords {
		re := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(kw) + `\s*[:=]`)
		patterns = append(patterns, re)
	}

	for _, kw := range keywords {
		re := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(kw) + `\s+[^\s]{3,}`)
		patterns = append(patterns, re)
	}

	emailPattern := regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`)
	patterns = append(patterns, emailPattern)

	creditPattern := regexp.MustCompile(`\b(?:\d[ -]*?){13,16}\b`)
	patterns = append(patterns, creditPattern)

	return &SensitiveChecker{
		patterns:      patterns,
		keywords:      keywords,
		checkVarNames: true,
	}
}

type SensitiveConfig struct {
	Keywords         []string `json:"keywords"`
	CheckEmail       bool     `json:"check_email"`
	CheckCreditCards bool     `json:"check_credit_cards"`
	CheckVarNames    bool     `json:"check_var_names"`
	CustomPatterns   []string `json:"custom_patterns"`
}

func NewSensitiveCheckerFromConfig(cfg SensitiveConfig) *SensitiveChecker {
	var patterns []*regexp.Regexp

	for _, kw := range cfg.Keywords {
		patterns = append(patterns,
			regexp.MustCompile(`(?i)\b`+regexp.QuoteMeta(kw)+`\s*[:=]`),
			regexp.MustCompile(`(?i)\b`+regexp.QuoteMeta(kw)+`\s+[^\s]{3,}`),
		)
	}

	if cfg.CheckEmail {
		patterns = append(patterns, regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`))
	}

	if cfg.CheckCreditCards {
		patterns = append(patterns, regexp.MustCompile(`\b(?:\d[ -]*?){13,16}\b`))
	}

	for _, custom := range cfg.CustomPatterns {
		patterns = append(patterns, regexp.MustCompile(custom))
	}

	return &SensitiveChecker{
		patterns:      patterns,
		keywords:      cfg.Keywords,
		checkVarNames: cfg.CheckVarNames,
	}
}

func (sc *SensitiveChecker) CheckSensitive(msg string) (ok bool, errMsg string) {
	for _, re := range sc.patterns {
		if re.MatchString(msg) {
			return false, "log message must not contain potentially sensitive data"
		}
	}
	return true, ""
}

func (sc *SensitiveChecker) CheckSensitiveConcat(parts []string, varNames []string) (ok bool, errMsg string) {

	for _, part := range parts {
		if ok, _ := sc.CheckSensitive(part); !ok {
			return false, "the log message contains sensitive data"
		}
	}

	if sc.checkVarNames {
		for _, varName := range varNames {
			for _, kw := range sc.keywords {
				if strings.Contains(strings.ToLower(varName), strings.ToLower(kw)) {
					return false, "log message must not contain potentially sensitive data"
				}
			}
		}
	}

	return true, ""
}
