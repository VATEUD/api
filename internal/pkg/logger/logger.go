package logger

import (
	"api/utils"
	"github.com/onrik/logrus/filename"
	"github.com/onrik/logrus/sentry"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

const logFileName = "logs/log.txt"

var Log *logrus.Logger

func New() error {
	Log = logrus.New()

	fHook := filename.NewHook()
	fHook.Field = "file"
	Log.AddHook(fHook)

	hook, err := sentry.NewHook(sentry.Options{
		Dsn: utils.Getenv("SENTRY_DSN", ""),
	}, logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel)

	if err != nil {
		return err
	}

	defer hook.Flush()

	Log.AddHook(hook)

	pathMap := lfshook.PathMap{
		logrus.PanicLevel: "logs/api_panic.log",
		logrus.InfoLevel:  "logs/api_info.log",
		logrus.ErrorLevel: "logs/api_error.log",
		logrus.DebugLevel: "logs/api_debug.log",
		logrus.FatalLevel: "logs/api_fatal.log",
		logrus.WarnLevel:  "logs/api_warn.log",
	}

	Log.AddHook(lfshook.NewHook(pathMap, &logrus.JSONFormatter{}))

	return nil
}
