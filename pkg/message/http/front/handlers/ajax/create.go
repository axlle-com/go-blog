package web

import (
	"net/http"

	"github.com/axlle-com/blog/app/errutil"
	"github.com/axlle-com/blog/app/http/response"
	"github.com/axlle-com/blog/pkg/message/form"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func (c *messageController) CreateMessage(ctx *gin.Context) {
	userUUID := ctx.GetString("user_uuid")
	if userUUID == "" {
		userUUID = ctx.GetString("guest_uuid")
	}

	formName := ctx.Param("form")
	if formName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": c.T(ctx, "ui.message.field_form_name")})
		return
	}

	tempForm, err := form.NewForm(formName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctx.ShouldBindBodyWith(tempForm, binding.JSON); err != nil {
		formError := errutil.NewErrors(err)
		ctx.JSON(
			http.StatusBadRequest,
			response.Fail(http.StatusBadRequest, formError.Message, formError.Errors),
		)
		ctx.Abort()
		return
	}

	switch newForm := tempForm.(type) {
	case *form.Contact:
		newForm.UserUUID = userUUID
		c.mailService.SendContact(newForm)
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected form type"})
	}

	ctx.JSON(http.StatusOK, gin.H{"message": c.T(ctx, "ui.message.message_sent")})
}
