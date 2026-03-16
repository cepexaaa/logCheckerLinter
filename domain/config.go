package domain

type Config struct {
	CheckLowercase    bool     `json:"check_lowercase"`
	CheckEnglish      bool     `json:"check_english"`
	CheckNoSpecials   bool     `json:"check_no_specials"`
	CheckSensitive    bool     `json:"check_sensitive"`
	SensitivePatterns []string `json:"sensitive_patterns"`
	AutoFix           bool     `json:"auto_fix"`
}

func DefaultConfig() *Config {
	return &Config{
		CheckLowercase:    true,
		CheckEnglish:      true,
		CheckNoSpecials:   true,
		CheckSensitive:    true,
		SensitivePatterns: []string{"password", "token", "secret", "key", "api_key", "credential"},
		AutoFix:           false,
	}
}
