package common

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
	Logger.Out = os.Stdout
	Logger.Formatter = new(logrus.JSONFormatter)
	Logger.Level = logrus.DebugLevel
}