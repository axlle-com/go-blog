package web

import (
	"errors"
	. "github.com/axlle-com/blog/pkg/user/http/models"
	user "github.com/axlle-com/blog/pkg/user/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func (c *controller) CreateUser(ctx *gin.Context) {

	var authInput AuthInput

	if err := ctx.ShouldBindJSON(&authInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userFound, err = c.userService.GetByEmail(authInput.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if userFound != nil && userFound.ID != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "username already used"})
		return
	}

	newUser := user.User{
		Email: authInput.Email,
	}

	if err := c.userService.Create(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": newUser})

}
