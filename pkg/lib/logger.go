package logger

import (
	"fmt"
	"os"
	"sync"

	"github.com/rs/xid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

type Config struct {
	OutputPaths []string
	Level       string
}

const (
	LevelDebug = "Debug"
	LevelInfo  = "Info"
	LevelWarn  = "Warn"
)

func Init(conf *Config) {
	var level zapcore.Level
	switch conf.Level {
	case "Debug":
		level = zap.DebugLevel
	case "Info":
		level = zap.InfoLevel
	case "Warn":
		level = zap.WarnLevel
	default:
		level = zap.InfoLevel
	}
	zapConf := zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
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
		OutputPaths: conf.OutputPaths,
	}

	var err error
	logger, err = zapConf.Build()
	if err != nil {
		fmt.Println("Failed to initialize logger")
		os.Exit(1)
	}

}

type TraceContext struct {
	mtx     *sync.Mutex
	traceId string
}

func (self *TraceContext) MarshalLogObject(enc zapcore.ObjectEncoder) (err error) {
	enc.AddString("traceId", self.traceId)
	return
}

func (self *TraceContext) GetTraceId() (traceId string) {
	return self.traceId
}

func NewTraceContext() (tctx *TraceContext) {
	return &TraceContext{
		mtx:     new(sync.Mutex),
		traceId: xid.New().String(),
	}
}

func Debug(tctx *TraceContext, msg string, fields ...zap.Field) {
	logger.Debug(msg, append([]zap.Field{zap.Inline(tctx)}, fields...)...)
}

func Info(tctx *TraceContext, msg string, fields ...zap.Field) {
	logger.Info(msg, append([]zap.Field{zap.Inline(tctx)}, fields...)...)
}

func Warn(tctx *TraceContext, msg string, fields ...zap.Field) {
	logger.Warn(msg, append([]zap.Field{zap.Inline(tctx)}, fields...)...)
}

func Error(tctx *TraceContext, msg string, fields ...zap.Field) {
	logger.Error(msg, append([]zap.Field{zap.Inline(tctx)}, fields...)...)
}

func Fatal(tctx *TraceContext, msg string, fields ...zap.Field) {
	logger.Fatal(msg, append([]zap.Field{zap.Inline(tctx)}, fields...)...)
}

func Sync() {
	logger.Sync()
}
