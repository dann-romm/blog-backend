package app

import (
	"github.com/sirupsen/logrus"
	"os"
)

func SetLogrus(level string) {
	logrusLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrusLevel)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// file, err := os.OpenFile("/logs/app.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	// logrus.SetOutput(file)
	logrus.SetOutput(os.Stdout)
}
