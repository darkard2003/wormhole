package itemhandelers

import (
	"github.com/darkard2003/wormhole/services/db"
	"github.com/gin-gonic/gin"
)

func RegisterItemRoutes(router *gin.RouterGroup, db db.DBInterface) {
	router.POST("/items", PushItem(db))
	router.GET("/items", PopItem(db))
}
