package logger

import (
	"context"
	log "github.com/sirupsen/logrus"
	"os"
)

var RequestId = "RequestId"

type Logger struct {
	*log.Logger
}

func NewLogger() *Logger {
	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(log.InfoLevel)

	return &Logger{
		logger,
	}
}

func (l *Logger) WithRequestId(ctx context.Context) *log.Entry {
	value := ctx.Value(RequestId)
	v, ok := value.(string)
	if !ok {
		v = ""
	}

	entry := l.WithFields(log.Fields{
		"RequestId": v,
	})

	return entry
}
