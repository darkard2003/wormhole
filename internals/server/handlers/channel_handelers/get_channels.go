package channelhandelers

import (
	"net/http"
	"strconv"

	"github.com/darkard2003/wormhole/internals/models"
	"github.com/darkard2003/wormhole/internals/services/db"
	"github.com/darkard2003/wormhole/internals/utils"
	"github.com/gin-gonic/gin"
)

type SanitedChannel struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Protected   bool   `json:"protected"`
}

func GetChannelHandeler(db db.DBInterface) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		userId, exists := ctx.Get("userId")
		if !exists || userId == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		limit := ctx.DefaultQuery("limit", "10")
		offset := ctx.DefaultQuery("offset", "0")
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
			return
		}
		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset"})
			return
		}
		total, channels, err := db.GetChannelsByUserId(userId.(int), limitInt, offsetInt)
		if err != nil {
			httpError := utils.DBToHttpError(err)
			ctx.JSON(httpError.Code, gin.H{"error": httpError.Response})
			return
		}
		sanitizedChannels := make([]*SanitedChannel, len(channels))
		for i, channel := range channels {
			sanitizedChannels[i] = SanitizeChannel(channel)
		}
		ctx.JSON(http.StatusOK, gin.H{
			"total":    total,
			"limit":    limitInt,
			"offset":   offsetInt,
			"channels": sanitizedChannels,
		})
	}
}

func SanitizeChannel(channel *models.Channel) *SanitedChannel {
	return &SanitedChannel{
		ID:          channel.ID,
		Name:        channel.Name,
		Description: channel.Description,
		Protected:   channel.Protected,
	}
}
