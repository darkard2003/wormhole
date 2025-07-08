package channelhandelers

import (
	"github.com/darkard2003/wormhole/models"
	"github.com/darkard2003/wormhole/services/dbservice"
	"github.com/gin-gonic/gin"
)

type SanitedChannel struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Protected   bool   `json:"protected"`
}

func GetChannels(ctx *gin.Context) {
	userId, exists := ctx.Get("userId")
	if !exists || userId == nil {
		ctx.JSON(400, gin.H{"error": "Unauthorized"})
		return
	}
	channels, err := dbservice.GetDBService().GetChannelsByUserId(userId.(int))
	sanitizedChannels := make([]*SanitedChannel, len(channels))
	for i, channel := range channels {
		sanitizedChannels[i] = SanitizeChannel(channel)
	}
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to retrieve channels"})
		return
	}
	ctx.JSON(200, gin.H{
		"total":    len(channels),
		"channels": sanitizedChannels,
	})

}

func SanitizeChannel(channel *models.Channel) *SanitedChannel {
	return &SanitedChannel{
		ID:          channel.ID,
		Name:        channel.Name,
		Description: channel.Description,
		Protected:   channel.Protected,
	}
}
