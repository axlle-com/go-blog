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

		switch status {
		case http.StatusUnauthorized:
			if isAPI {
				ctx.JSON(status, gin.H{"message": "Unauthorized"})
			} else {
				ctx.HTML(status, e.view.View("error"), gin.H{"title": "Unauthorized", "error": status})
			}

		case http.StatusForbidden:
			if isAPI {
				ctx.JSON(status, gin.H{"message": "Forbidden"})
			} else {
				ctx.HTML(status, e.view.View("error"), gin.H{"title": "Forbidden", "error": status})
			}

		case http.StatusInternalServerError:
			errMsg := ctx.Errors.String()
			if errMsg == "" {
				errMsg = "Internal Server Error"
			}
			if isAPI {
				ctx.JSON(status, gin.H{"message": errMsg})
			} else {
				ctx.HTML(status, e.view.View("error"), gin.H{"message": errMsg, "error": status})
			}
		}
	}
}
