package rules

import "testing"

func TestCheckEnglish(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		wantOk  bool
		wantErr string
	}{
		{
			name:    "empty message",
			msg:     "",
			wantOk:  true,
			wantErr: "",
		},
		{
			name:    "ascii only",
			msg:     "starting server on port 8080",
			wantOk:  true,
			wantErr: "",
		},
		{
			name:    "russian characters",
			msg:     "запуск сервера",
			wantOk:  false,
			wantErr: "log message must be in English (ASCII only)",
		},
		{
			name:    "chinese characters",
			msg:     "启动服务器",
			wantOk:  false,
			wantErr: "log message must be in English (ASCII only)",
		},
		{
			name:    "mixed ascii and russian",
			msg:     "server запуск",
			wantOk:  false,
			wantErr: "log message must be in English (ASCII only)",
		},
		{
			name:    "french accents",
			msg:     "démarrage du serveur",
			wantOk:  false,
			wantErr: "log message must be in English (ASCII only)",
		},
		{
			name:    "german umlauts",
			msg:     "überprüfung fehlgeschlagen",
			wantOk:  false,
			wantErr: "log message must be in English (ASCII only)",
		},
		{
			name:    "emoji",
			msg:     "server started 🚀",
			wantOk:  false,
			wantErr: "log message must be in English (ASCII only)",
		},
		{
			name:    "ascii with numbers",
			msg:     "response time 123ms",
			wantOk:  true,
			wantErr: "",
		},
		{
			name:    "ascii with symbols",
			msg:     "file not found: /etc/config",
			wantOk:  true,
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOk, gotErr := CheckEnglish(tt.msg)

			if gotOk != tt.wantOk {
				t.Errorf("CheckEnglish() gotOk = %v, want %v", gotOk, tt.wantOk)
			}

			if gotErr != tt.wantErr {
				t.Errorf("CheckEnglish() gotErr = %q, want %q", gotErr, tt.wantErr)
			}
		})
	}
}
