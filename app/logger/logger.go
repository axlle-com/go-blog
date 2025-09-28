package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/axlle-com/blog/app/config"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/jc633/rotatelogs"
	"github.com/sirupsen/logrus"
)

type WriterHook struct {
	Writer    io.Writer
	LogLevels []logrus.Level
	Formatter logrus.Formatter
}

func (hook *WriterHook) Fire(entry *logrus.Entry) error {
	line, err := hook.Formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write(line)
	return err
}

func (hook *WriterHook) Levels() []logrus.Level {
	return hook.LogLevels
}

var (
	once         sync.Once
	logrusLogger *logrus.Logger
	conf         = config.Config()
)

func getLogger() *logrus.Logger {
	once.Do(func() {
		logrusLogger = logrus.New()

		// 1) Настраиваем консольный форматтер (с цветом):
		consoleFmt := &logrus.TextFormatter{
			FullTimestamp:             true,
			TimestampFormat:           "2006-01-02 15:04:05.000",
			DisableQuote:              true,
			ForceColors:               true,
			EnvironmentOverrideColors: true,
			PadLevelText:              true,
		}
		logrusLogger.SetFormatter(consoleFmt)
		logrusLogger.SetOutput(os.Stdout)

		// 2) Настраиваем ротацию
		rl, err := rotatelogs.New(
			conf.RuntimeFolder("logs/app-%Y-%m-%d.log"),
			rotatelogs.WithLinkName(conf.RuntimeFolder("logs/app.log")),
			rotatelogs.WithRotationTime(24*time.Hour),
			rotatelogs.WithMaxAge(7*24*time.Hour),
		)

		if err == nil {
			// 3) Добавляем Hook, который пишет в rl своим “plain” форматтером
			fileFmt := &logrus.TextFormatter{
				FullTimestamp:   true,
				TimestampFormat: "2006-01-02 15:04:05.000",
				DisableQuote:    true,
				DisableColors:   true, // критично
				PadLevelText:    true,
			}
			logrusLogger.AddHook(&WriterHook{
				Writer:    rl,
				LogLevels: logrus.AllLevels,
				Formatter: fileFmt,
			})
		}

		logrusLogger.SetLevel(logrus.TraceLevel)
	})

	return logrusLogger
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

func WithRequest(ctx *gin.Context) *logrus.Entry {
	base := getLogger()
	fields := logrus.Fields{
		"request_uuid": ctx.GetString("request_uuid"),
	}
	return base.WithFields(fields)
}
