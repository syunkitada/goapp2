package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestBasic(t *testing.T) {
	a := assert.New(t)
	a.Equal(true, true)

	Init(&Config{})
}

func BenchmarkLogger(b *testing.B) {
	Init(&Config{
		Level:       LevelInfo,
		OutputPaths: []string{"/tmp/output.log"},
	})
	tctx := NewTraceContext()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info(tctx, "test")
	}
}

func BenchmarkZap(b *testing.B) {
	zapConf := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths: []string{"/tmp/output.log"},
	}
	logger, _ := zapConf.Build()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("test")
	}
}
