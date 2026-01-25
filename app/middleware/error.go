package middleware

import (
	"net/http"
	"strings"

	"github.com/axlle-com/blog/app/logger"
	"github.com/axlle-com/blog/app/models/contract"
	"github.com/gin-gonic/gin"
)

func NewError(view contract.View) *Error {
	return &Error{
		view: view,
	}
}

type Error struct {
	view contract.View
}

func (e *Error) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		status := ctx.Writer.Status()
		if status < 400 {
			return
		}

		written := ctx.Writer.Written()
		isAborted := ctx.IsAborted()
		accept := ctx.GetHeader("Accept")
		isAPI := strings.Contains(accept, "application/json") ||
			strings.Contains(ctx.GetHeader("X-Requested-With"), "XMLHttpRequest")

		logger.WithRequest(ctx).Debugf("[ErrorMiddleware] Status: %d, Written: %v, IsAborted: %v, Path: %s, isAPI: %v",
			status, written, isAborted, ctx.Request.URL.Path, isAPI)

		if written {
			logger.WithRequest(ctx).Debugf("[ErrorMiddleware] Response already written, skipping")
			return
		}

		tFunc := getT(ctx)
		withSettings := func(data gin.H) gin.H {
			if data == nil {
				data = gin.H{}
			}

			if settings, ok := data["settings"]; ok {
				switch s := settings.(type) {
				case gin.H:
					s["T"] = tFunc
				case map[string]any:
					s["T"] = tFunc
				default:
					data["settings"] = gin.H{"T": tFunc}
				}
			} else {
				data["settings"] = gin.H{"T": tFunc}
			}

			return data
		}

		switch status {
		case http.StatusUnauthorized:
			if isAPI {
				ctx.JSON(status, gin.H{"message": "Unauthorized"})
			} else {
				ctx.HTML(status, e.view.ViewStatic("error"), withSettings(gin.H{"title": "Unauthorized", "error": status}))
			}

		case http.StatusForbidden:
			if isAPI {
				ctx.JSON(status, gin.H{"message": "Forbidden"})
			} else {
				ctx.HTML(status, e.view.ViewStatic("error"), withSettings(gin.H{"title": "Forbidden", "error": status}))
			}

		case http.StatusInternalServerError:
			errMsg := ctx.Errors.String()
			if errMsg == "" {
				errMsg = "Internal Server Error"
			}

			if isAPI {
				ctx.JSON(status, gin.H{"message": errMsg})
			} else {
				ctx.HTML(status, e.view.ViewStatic("error"), withSettings(gin.H{"message": errMsg, "error": status}))
			}
		}
	}
}

func getT(ctx *gin.Context) func(id string, data map[string]any, n ...int) string {
	if v, ok := ctx.Get(CtxTKey); ok {
		if f, ok := v.(func(string, map[string]any, ...int) string); ok {
			return f
		}
	}

	return func(id string, _ map[string]any, _ ...int) string { return id }
}
