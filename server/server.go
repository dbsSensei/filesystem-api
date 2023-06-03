package server

import (
	"fmt"
	"github.com/dbsSensei/filesystem-api/config"
	"github.com/dbsSensei/filesystem-api/controllers"
	"github.com/dbsSensei/filesystem-api/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Init(c *config.Config, db *gorm.DB) error {
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

	// Setup Static File
	server.Static("/public", "./public")

	// Run
	err := server.Run(c.HTTPServerAddress)
	if err != nil {
		return fmt.Errorf("error while running server %+e", err)
	}

	return nil
}
