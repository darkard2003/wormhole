package main

import (
	"log"
	"os"

	"github.com/darkard2003/wormhole/handlers"
	"github.com/darkard2003/wormhole/handlers/auth_handelers"
	channelhandelers "github.com/darkard2003/wormhole/handlers/channel_handelers"
	"github.com/darkard2003/wormhole/middleware"
	"github.com/darkard2003/wormhole/services/dbservice"
	"github.com/darkard2003/wormhole/services/envservice"
	"github.com/gin-gonic/gin"
)

func init() {
	envservice.LoadEnv()
	log.Println("Environment variables loaded successfully")
	db := dbservice.GetDBService()
	err := db.InitializeMySql()
	if err != nil {
		log.Println("Error initializing database:", err)
		os.Exit(1)
	}
	log.Println("Database initialized successfully")
}

func main() {
	r := gin.Default()

	r.GET("/", handlers.IndexHandler)

	apiv1 := r.Group("/api/v1")
	apiv1.GET("/ping", handlers.PingHandler)
	apiv1.GET("/health", handlers.HealthCheckHandler)

	apiv1.POST("/signup", authhandelers.SignUpHandler)
	apiv1.POST("/signin", authhandelers.SignInHandler)

	authenticatedRoute := apiv1.Group("/user")
	authenticatedRoute.Use(middleware.AuthMiddleware())
	authenticatedRoute.GET("/status", handlers.AuthStatus)

	authenticatedRoute.GET("/channels", channelhandelers.GetChannels)
	authenticatedRoute.POST("/channels", channelhandelers.CreateChannel)

	r.Run()
}
