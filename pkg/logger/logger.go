package logger

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

var entry *logrus.Entry
var HookTargets []io.Writer

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

func (h *writerHook) Fire(entr *logrus.Entry) error {
	str, err := entr.String()
	if err != nil {
		return err
	}
	for _, w := range h.Writer {
		w.Write([]byte(str))
	}
	return err
}

func (h *writerHook) Levels() []logrus.Level {
	return h.LogLevels
}

func InitLogger(path, filename string) error {
	l := logrus.New()
	l.SetReportCaller(true)

	l.Formatter = &easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[%lvl%]: %time% - %msg% \n",
	}

	fp := filepath.Join(path, filename)
	logFile, err := os.OpenFile(fp, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return err
	}

	HookTargets = append(HookTargets, logFile)
	HookTargets = append(HookTargets, os.Stdout)

	l.SetOutput(io.Discard)
	l.AddHook(&writerHook{
		Writer:    HookTargets,
		LogLevels: logrus.AllLevels,
	})
	l.SetLevel(logrus.TraceLevel)
	entry = logrus.NewEntry(l)

	l.Debug("App logging initialized")
	return nil
}

