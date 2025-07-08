package handlers

import "github.com/gin-gonic/gin"

func HealthCheckHandler(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"status": "ok"})
}
