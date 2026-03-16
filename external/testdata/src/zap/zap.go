package zap

import (
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()

	logger.Info("Starting server")     // want "the log message must begin with a lowercase letter"
	logger.Error("ошибка подключения") // want "log message must be in English \\(ASCII only\\)"
	logger.Warn("warning!!!")          // want "log message must not contain repeated punctuation marks"
	logger.Info("user password: 123")  // want "the log message contains sensitive data"

	logger.Info("starting server")
	logger.Error("connection failed")
	logger.Warn("warning message")
	logger.Info("user authenticated successfully")

	logger.Infof("Server started on port %d", 8080)   // want "the log message must begin with a lowercase letter"
	logger.Errorf("Failed to connect: %v", "timeout") // want "the log message must begin with a lowercase letter"

	sugar := logger.Sugar()
	sugar.Info("Starting server") // want "the log message must begin with a lowercase letter"
	sugar.Infof("Port: %d", 8080) // want "the log message must begin with a lowercase letter"
}
