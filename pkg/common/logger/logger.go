package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
)

type LoggerInterface interface {
	Request(*gin.Context)
	Error(error)
	Info(string)
}

type log struct {
	logger *logrus.Logger
}

func New() LoggerInterface {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
	return &log{logger: logger}
}

func (f *log) Request(c *gin.Context) {
	f.logger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
	}).Info("Request received")
}

func (f *log) Error(err error) {
	f.logger.WithFields(logrus.Fields{
		"error": err,
	}).Error("An error occurred")
}

func (f *log) Info(message string) {
	f.logger.Info(message)
}
