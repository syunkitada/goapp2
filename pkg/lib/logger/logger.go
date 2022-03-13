package logger

import (
	"fmt"
	"sync"

	"github.com/rs/xid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/syunkitada/goapp2/pkg_infra/lib/infra_logger"
	"github.com/syunkitada/goapp2/pkg_infra/lib/infra_os"
)

var disableExit bool

var logger *zap.Logger

type Config struct {
	OutputPaths []string
	Level       string
	Encoding    string
	DisableExit bool
}

const (
	LevelDebug = "Debug"
	LevelInfo  = "Info"
	LevelWarn  = "Warn"
)

func NewZapCoreLevel(levelStr string) (level zapcore.Level) {
	switch levelStr {
	case "Debug":
		level = zap.DebugLevel
	case "Info":
		level = zap.InfoLevel
	case "Warn":
		level = zap.WarnLevel
	default:
		level = zap.InfoLevel
	}
	return
}

func Init(conf *Config) {
	if conf.Encoding == "" {
		conf.Encoding = "json"
	}
	disableExit = conf.DisableExit
	zapConf := zap.Config{
		Level:       zap.NewAtomicLevelAt(NewZapCoreLevel(conf.Level)),
		Development: false,
		Encoding:    conf.Encoding,
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
		infra_os.Exit(disableExit, 1)
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
	infra_logger.Fatal(disableExit, logger, msg, append([]zap.Field{zap.Inline(tctx)}, fields...)...)
}

func Sync() {
	logger.Sync()
}
