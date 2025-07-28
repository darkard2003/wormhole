package main

import (
	"log"
	"os"

	"github.com/darkard2003/wormhole/internals/server/handlers"
	authhandelers "github.com/darkard2003/wormhole/internals/server/handlers/auth_handelers"
	channelhandelers "github.com/darkard2003/wormhole/internals/server/handlers/channel_handelers"
	itemhandelers "github.com/darkard2003/wormhole/internals/server/handlers/item_handelers"
	"github.com/darkard2003/wormhole/internals/server/middleware"
	"github.com/darkard2003/wormhole/internals/services/db"
	"github.com/darkard2003/wormhole/internals/services/db/mysqldb"
	"github.com/darkard2003/wormhole/internals/services/envservice"
	"github.com/darkard2003/wormhole/internals/services/jwtservice"
	storageservice "github.com/darkard2003/wormhole/internals/services/storage_service"
	localstorage "github.com/darkard2003/wormhole/internals/services/storage_service/local_storage"
	"github.com/gin-gonic/gin"
)

var appDb db.DBInterface
var storage storageservice.StorageInterface

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

	storePath, err := envservice.GetEnv("STORE_PATH")
	if err != nil {
		log.Println("Error getting store path:", err)
		os.Exit(1)
	}
	storage = localstorage.NewLocalStorage(storePath)
	log.Println("Storage service initialized successfully")

	accessSecrectKey, err := envservice.GetEnv("ACCESS_SECRET_KEY")
	if err != nil {
		log.Fatal("Error getting access secret key:", err)
	}
	refreshSecrectKey, err := envservice.GetEnv("REFRESH_SECRET_KEY")
	if err != nil {
		log.Fatal("Error getting refresh secret key:", err)
	}
	jwtservice.InitJWTService(accessSecrectKey, refreshSecrectKey)
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
	itemhandelers.RegisterItemRoutes(authenticatedRoute, appDb, storage)

	r.Run()
}
