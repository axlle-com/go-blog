package logger

import (
	"fmt"
	"io"
	"os"
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

func (hook *WriterHook) Levels() []logrus.Level { return hook.LogLevels }

var (
	once         sync.Once
	logrusLogger *logrus.Logger
	conf         = config.Config()
)

func getLogger() *logrus.Logger {
	once.Do(func() {
		logrusLogger = logrus.New()

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
		logrusLogger.SetLevel(logrus.TraceLevel)

		if shouldLogToFile() {
			if rl, err := rotatelogs.New(
				conf.RuntimeFolder("logs/app-%Y-%m-%d.log"),
				rotatelogs.WithLinkName(conf.RuntimeFolder("logs/app.log")),
				rotatelogs.WithRotationTime(24*time.Hour),
				rotatelogs.WithMaxAge(7*24*time.Hour),
			); err == nil {
				fileFmt := &logrus.TextFormatter{
					FullTimestamp:   true,
					TimestampFormat: "2006-01-02 15:04:05.000",
					DisableQuote:    true,
					DisableColors:   true,
					PadLevelText:    true,
				}
				logrusLogger.AddHook(&WriterHook{
					Writer:    rl,
					LogLevels: logrus.AllLevels,
					Formatter: fileFmt,
				})
			}
		}
	})
	return logrusLogger
}

func shouldLogToFile() bool {
	if value := os.Getenv("LOG_TO_FILE"); value != "" {
		if value == "1" || value == "true" || value == "TRUE" {
			return true
		}

		if number, err := strconv.Atoi(value); err == nil && number != 0 {
			return true
		}

		return false
	}

	return true
}

func log(level logrus.Level, args ...any) {
	if conf.LogLevel() < int(level) {
		return
	}
	getLogger().Log(level, args...)
}

func logf(level logrus.Level, format string, a ...any) {
	if conf.LogLevel() < int(level) {
		return
	}
	getLogger().Log(level, fmt.Sprintf(format, a...))
}

func Error(args ...any)                { log(logrus.ErrorLevel, args...) }
func Errorf(format string, a ...any)   { logf(logrus.ErrorLevel, format, a...) }
func Warning(args ...any)              { log(logrus.WarnLevel, args...) }
func Warningf(format string, a ...any) { logf(logrus.WarnLevel, format, a...) }
func Fatal(args ...any)                { log(logrus.FatalLevel, args...) }
func Fatalf(format string, a ...any)   { logf(logrus.FatalLevel, format, a...) }
func Info(args ...any)                 { log(logrus.InfoLevel, args...) }
func Infof(format string, a ...any)    { logf(logrus.InfoLevel, format, a...) }
func Debug(args ...any)                { log(logrus.DebugLevel, args...) }
func Debugf(format string, a ...any)   { logf(logrus.DebugLevel, format, a...) }
func Dump(args ...any)                 { spew.Dump(args...) }

func WithRequest(ctx *gin.Context) *logrus.Entry {
	return getLogger().WithFields(logrus.Fields{
		"request_uuid": ctx.GetString("request_uuid"),
	})
}
