package itemhandelers

import (
	"encoding/json"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/darkard2003/wormhole/internals/models"
	"github.com/darkard2003/wormhole/internals/services/db"
	storageservice "github.com/darkard2003/wormhole/internals/services/storage_service"
	"github.com/darkard2003/wormhole/internals/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type PopItemRequest struct {
	ChannelName     string `json:"channel_name"`
	ChannelPassword string `json:"channel_password"`
}

func PopItem(db db.DBInterface, s storageservice.StorageInterface) gin.HandlerFunc {
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
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Channel password required"})
				return
			}
			err := bcrypt.CompareHashAndPassword([]byte(channel.Password), []byte(request.ChannelPassword))
			if err != nil {
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

		switch item := item.(type) {
		case *models.TextItem:
			textItem := item
			ctx.JSON(http.StatusOK, gin.H{"item": textItem})
			return
		case *models.FileItem:
			fileId := item.FileId
			data, err := s.GetBlob(fileId)
			if err != nil {
				log.Println("Error getting blob:", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}

			mw := multipart.NewWriter(ctx.Writer)
			ctx.Header("Content-Type", mw.FormDataContentType())

			jsonPart, err := mw.CreateFormField("json")
			if err != nil {
				log.Println("Error creating json part:", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}
			json.NewEncoder(jsonPart).Encode(item)

			filePart, err := mw.CreateFormFile("file", item.FileName)
			if err != nil {
				log.Println("Error creating file part:", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}
			filePart.Write(data)
			mw.Close()
			return
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item type"})
		}
	}
}
