package authhandelers

import (
	"log"

	"github.com/darkard2003/wormhole/models"
	"github.com/darkard2003/wormhole/services/dbservice"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserInput struct {
	Username string  `json:"username" binding:"required"`
	Password string  `json:"password" binding:"required"`
	Email    *string `json:"email,omitempty" binding:"omitempty,email"`
}

func SignUpHandler(ctx *gin.Context) {
	var userInput UserInput
	if err := ctx.ShouldBindJSON(&userInput); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	user := &models.User{}
	user.Username = userInput.Username
	user.Email = userInput.Email

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		ctx.JSON(500, gin.H{"error": "Failed to hash password"})
	}
	user.Password = string(passwordHash)

	db := dbservice.GetDBService()
	err = db.CreateUser(user)
	if err != nil {
		log.Println("Error creating user:", err)
		ctx.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}
	ctx.JSON(200, gin.H{"message": "User signed up successfully"})
}
