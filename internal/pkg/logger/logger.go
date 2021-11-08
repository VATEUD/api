package logger

import (
	"api/utils"
	"github.com/onrik/logrus/sentry"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func New() error {
	Log = logrus.New()

	hook, err := sentry.NewHook(sentry.Options{
		Dsn: utils.Getenv("SENTRY_DSN", ""),
	}, logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel)

	if err != nil {
		return err
	}

	defer hook.Flush()

	Log.AddHook(hook)

	return nil
}
