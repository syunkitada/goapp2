package logger

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestBasic(t *testing.T) {
	a := assert.New(t)
	a.Equal(true, true)

	logFile := "/tmp/output.log"
	_, err := os.Stat(logFile)
	if err == nil {
		err = os.Remove(logFile)
		a.NoError(err)
	}

	Init(&Config{Encoding: "dummy", DisableExit: true}) // failed
	Init(&Config{OutputPaths: []string{logFile}, Level: "Debug", DisableExit: true})

	tctx := NewTraceContext()
	traceId := tctx.GetTraceId()
	a.Greater(len(traceId), 5)

	Debug(tctx, "debugmsg")
	Info(tctx, "infomsg")
	Warn(tctx, "warnmsg")
	Error(tctx, "errormsg")
	Fatal(tctx, "fatalmsg")
	Sync()

	bytes, err := ioutil.ReadFile(logFile)
	a.NoError(err)
	lines := strings.Split(string(bytes), "\n")
	a.Equal(5, len(lines))
}

func TestNewZapCoreLevel(t *testing.T) {
	a := assert.New(t)
	a.Equal(zap.DebugLevel, NewZapCoreLevel("Debug"))
	a.Equal(zap.InfoLevel, NewZapCoreLevel("Info"))
	a.Equal(zap.WarnLevel, NewZapCoreLevel("Warn"))
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
