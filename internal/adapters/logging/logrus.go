package logging

import (
	"github.com/sirupsen/logrus"
	"os"
)

var logger *logrus.Logger

func InitLogger() {
	logger = logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		ForceQuote:      true,
		PadLevelText:    true,
		TimestampFormat: "2006-01-02 15:04:05",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})
}

func GetLogger() *logrus.Logger {
	return logger
}
