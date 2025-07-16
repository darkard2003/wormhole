package authhandelers

import (
	"log"
	"net/http"

	"github.com/darkard2003/wormhole/models"
	"github.com/darkard2003/wormhole/services/db"
	"github.com/darkard2003/wormhole/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserInput struct {
	Username string  `json:"username" binding:"required"`
	Password string  `json:"password" binding:"required"`
	Email    *string `json:"email,omitempty" binding:"omitempty,email"`
}

func SignUpHandlerHandler(db db.DBInterface) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		var userInput UserInput
		if err := ctx.ShouldBindJSON(&userInput); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		user := &models.User{}
		user.Username = userInput.Username
		user.Email = userInput.Email

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Error hashing password:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		}
		user.Password = string(passwordHash)

		id, err := db.CreateUser(user.Username, user.Password, user.Email)
		if err != nil {
			httpError := utils.DBToHttpError(err)
			ctx.JSON(httpError.Code, gin.H{"error": httpError.Response})
			return
		}

		id, err = db.CreateChannel(id, "default", "Default Channel", false, "")

		if err != nil {
			httpError := utils.DBToHttpError(err)
			ctx.JSON(httpError.Code, gin.H{"error": httpError.Response})
			return
		}

		ctx.JSON(200, gin.H{"message": "User signed up successfully"})
	}
}
