package slog

import "log/slog"

func main() {
	slog.Info("Starting server")    // want "the log message must begin with a lowercase letter"
	slog.Info("starting server")    // ok
	slog.Error("ошибка")            // want "log message must be in English \\(ASCII only\\)"
	slog.Warn("warning!!!")         // want "log message must not contain repeated punctuation marks"
	slog.Info("user password: 123") // want "the log message contains sensitive data"
}
