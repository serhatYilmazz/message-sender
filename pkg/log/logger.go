package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

func NewLogger(environment string) *logrus.Logger {
	logger := logrus.New()

	switch environment {
	case "development":
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
		})
		logger.SetLevel(logrus.DebugLevel)
	case "prod":
		logger.SetFormatter(&logrus.JSONFormatter{})
		logger.SetLevel(logrus.InfoLevel)
	}

	logrus.SetOutput(os.Stdout)
	return logger
}
