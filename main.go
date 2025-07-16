package main

import (
	"log"
	"os"

	"github.com/darkard2003/wormhole/handlers"
	"github.com/darkard2003/wormhole/handlers/auth_handelers"
	channelhandelers "github.com/darkard2003/wormhole/handlers/channel_handelers"
	"github.com/darkard2003/wormhole/middleware"
	"github.com/darkard2003/wormhole/services/db"
	"github.com/darkard2003/wormhole/services/envservice"
	"github.com/darkard2003/wormhole/services/mysqldb"
	"github.com/gin-gonic/gin"
)

var appDb db.DBInterface

func init() {
	envservice.LoadEnv()
	log.Println("Environment variables loaded successfully")
	appDb = &mysqldb.MySqlRepo{}
	err := appDb.Initialize()
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
	handlers.RegisterAPIRoutes(apiv1)

	authhandelers.RegisterAuthRoutes(apiv1, appDb)
	authenticatedRoute := apiv1.Group("/user")
	authenticatedRoute.Use(middleware.AuthMiddleware())
	authenticatedRoute.GET("/status", authhandelers.AuthStatus)
	channelhandelers.RegisterChannelRoutes(authenticatedRoute, appDb)

	r.Run()
}
