package logger

import (
	"os"
	"runtime"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Setup log rotation
	path := "logs/app.log"
	var writer *rotatelogs.RotateLogs
	var err error

	// On Windows, don't create symlinks to avoid permission issues
	if runtime.GOOS == "windows" {
		writer, err = rotatelogs.New(
			path+".%Y%m%d.log",
			rotatelogs.WithRotationTime(24*time.Hour),
			rotatelogs.WithMaxAge(7*24*time.Hour),
		)
	} else {
		writer, err = rotatelogs.New(
			path+".%Y%m%d.log",
			rotatelogs.WithLinkName(path),
			rotatelogs.WithRotationTime(24*time.Hour),
			rotatelogs.WithMaxAge(7*24*time.Hour),
		)
	}

	if err == nil {
		log.SetOutput(writer)
	}

	return log
}
