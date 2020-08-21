package common

import (
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	globalLogger *logrus.Logger
	onceGL       sync.Once
)

// Logger returns the global logger.
func Logger() *logrus.Logger {
	onceGL.Do(func() {
		globalLogger = logrus.New()
	})
	return globalLogger
}
