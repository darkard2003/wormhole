package handlers

import "github.com/gin-gonic/gin"

func IndexHandler(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"message": "Welcome to the Wormhole"})
}
