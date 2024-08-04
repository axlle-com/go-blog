package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
)

type Logger interface {
	Request(*gin.Context)
	Error(error)
	Info(string)
}

type log struct {
	logger *logrus.Logger
	error
	file     string
	line     string
	function string
	message  string
}

func New() Logger {
	logger := logrus.New()
	//logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
		line = 0
	}

	function := "unknown"
	if pc, _, _, ok := runtime.Caller(1); ok {
		function = runtime.FuncForPC(pc).Name()
	}
	return &log{
		logger:   logger,
		file:     file,
		line:     string(rune(line)),
		function: function,
	}
}

func (f *log) Request(c *gin.Context) {
	f.logger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
	}).Info("Request received")
}

func (f *log) Error(err error) {
	f.logger.WithFields(logrus.Fields{
		"error":    err,
		"file":     f.file,
		"line":     f.line,
		"function": f.function,
	}).Error("An error occurred")
}

func (f *log) Info(message string) {
	f.logger.Info(message)
}