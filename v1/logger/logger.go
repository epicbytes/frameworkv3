package logger

import (
	"github.com/mattn/go-colorable"
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
