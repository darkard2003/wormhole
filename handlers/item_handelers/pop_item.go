package itemhandelers

import (
	"log"
	"net/http"

	"github.com/darkard2003/wormhole/services/db"
	"github.com/darkard2003/wormhole/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type PopItemRequest struct {
	ChannelName     string `json:"channel_name"`
	ChannelPassword string `json:"channel_password"`
}

func PopItem(db db.DBInterface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, exists := ctx.Get("userId")
		if !exists || userId == nil {
			ctx.JSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		userIdInt, ok := userId.(int)
		if !ok {
			log.Println("Error casting userId to int")
			ctx.JSON(500, gin.H{"error": "Internal server error"})
		}

		var request PopItemRequest
		ctx.ShouldBindJSON(&request)

		if request.ChannelName == "" {
			request.ChannelName = "default"
		}

		channel, err := db.GetChannelByName(userIdInt, request.ChannelName)
		if err != nil {
			httpError := utils.DBToHttpError(err)
			ctx.JSON(httpError.Code, gin.H{"error": httpError.Response})
			return
		}

		if channel == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Channel not found"})
			return
		}

		if channel.Protected {
			if request.ChannelPassword == "" {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Channel password required"})
				return
			}
			passHashBytes, err := bcrypt.GenerateFromPassword([]byte(request.ChannelPassword), bcrypt.DefaultCost)
			if err != nil {
				log.Println("Error hashing password:", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}

			if string(passHashBytes) != channel.Password {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Channel password incorrect"})
				return
			}
		}

		item, err := db.PopLatestItem(channel.ID)
		if err != nil {
			httpError := utils.DBToHttpError(err)
			ctx.JSON(httpError.Code, gin.H{"error": httpError.Response})
			return
		}
		if item == nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "No items found"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"item": item})
	}
}
