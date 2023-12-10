package tasks

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ZerologAdapter struct {
	logger *zerolog.Event
}

func NewZerologAdapter() *ZerologAdapter {
	return &ZerologAdapter{}
}

func (l *ZerologAdapter) fields(keyvals []interface{}) {
	if len(keyvals)%2 != 0 {
		l.logger.Err(fmt.Errorf("odd number of keyvals pairs: %v", keyvals))
		return
	}

	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			key = fmt.Sprintf("%v", keyvals[i])
		}
		l.logger.Any(key, keyvals[i+1])
	}
}

func (l *ZerologAdapter) Debug(msg string, keyvals ...interface{}) {
	l.logger = log.Debug()
	l.fields(keyvals)
	l.logger.Msg(msg)
}

func (l *ZerologAdapter) Info(msg string, keyvals ...interface{}) {
	l.logger = log.Info()
	l.fields(keyvals)
	l.logger.Msg(msg)
}

func (l *ZerologAdapter) Warn(msg string, keyvals ...interface{}) {
	l.logger = log.Warn()
	l.fields(keyvals)
	l.logger.Msg(msg)
}

func (l *ZerologAdapter) Error(msg string, keyvals ...interface{}) {
	l.logger = log.Error()
	l.fields(keyvals)
	l.logger.Msg(msg)
}
