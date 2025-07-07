package main

import (
	"fmt"
	"os"

	"github.com/darkard2003/wormhole/handelers"
	"github.com/darkard2003/wormhole/services/dbservice"
	"github.com/darkard2003/wormhole/services/envservice"
	"github.com/gin-gonic/gin"
)

func init() {
	envservice.LoadEnv()
	fmt.Println("Environment variables loaded successfully")
	db := dbservice.GetDBService()
	err := db.InitializeMySql()
	if err != nil {
		fmt.Println("Error initializing database:", err)
		os.Exit(1)
	}
	fmt.Println("Database initialized successfully")
}

func main() {
	r := gin.Default()

	r.GET("/", handelers.IndexHandler)

	apiv1 := r.Group("/api/v1")
	apiv1.GET("/ping", handelers.PingHandler)
	apiv1.GET("/health", handelers.HealthCheckHandler)

	apiv1.POST("/sign_up", handelers.SignUpHandeler)
	apiv1.POST("/sign_in", handelers.SignInHandeler)

	r.Run()
}
