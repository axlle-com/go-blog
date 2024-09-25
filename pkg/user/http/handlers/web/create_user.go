package web

import (
	. "github.com/axlle-com/blog/pkg/user/http/models"
	user "github.com/axlle-com/blog/pkg/user/models"
	"github.com/axlle-com/blog/pkg/user/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func CreateUser(c *gin.Context) {

	var authInput AuthInput
	userRepo := repository.NewRepo()

	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userFound, err = userRepo.GetByEmail(authInput.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if userFound != nil && userFound.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already used"})
		return
	}

	newUser := user.User{
		Email: authInput.Email,
	}

	if err := userRepo.Create(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": newUser})

}
