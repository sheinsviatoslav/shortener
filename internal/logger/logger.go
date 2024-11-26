package logger

import (
	"go.uber.org/zap"
)

func Get() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	return logger
}
