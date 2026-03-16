package testdata

import "go.uber.org/zap"

func main() {
	logger, _ := zap.NewProduction()
	logger.Info("Starting server")
	logger.Info("starting server")
	logger.Error("Ошибка")
	logger.Warn("warning!!!")
	logger.Info("user password: 123")
}
