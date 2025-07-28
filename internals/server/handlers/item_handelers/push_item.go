package itemhandelers

import (
	"io"
	"log"
	"net/http"

	"github.com/darkard2003/wormhole/internals/models"
	"github.com/darkard2003/wormhole/internals/services/db"
	storageservice "github.com/darkard2003/wormhole/internals/services/storage_service"
	"github.com/darkard2003/wormhole/internals/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type PushItemRequest struct {
	ChannelName     string `form:"channel_name"`
	ChannelPassword string `form:"channel_password"`
	Type            string `form:"type" binding:"required"`
	Title           string `form:"title" binding:"required"`
	Salt            string `form:"salt" binding:"required"`
	IV              string `form:"iv" binding:"required"`
	Encryption      string `form:"encryption" binding:"required"`
	TextContent     string `form:"text_content"`
	Filename        string `form:"filename"`
	MimeType        string `form:"mimetype"`
	FileSize        int64  `form:"filesize"`
	FileCreatedAt   string `form:"file_created_at"`
	FileUpdatedAt   string `form:"file_updated_at"`
}

func PushItem(db db.DBInterface, s storageservice.StorageInterface) gin.HandlerFunc {
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
		if err := ctx.ShouldBind(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

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
			fileheader, err := ctx.FormFile("file")
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
				return
			}

			file, err := fileheader.Open()
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
				return
			}
			defer file.Close()

			data, err := io.ReadAll(file)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
				return
			}

			id, err := s.StoreBlob(data)
			if err != nil {
				log.Println("Error storing blob:", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}

			item = &models.FileItem{
				Item: models.Item{
					UserID:             userIdInt,
					ChannelID:          channel.ID,
					Type:               request.Type,
					Title:              request.Title,
					Salt:               request.Salt,
					IV:                 request.IV,
					EncryptionMetadata: request.Encryption,
				},
				FileName:      request.Filename,
				FileSize:      request.FileSize,
				BlobSize:      fileheader.Size,
				MimeType:      request.MimeType,
				FileCreatedAt: request.FileCreatedAt,
				FileUpdatedAt: request.FileUpdatedAt,
				FileId:        id,
			}
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item type"})
			return
		}

		_, err = db.CreateItem(item)
		if err != nil {
			httpError := utils.DBToHttpError(err)
			ctx.JSON(httpError.Code, gin.H{"error": httpError.Response})
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{"message": "Item created successfully"})
	}
}
