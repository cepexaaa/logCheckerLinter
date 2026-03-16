package rules

import (
	"testing"
)

func TestSensitiveChecker_CheckSensitive(t *testing.T) {
	defaultKeywords := []string{"password", "token", "secret", "key", "api_key", "credential"}

	tests := []struct {
		name     string
		keywords []string
		msg      string
		wantOk   bool
		wantErr  string
	}{
		{
			name:     "empty message",
			keywords: defaultKeywords,
			msg:      "",
			wantOk:   true,
			wantErr:  "",
		},
		{
			name:     "no sensitive data",
			keywords: defaultKeywords,
			msg:      "user logged in successfully",
			wantOk:   true,
			wantErr:  "",
		},
		{
			name:     "password with colon",
			keywords: defaultKeywords,
			msg:      "user password: 12345",
			wantOk:   false,
			wantErr:  "log message must not contain potentially sensitive data",
		},
		{
			name:     "password with equals",
			keywords: defaultKeywords,
			msg:      "password=12345",
			wantOk:   false,
			wantErr:  "log message must not contain potentially sensitive data",
		},
		{
			name:     "token with colon",
			keywords: defaultKeywords,
			msg:      "auth token: abc123",
			wantOk:   false,
			wantErr:  "log message must not contain potentially sensitive data",
		},
		{
			name:     "api key with equals",
			keywords: defaultKeywords,
			msg:      "api_key=sk-123456",
			wantOk:   false,
			wantErr:  "log message must not contain potentially sensitive data",
		},
		{
			name:     "case insensitive",
			keywords: defaultKeywords,
			msg:      "USER PASSWORD: 123",
			wantOk:   false,
			wantErr:  "log message must not contain potentially sensitive data",
		},
		{
			name:     "word boundary",
			keywords: defaultKeywords,
			msg:      "password123: value",
			wantOk:   true,
			wantErr:  "",
		},
		{
			name:     "custom keywords",
			keywords: []string{"secret", "private"},
			msg:      "private data: 123",
			wantOk:   false,
			wantErr:  "log message must not contain potentially sensitive data",
		},
		{
			name:     "multiple sensitive words",
			keywords: defaultKeywords,
			msg:      "token=abc and password=123",
			wantOk:   false,
			wantErr:  "log message must not contain potentially sensitive data",
		},
		{
			name:     "whitespace handling",
			keywords: defaultKeywords,
			msg:      "password:123",
			wantOk:   false,
			wantErr:  "log message must not contain potentially sensitive data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checker := NewSensitiveChecker(tt.keywords)
			gotOk, gotErr := checker.CheckSensitive(tt.msg)

			if gotOk != tt.wantOk {
				t.Errorf("CheckSensitive() gotOk = %v, want %v", gotOk, tt.wantOk)
			}

			if gotErr != tt.wantErr {
				t.Errorf("CheckSensitive() gotErr = %q, want %q", gotErr, tt.wantErr)
			}
		})
	}
}

func TestSensitiveChecker_CheckSensitiveConcat(t *testing.T) {
	defaultKeywords := []string{"password", "token", "secret"}

	tests := []struct {
		name     string
		keywords []string
		strParts []string
		varNames []string
		wantOk   bool
		wantErr  string
	}{
		{
			name:     "no sensitive parts",
			keywords: defaultKeywords,
			strParts: []string{"user ", " logged in"},
			varNames: []string{"username"},
			wantOk:   true,
			wantErr:  "",
		},
		{
			name:     "sensitive in string part",
			keywords: defaultKeywords,
			strParts: []string{"password: ", "123"},
			varNames: []string{"pwd"},
			wantOk:   false,
			wantErr:  "the log message contains sensitive data",
		},
		{
			name:     "sensitive variable name",
			keywords: defaultKeywords,
			strParts: []string{"user ", " value"},
			varNames: []string{"password"},
			wantOk:   false,
			wantErr:  "log message must not contain potentially sensitive data",
		},
		{
			name:     "multiple variables",
			keywords: defaultKeywords,
			strParts: []string{"data: "},
			varNames: []string{"token", "username"},
			wantOk:   false,
			wantErr:  "log message must not contain potentially sensitive data",
		},
		{
			name:     "sensitive in both",
			keywords: defaultKeywords,
			strParts: []string{"secret: ", "value"},
			varNames: []string{"password"},
			wantOk:   false,
			wantErr:  "the log message contains sensitive data",
		},
		{
			name:     "partial match variable",
			keywords: defaultKeywords,
			strParts: []string{"data"},
			varNames: []string{"userPassword"},
			wantOk:   false,
			wantErr:  "log message must not contain potentially sensitive data",
		},
		{
			name:     "safe variable name",
			keywords: defaultKeywords,
			strParts: []string{"count: "},
			varNames: []string{"attempts"},
			wantOk:   true,
			wantErr:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checker := NewSensitiveChecker(tt.keywords)
			checker.checkVarNames = true

			gotOk, gotErr := checker.CheckSensitiveConcat(tt.strParts, tt.varNames)

			if gotOk != tt.wantOk {
				t.Errorf("CheckSensitiveConcat() gotOk = %v, want %v", gotOk, tt.wantOk)
			}

			if gotErr != tt.wantErr {
				t.Errorf("CheckSensitiveConcat() gotErr = %q, want %q", gotErr, tt.wantErr)
			}
		})
	}
}

func TestSensitiveChecker_WithCustomPatterns(t *testing.T) {
	checker := NewSensitiveChecker([]string{"password"})

	t.Run("email detection", func(t *testing.T) {
		msg := "user email: test@example.com"
		ok, err := checker.CheckSensitive(msg)

		if ok {
			t.Error("CheckSensitive() should detect email")
		}
		if err == "" {
			t.Error("CheckSensitive() should return error for email")
		}
	})

	t.Run("credit card detection", func(t *testing.T) {
		msg := "card: 4111-1111-1111-1111"
		ok, err := checker.CheckSensitive(msg)

		if ok {
			t.Error("CheckSensitive() should detect credit card")
		}
		if err == "" {
			t.Error("CheckSensitive() should return error for credit card")
		}
	})

	t.Run("simple number not detected", func(t *testing.T) {
		msg := "count: 42"
		ok, err := checker.CheckSensitive(msg)

		if !ok {
			t.Error("CheckSensitive() should not detect simple number")
		}
		if err != "" {
			t.Errorf("CheckSensitive() returned error: %q", err)
		}
	})
}
