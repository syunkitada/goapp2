package logger

import (
	"os"

	"go.uber.org/zap"
)

func FatalInternal(disable bool, logger *zap.Logger, msg string, fields ...zap.Field) {
	if !disable {
		logger.Fatal(msg, fields...)
	}
}

func Exit(disable bool, code int) {
	if !disable {
		os.Exit(code)
	}
}
