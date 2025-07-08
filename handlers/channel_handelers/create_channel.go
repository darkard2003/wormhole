package channelhandelers

import (
	"log"

	"github.com/darkard2003/wormhole/models"
	"github.com/darkard2003/wormhole/services/dbservice"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CreateChannel(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")

	if !exists || userId == nil {
		ctx.JSON(400, gin.H{"error": "Unauthorized"})
		return
	}

	var request struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description,omitempty"`
		Protected   bool    `json:"protected,omitempty"`
		Password    *string `json:"password,omitempty"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	passwordHash := ""
	if request.Password != nil {
		passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(*request.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Error hashing password:", err)
			ctx.JSON(500, gin.H{"error": "Failed to hash password"})
			return
		}

		passwordHash = string(passwordHashBytes)
	}

	channel := &models.Channel{
		UserID:      userId.(int),
		Name:        request.Name,
		Description: request.Description,
		Protected:   request.Protected,
		Password:    passwordHash,
	}

	err := dbservice.GetDBService().CreateChannel(channel, userId.(int))

	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to create channel"})
		return
	}
	ctx.JSON(201, gin.H{
		"message": "Channel created successfully",
	})
}
