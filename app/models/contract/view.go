package contract

import "github.com/gin-gonic/gin"

type View interface {
	SetRouter(router *gin.Engine)
	Load()
	RenderToString(name string, data any) (string, error)
	View(name string) string
}
