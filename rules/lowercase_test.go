package rules

import (
	"go/token"
	"testing"
)

func TestCheckLowercase(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		wantOk  bool
		wantFix string
		wantErr string
	}{
		{
			name:    "empty message",
			msg:     "",
			wantOk:  true,
			wantFix: "",
			wantErr: "",
		},
		{
			name:    "starts with lowercase letter",
			msg:     "starting server",
			wantOk:  true,
			wantFix: "",
			wantErr: "",
		},
		{
			name:    "starts with uppercase letter",
			msg:     "Starting server",
			wantOk:  false,
			wantFix: "starting server",
			wantErr: "the log message must begin with a lowercase letter",
		},
		{
			name:    "starts with digit",
			msg:     "2 users connected",
			wantOk:  false,
			wantFix: "",
			wantErr: "the log message must begin with a lowercase letter",
		},
		{
			name:    "starts with symbol",
			msg:     "!important message",
			wantOk:  false,
			wantFix: "",
			wantErr: "the log message must begin with a lowercase letter",
		},
		{
			name:    "unicode uppercase",
			msg:     "Привет мир",
			wantOk:  false,
			wantFix: "привет мир",
			wantErr: "the log message must begin with a lowercase letter",
		},
		{
			name:    "unicode lowercase",
			msg:     "привет мир",
			wantOk:  true,
			wantFix: "",
			wantErr: "",
		},
		{
			name:    "single uppercase letter",
			msg:     "A",
			wantOk:  false,
			wantFix: "a",
			wantErr: "the log message must begin with a lowercase letter",
		},
		{
			name:    "single lowercase letter",
			msg:     "a",
			wantOk:  true,
			wantFix: "",
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOk, gotFix, gotErr := CheckLowercase(tt.msg, token.NoPos)

			if gotOk != tt.wantOk {
				t.Errorf("CheckLowercase() gotOk = %v, want %v", gotOk, tt.wantOk)
			}

			if gotFix != tt.wantFix {
				t.Errorf("CheckLowercase() gotFix = %q, want %q", gotFix, tt.wantFix)
			}

			if gotErr != tt.wantErr && !(gotErr == "" && tt.wantErr == "") {
				t.Errorf("CheckLowercase() gotErr = %q, want %q", gotErr, tt.wantErr)
			}
		})
	}
}

func TestCheckLowercaseWithPosition(t *testing.T) {
	pos := token.Pos(100)
	ok, fix, err := CheckLowercase("Test", pos)

	if ok {
		t.Error("CheckLowercase() should return false for uppercase")
	}

	if fix != "test" {
		t.Errorf("CheckLowercase() fix = %q, want %q", fix, "test")
	}

	if err == "" {
		t.Error("CheckLowercase() should return error message")
	}
}
