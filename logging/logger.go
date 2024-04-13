package logging

import (
	"github.com/sirupsen/logrus"
	"sync"
)

var instance *logrus.Logger
var once sync.Once

func GetLogger() *logrus.Logger {
	once.Do(func() {
		instance = logrus.StandardLogger()
	})
	return instance
}
