package channelhandelers

import (
	"github.com/darkard2003/wormhole/internals/services/db"
	"github.com/gin-gonic/gin"
)

func RegisterChannelRoutes(r *gin.RouterGroup, db db.DBInterface) {
	r.GET("/channels", GetChannelHandeler(db))
	r.POST("/channels", CreateChannelHandler(db))
	r.DELETE("/channels", DeleteChannelHandler(db))

}
