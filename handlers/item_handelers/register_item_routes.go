package itemhandelers

import (
	"github.com/darkard2003/wormhole/services/db"
	storageservice "github.com/darkard2003/wormhole/services/storage_service"
	"github.com/gin-gonic/gin"
)

func RegisterItemRoutes(router *gin.RouterGroup, db db.DBInterface, s storageservice.StorageInterface) {
	router.POST("/items", PushItem(db, s))
	router.GET("/items", PopItem(db, s))
}
