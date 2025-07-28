package authhandelers

import (
	"github.com/darkard2003/wormhole/internals/server/middleware"
	"github.com/darkard2003/wormhole/internals/services/db"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.RouterGroup, db db.DBInterface) {
	r.POST("/signup", SignUpHandlerHandler(db))
	r.POST("/signin", SignInHandlerHandler(db))

	r.GET("/refresh", middleware.RefreshMiddleware(), RefreshHandler(db))
}
