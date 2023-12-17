package logger

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

func InitLogger() *logrus.Logger {
	l := logrus.New()
	l.SetOutput(os.Stdout)

	return l
}

const loggerKeyContext = "logrus.Logger"

func ContextWithLogger(ctx context.Context, l *logrus.Logger) context.Context {
	return context.WithValue(ctx, loggerKeyContext, l)
}

func Warn(ctx context.Context, message string) {
	l := loggerFromContext(ctx)

	if l != nil {
		l.Warn(message)
	} else {
		logrus.Warn(message)
	}
}

func loggerFromContext(ctx context.Context) *logrus.Logger {
	if l, ok := ctx.Value(loggerKeyContext).(*logrus.Logger); ok {
		return l
	}

	return &logrus.Logger{}
}
