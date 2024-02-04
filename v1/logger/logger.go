package logger

import (
	"github.com/mattn/go-colorable"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() (*zap.Logger, error) {
	lg := zap.NewDevelopmentEncoderConfig()
	lg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(lg),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.DebugLevel,
	)), nil
}

func Decorate() fx.Option {
	return fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
		return &fxevent.ZapLogger{Logger: log}
	})
}
