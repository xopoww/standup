package logging

import (
	"github.com/golang-migrate/migrate/v4"
	"go.uber.org/zap"
)

type migrateLogger struct {
	l       *zap.Logger
	verbose bool
}

func MigrateLogger(l *zap.Logger, verbose bool) migrate.Logger {
	return migrateLogger{
		l:       l,
		verbose: verbose,
	}
}

func (l migrateLogger) Printf(format string, v ...interface{}) {
	l.l.Sugar().Debugf(format, v...)
}

func (l migrateLogger) Verbose() bool {
	return l.verbose
}
