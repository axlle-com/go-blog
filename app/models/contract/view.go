package contract

import "github.com/gin-gonic/gin"

type View interface {
	SetRouter(router *gin.Engine)
	Load()
	SetStatic()
	RenderToString(name string, data any) (string, error)
	View(name string) string
	ViewStatic(name string) string
}
