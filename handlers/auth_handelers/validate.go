package authhandelers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthStatus(ctx *gin.Context) {
	userId, err := ctx.Get("userId")
	if !err {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	if userId == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Authorized",
		"userId":  userId,
	})
}
