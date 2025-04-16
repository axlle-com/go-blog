package logger

import (
	"fmt"
	"github.com/axlle-com/blog/app/config"
	"github.com/davecgh/go-spew/spew"
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
		logger.SetLevel(logrus.TraceLevel)
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
	getLogger().Log(level, args...)
	return
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

	fields := logrus.Fields{
		"file":     file,
		"line":     strconv.Itoa(line),
		"function": fnName,
	}
	getLogger().WithFields(fields).Log(level, args...)
}

func log(level logrus.Level, args ...any) {
	if conf.LogLevel() < int(level) {
		return
	}

	getLogger().Log(level, args...)
}

func Error(args ...any) {
	logWithCaller(logrus.ErrorLevel, args...)
}

func Errorf(format string, a ...any) {
	logWithCaller(logrus.ErrorLevel, fmt.Errorf(format, a...))
}

func Warning(args ...any) {
	logWithCaller(logrus.WarnLevel, args...)
}

func Warningf(format string, a ...any) {
	logWithCaller(logrus.WarnLevel, fmt.Errorf(format, a...))
}

func Fatal(args ...any) {
	logWithCaller(logrus.FatalLevel, args...)
}

func Fatalf(format string, a ...any) {
	logWithCaller(logrus.FatalLevel, fmt.Errorf(format, a...))
}

func Info(args ...any) {
	log(logrus.InfoLevel, args...)
}

func Infof(format string, a ...any) {
	log(logrus.InfoLevel, fmt.Sprintf(format, a...))
}

func Debug(args ...any) {
	log(logrus.DebugLevel, args...)
}

func Dump(args ...any) {
	spew.Dump(args)
}

func Debugf(format string, a ...any) {
	log(logrus.DebugLevel, fmt.Sprintf(format, a...))
}
