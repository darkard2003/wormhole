package authhandelers

import (
	"log"
	"net/http"

	"github.com/darkard2003/wormhole/internals/services/db"
	"github.com/darkard2003/wormhole/internals/services/jwtservice"
	"github.com/gin-gonic/gin"
)

func RefreshHandler(db db.DBInterface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, exists := ctx.Get("userId")
		if !exists || userId == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		userIdInt, ok := userId.(int)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		userName, exists := ctx.Get("userName")
		if !exists || userName == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		userNameStr, ok := userName.(string)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})

			return
		}

		jwtService := jwtservice.GetJWTService()

		token, exp, err := jwtService.GenerateToken(userIdInt, userNameStr)

		if err != nil {
			log.Println("Error generating token:", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"token":      token,
			"expires_at": exp,
			"message":    "Access token refreshed successfully",
		})
	}
}
