package channelhandelers

import (
	"github.com/darkard2003/wormhole/interfaces"
	"github.com/gin-gonic/gin"
)

func RegisterChannelRoutes(r *gin.RouterGroup, db interfaces.DBInterface) {
	r.GET("/channels", GetChannelHandeler(db))
	r.POST("/channels", CreateChannelHandler(db))
	r.DELETE("/channels", DeleteChannelHandler(db))

}
