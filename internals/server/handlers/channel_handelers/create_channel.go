package channelhandelers

import (
	"log"

	"github.com/darkard2003/wormhole/internals/services/db"
	"github.com/darkard2003/wormhole/internals/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CreateChannelHandler(db db.DBInterface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, exists := ctx.Get("userId")
		userIdInt, ok := userId.(int)
		if !ok {
			log.Println("Error casting userId to int")
			ctx.JSON(500, gin.H{"error": "Internal server error"})
			return
		}

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

		_, err := db.CreateChannel(userIdInt, request.Name, request.Description, request.Protected, passwordHash)

		if err != nil {
			httpError := utils.DBToHttpError(err)
			ctx.JSON(httpError.Code, gin.H{"error": httpError.Response})
			return
		}
		ctx.JSON(201, gin.H{
			"message": "Channel created successfully",
		})

	}
}
