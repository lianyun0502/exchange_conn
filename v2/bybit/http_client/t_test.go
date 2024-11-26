package bybit_test

import (
	"os"
	"github.com/sirupsen/logrus"
)

var apiKey = "L7ksyiOdEgqg0gwIbf"
var secretKey = "0CVhyQmkwUDKWLcAP6NhtH7jB0P8XqSIVxE1"

var logger = &logrus.Logger{
	Out: os.Stderr,
	Formatter: &logrus.TextFormatter{
		ForceColors:     true,
		TimestampFormat: "2006-01-02 15:04:05.000000",
		FullTimestamp:   true,
	},
	Level: logrus.DebugLevel,
	Hooks: make(logrus.LevelHooks),
}