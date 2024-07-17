package web

import (
	. "github.com/axlle-com/blog/pkg/user/http/models"
	"github.com/axlle-com/blog/pkg/user/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func Auth(c *gin.Context) {
	var authInput AuthInput
	userRepo := repository.NewUserRepository()

	if err := c.ShouldBindJSON(&authInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userFound, err = userRepo.GetUserByEmail(authInput.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if userFound == nil || userFound.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid password or login"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFound.PasswordHash), []byte(authInput.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid password or login"})
		return
	}

	token, err := userFound.SetAuthToken()
	if err != nil {
		log.Println("failed to generate token")
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to generate token"})
	}

	if err := userRepo.UpdateUser(userFound); err != nil {
		return
	}

	c.JSON(200, gin.H{
		"token": token,
	})
}
