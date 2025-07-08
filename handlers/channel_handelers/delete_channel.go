package channelhandelers

import (
	"strconv"

	"github.com/darkard2003/wormhole/services/dbservice"
	"github.com/gin-gonic/gin"
)

func DeleteChannel(c *gin.Context) {
	channelId := c.Query("id")
	userId, exists := c.Get("userId")
	userId, ok := userId.(int)
	if !ok {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}
	if !exists {
		c.JSON(400, gin.H{"error": "Unauthorized"})
		return
	}
	channelIdInt, err := strconv.Atoi(channelId)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid channel ID"})
		return
	}
	err = dbservice.GetDBService().DeleteChannel(channelIdInt, userId.(int))
	if err != nil {
		if err.Error() == "channel not found" {
			c.JSON(404, gin.H{"error": "Channel not found"})
			return
		}
		c.JSON(500, gin.H{"error": "Failed to delete channel"})
		return
	}
	c.JSON(200, gin.H{"message": "Channel deleted successfully"})
}
