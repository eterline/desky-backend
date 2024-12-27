package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

var entry *logrus.Entry
var HookTargets []io.Writer

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

func InitLogger(path, filename string, level uint) error {
	l := logrus.New()
	l.SetReportCaller(true)
	l.SetLevel(logrus.Level(level))

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

func InitWithConfig(config LoggingConfig) error {
	logfile := "trace.log"

	if config.Renew {
		logfile = "trace." + time.Now().Format(time.RFC3339) + ".log"
	}

	if err := InitLogger(config.Dir, logfile, config.Level); err != nil {

		errDefault := InitLogger("", "trace.log", 0)

		if errDefault != nil {
			log.Fatalf("can't start app: %s", err)
		}

		ReturnEntry().Errorf("can't use log config: %s", err)
		return err
	}

	return nil
}
