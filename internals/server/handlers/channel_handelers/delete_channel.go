package channelhandelers

import (
	"net/http"
	"strconv"

	"github.com/darkard2003/wormhole/internals/services/db"
	"github.com/darkard2003/wormhole/utils"
	"github.com/gin-gonic/gin"
)

func DeleteChannelHandler(db db.DBInterface) gin.HandlerFunc {

	return func(c *gin.Context) {
		channelId := c.Query("id")
		userId, exists := c.Get("userId")
		userId, ok := userId.(int)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
			return
		}
		channelIdInt, err := strconv.Atoi(channelId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid channel ID"})
			return
		}
		err = db.DeleteChannel(channelIdInt, userId.(int))
		if err != nil {
			httpError := utils.DBToHttpError(err)
			c.JSON(httpError.Code, gin.H{"error": httpError.Response})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Channel deleted successfully"})
	}
}
