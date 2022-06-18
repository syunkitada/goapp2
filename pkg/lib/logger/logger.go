package logger

import (
	"fmt"
	"sync"

	"github.com/rs/xid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/syunkitada/goapp2/pkg/lib/struct_utils"
	"github.com/syunkitada/goapp2/pkg_infra/lib/infra_logger"
	"github.com/syunkitada/goapp2/pkg_infra/lib/infra_os"
)

var disableExit bool

var logger *zap.Logger
var sugar *zap.SugaredLogger

type Config struct {
	OutputPaths []string
	Level       string
	Encoding    string
	DisableExit bool
}

var conf = Config{
	OutputPaths: []string{"stdout"},
	Level:       "info",
	Encoding:    "json",
	DisableExit: false,
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

func Init(conf2 *Config) {
	struct_utils.MergeStruct(conf, conf2)

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
		fmt.Println("Failed to initialize logger", err.Error())
		infra_os.Exit(disableExit, 1)
	}

	sugar = logger.Sugar()
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

func Debugf(tctx *TraceContext, msg string, fields ...interface{}) {
	sugar.Debugf(msg, fields...)
}

func Info(tctx *TraceContext, msg string, fields ...zap.Field) {
	logger.Info(msg, append([]zap.Field{zap.Inline(tctx)}, fields...)...)
}

func Infof(tctx *TraceContext, msg string, fields ...interface{}) {
	sugar.Infof(msg, fields...)
}

func Warn(tctx *TraceContext, msg string, fields ...zap.Field) {
	logger.Warn(msg, append([]zap.Field{zap.Inline(tctx)}, fields...)...)
}

func Warnf(tctx *TraceContext, msg string, fields ...interface{}) {
	sugar.Warnf(msg, fields...)
}

func Error(tctx *TraceContext, msg string, fields ...zap.Field) {
	logger.Error(msg, append([]zap.Field{zap.Inline(tctx)}, fields...)...)
}

func Errorf(tctx *TraceContext, msg string, fields ...interface{}) {
	sugar.Errorf(msg, fields...)
}

func Fatal(tctx *TraceContext, msg string, fields ...zap.Field) {
	infra_logger.Fatal(disableExit, logger, msg, append([]zap.Field{zap.Inline(tctx)}, fields...)...)
}

func Fatalf(tctx *TraceContext, msg string, fields ...interface{}) {
	sugar.Fatalf(msg, fields...)
}

func Sync() {
	logger.Sync()
}
