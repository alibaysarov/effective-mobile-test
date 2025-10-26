package main

import (
	"effective-mobile/internal/controller"
	"effective-mobile/internal/logger"
	"effective-mobile/internal/service"
	"effective-mobile/internal/storage/postgres"
	"effective-mobile/internal/storage/postgres/repository"
	"fmt"
	"log"
	"os"

	_ "effective-mobile/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title EFFECTIVE MOBILE Subscription Service API
// @version 1.0
// @description API для сервиса подписок
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:9001
// @BasePath /api/v1
// @schemes http
func main() {
	//logger init
	file, err := os.OpenFile("info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	logger, err := logger.InitLogger("wallet-service", file)
	if err != nil {
		log.Fatal(err)
	}

	//db init
	dbUrl := os.Getenv("DATABASE_URL")

	db, err := postgres.DatabaseConnect(dbUrl)
	if err != nil {
		logger.Error("Error while connecting to database" + err.Error())
		log.Fatal(err)
	}
	fmt.Println("connected to db")

	defer db.Close()

	subscribeRepository := repository.NewSubscribeRepository(db)
	subscribeService := service.NewSubscribeService(subscribeRepository, logger)
	subscribeController := controller.NewSubscribeController(subscribeService)

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")
	{
		api.GET("/subscriptions", subscribeController.GetAll)
		api.POST("/subscriptions", subscribeController.Create)
		api.GET("/subscriptions/:id", subscribeController.GetOne)
		api.PUT("/subscriptions/:id", subscribeController.Update)
		api.DELETE("/subscriptions/:id", subscribeController.Delete)
		api.GET("/subscriptions/total", subscribeController.GetSum)
	}
	PORT := os.Getenv("APP_PORT")
	if PORT == "" {
		PORT = "9001"
		logger.Warn("APP_PORT not set, using default: 9001")
	}
	fmt.Printf("Starting server on port: %s\n", PORT)
	router.Run(":" + PORT)
}
