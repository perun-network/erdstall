package caching

import "github.com/sirupsen/logrus"

func init() {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
}
