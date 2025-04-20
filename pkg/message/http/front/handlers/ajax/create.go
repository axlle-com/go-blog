package web

import (
	"github.com/axlle-com/blog/pkg/message/form"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

func (c *messageController) CreateMessage(ctx *gin.Context) {
	userUUID := ctx.GetString("user_uuid")
	if userUUID == "" {
		userUUID = ctx.GetString("guest_uuid")
	}

	var name form.Name
	if err := ctx.ShouldBindBodyWith(&name, binding.JSON); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "поле form_name: " + err.Error()})
		return
	}

	tempForm, err := name.NewForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctx.ShouldBindBodyWith(tempForm, binding.JSON); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch newForm := tempForm.(type) {
	case *form.Contact:
		newForm.UserUUID = userUUID
		c.mailService.SendContact(newForm)
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected form type"})
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Сообщение успешно отправлено"})
	return
}
