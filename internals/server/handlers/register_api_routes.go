package handlers

import "github.com/gin-gonic/gin"

func RegisterAPIRoutes(r *gin.RouterGroup) {
	r.GET("/ping", PingHandler)
	r.GET("/health", HealthCheckHandler)
}
