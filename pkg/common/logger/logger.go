package logger

import (
	"errors"
	"github.com/axlle-com/blog/pkg/common/config"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strconv"
	"sync"
)

var (
	once         sync.Once
	globalLogger *logrus.Logger
	conf         = config.Config()
)

func getLogger() *logrus.Logger {
	once.Do(func() {
		logger := logrus.New()
		logger.SetOutput(os.Stdout)
		logger.SetLevel(logrus.DebugLevel)
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05.000",
		})

		globalLogger = logger
	})

	return globalLogger
}

func logWithCaller(level logrus.Level, args ...any) {
	if conf.LogLevel() < int(level) {
		return
	}

	// runtime.Caller(3) - почему 3?
	//   1) сам logWithCaller
	//   2) функция-обёртка (например, Error() или Fatal() ниже)
	//   3) конечная точка вызова в прикладном коде
	pc, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "unknown"
		line = 0
	}

	fnName := "unknown"
	fn := runtime.FuncForPC(pc)
	if fn != nil {
		fnName = fn.Name()
	}

	getLogger().WithFields(logrus.Fields{
		"file":     file,
		"line":     strconv.Itoa(line),
		"function": fnName,
	}).Log(level, args...)
}

func log(level logrus.Level, args ...any) {
	if conf.LogLevel() < int(level) {
		return
	}

	getLogger().Log(level, args...)
}

func Error(err error) {
	logWithCaller(logrus.ErrorLevel, err)
}

func Errorf(msg string) {
	logWithCaller(logrus.ErrorLevel, errors.New(msg))
}

func Fatal(err error) {
	logWithCaller(logrus.FatalLevel, err)
}

func Info(args ...any) {
	log(logrus.InfoLevel, args...)
}

func Debug(args ...any) {
	log(logrus.DebugLevel, args...)
}
