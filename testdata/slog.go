package testdata

import "log/slog"

func main() {
	slog.Info("Starting server")
	slog.Info("starting server")
	slog.Error("Ошибка")
	slog.Warn("warning!!!")
	slog.Info("user password: 123")
}
