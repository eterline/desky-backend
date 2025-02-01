package logger

import (
	"io"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

const (
	_ EnvValue = iota
	LOCAL
	DEVELOP
	PRODUCTION
)

var (
	entry       *logrus.Entry
	HookTargets []io.Writer

	once sync.Once
)

type LogWorker struct {
	*logrus.Entry
}

func ReturnEntry() LogWorker {
	return LogWorker{entry}
}

func InitLogger(options ...LoggerOptionFunc) error {
	opts := mustOptions(options...)

	logFile, err := os.OpenFile(returnName(opts), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	l := logrus.New()

	if opts.isFormat {
		l.Formatter = returnFormatter(opts)
	}

	HookTargets = append(HookTargets, logFile, os.Stdout)

	l.SetOutput(io.Discard)
	l.AddHook(&writerHook{
		Writer:    HookTargets,
		LogLevels: logrus.AllLevels,
	})

	l.SetLevel(opts.level)
	l.SetReportCaller(true)

	once.Do(func() {
		entry = logrus.NewEntry(l)
		l.Debugf(
			"logger initialized with log level: %s and env: %s",
			opts.level, opts.env,
		)
	})

	return nil
}
