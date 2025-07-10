package authhandelers

import (
	"github.com/darkard2003/wormhole/interfaces"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.RouterGroup, db interfaces.DBInterface) {
	r.POST("/signup", SignUpHandlerHandler(db))
	r.POST("/signin", SignInHandlerHandler(db))
}
