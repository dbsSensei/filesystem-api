package server

import (
	"fmt"
	"github.com/dbsSensei/filesystem-api/config"
	"github.com/dbsSensei/filesystem-api/controllers"
	_ "github.com/dbsSensei/filesystem-api/docs/v1"
	"github.com/dbsSensei/filesystem-api/routers"
	"github.com/dbsSensei/filesystem-api/service"
	"github.com/dbsSensei/filesystem-api/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func Init(c *config.Config, db *gorm.DB, s *service.Services) error {
	// Setup Router
	server := gin.New()
	server.Use(utils.Logger())
	server.Use(gin.Recovery())

	// Setup Cors
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*", "http://localhost"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Role", "Authorization"}
	server.Use(cors.New(corsConfig))

	// Health Check
	serverController := controllers.NewServerController(c, db)
	server.GET("/health", serverController.HealthCheck)

	// Setup swagger documentation
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// Setup Static File
	server.Static("/public", "./public")

	// Setup Routers
	routers.V1(server, c, db, s)

	// Run
	err := server.Run(c.HTTPServerAddress)
	if err != nil {
		return fmt.Errorf("error while running server %+e", err)
	}

	return nil
}
