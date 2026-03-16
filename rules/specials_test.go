package rules

import "testing"

func TestCheckNoSpecials(t *testing.T) {
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
			name:    "normal message",
			msg:     "server started successfully",
			wantOk:  true,
			wantErr: "",
		},
		{
			name:    "single exclamation",
			msg:     "warning",
			wantOk:  true,
			wantErr: "",
		},
		{
			name:    "double exclamation",
			msg:     "warning!!",
			wantOk:  false,
			wantErr: "log message must not contain repeated punctuation marks",
		},
		{
			name:    "triple exclamation",
			msg:     "error!!!",
			wantOk:  false,
			wantErr: "log message must not contain repeated punctuation marks",
		},
		{
			name:    "double question mark",
			msg:     "really??",
			wantOk:  false,
			wantErr: "log message must not contain repeated punctuation marks",
		},
		{
			name:    "triple question mark",
			msg:     "what???",
			wantOk:  false,
			wantErr: "log message must not contain repeated punctuation marks",
		},
		{
			name:    "double dot",
			msg:     "loading..",
			wantOk:  false,
			wantErr: "log message must not contain repeated punctuation marks",
		},
		{
			name:    "triple dot",
			msg:     "waiting...",
			wantOk:  false,
			wantErr: "log message must not contain repeated punctuation marks",
		},
		{
			name:    "mixed punctuation",
			msg:     "what?!",
			wantOk:  false,
			wantErr: "log message must not contain repeated punctuation marks",
		},
		{
			name:    "multiple punctuation types",
			msg:     "error!?",
			wantOk:  false,
			wantErr: "log message must not contain repeated punctuation marks",
		},
		{
			name:    "punctuation in middle",
			msg:     "check!! this",
			wantOk:  false,
			wantErr: "log message must not contain repeated punctuation marks",
		},
		{
			name:    "emoji with ascii",
			msg:     "done 🎉",
			wantOk:  true,
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOk, gotErr := CheckNoSpecials(tt.msg)

			if gotOk != tt.wantOk {
				t.Errorf("CheckNoSpecials() gotOk = %v, want %v", gotOk, tt.wantOk)
			}

			if gotErr != tt.wantErr {
				t.Errorf("CheckNoSpecials() gotErr = %q, want %q", gotErr, tt.wantErr)
			}
		})
	}
}
