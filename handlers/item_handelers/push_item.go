package itemhandelers

import (
	"log"
	"net/http"

	"github.com/darkard2003/wormhole/models"
	"github.com/darkard2003/wormhole/services/db"
	"github.com/darkard2003/wormhole/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type PushItemRequest struct {
	ChannelName     string `json:"channel_name" binding:""`
	ChannelPassword string `json:"channel_password"`
	Type            string `json:"type" binding:"required"`
	Title           string `json:"title"`
	Salt            string `json:"salt" binding:"required"`
	IV              string `json:"iv" binding:"required"`
	Encryption      string `json:"encryption" binding:"required"`
	TextContent     string `json:"text_content"`
	Filename        string `json:"filename"`
	Filesize        int64  `json:"filesize"`
	MimeType        string `json:"mimetype"`
	FileCreatedAt   string `json:"file_created_at"`
	FileUpdatedAt   string `json:"file_updated_at"`
}

func PushItem(db db.DBInterface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, exists := ctx.Get("userId")
		if !exists || userId == nil {
			log.Println("Unauthorized")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		userIdInt, ok := userId.(int)
		if !ok {
			log.Println("Error casting userId to int")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		var request PushItemRequest
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
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Channel password required"})
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

		var item any

		switch request.Type {
		case "text":
			item = &models.TextItem{
				Item: models.Item{
					UserID:             userIdInt,
					ChannelID:          channel.ID,
					Type:               request.Type,
					Title:              request.Title,
					Salt:               request.Salt,
					IV:                 request.IV,
					EncryptionMetadata: request.Encryption,
				},
				Content: request.TextContent,
			}
		case "file":
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Not implemented"})
			return
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item type"})
			return
		}

		_, err = db.CreateItem(item)
		if err != nil {
			httpError := utils.DBToHttpError(err)
			ctx.JSON(httpError.Code, gin.H{"error": httpError.Response})
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Item created successfully"})
	}
}
