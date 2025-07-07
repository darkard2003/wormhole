package handelers

import (
	"crypto"
	"encoding/hex"
	"fmt"

	"github.com/darkard2003/wormhole/models"
	"github.com/darkard2003/wormhole/services/dbservice"
	"github.com/gin-gonic/gin"
)

type UserInput struct {
	Username string  `json:"username" binding:"required"`
	Password string  `json:"password" binding:"required"`
	Email    *string `json:"email,omitempty" binding:"omitempty,email"`
}

func SignUpHandeler(ctx *gin.Context) {
	var userInput UserInput
	if err := ctx.ShouldBindJSON(&userInput); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	user := &models.User{}
	user.Username = userInput.Username
	user.Email = userInput.Email

	passwordHash := crypto.SHA256.New()
	passwordHash.Write([]byte(userInput.Password))
	user.Password = hex.EncodeToString(passwordHash.Sum(nil))

	db := dbservice.GetDBService()
	err := db.CreateUser(user)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to create user"})
		fmt.Println("Error creating user:", err)
		return
	}
	ctx.JSON(200, gin.H{"message": "User signed up successfully"})
}
