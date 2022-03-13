package infra_logger

import "go.uber.org/zap"

func Fatal(disable bool, logger *zap.Logger, msg string, fields ...zap.Field) {
	if !disable {
		logger.Fatal(msg, fields...)
	}
}
