package models

import (
	"bytes"
	"github.com/gin-gonic/gin"
)

type ResponseWriterWrapper struct {
	gin.ResponseWriter
	Buffer *bytes.Buffer
}

func (rw *ResponseWriterWrapper) Write(data []byte) (int, error) {
	return rw.Buffer.Write(data)
}
