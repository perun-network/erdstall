package operator

import "github.com/sirupsen/logrus"

func init() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
}
