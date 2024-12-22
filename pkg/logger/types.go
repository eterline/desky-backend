package logger

import (
	"io"

	"github.com/sirupsen/logrus"
)

type LogWorker struct {
	*logrus.Entry
}

func ReturnEntry() LogWorker {
	return LogWorker{entry}
}

type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

type LoggingConfig struct {
	Renew bool   // If true log file will be created with starting time in name
	Dir   string // Directory for logging file

	_ struct{}
}
