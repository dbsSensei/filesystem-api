package routers

import (
	"github.com/dbsSensei/filesystem-api/config"
	"github.com/dbsSensei/filesystem-api/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// V1 Routes Docs
// @title Filesystem API
// @version 1.0
// @description This server provides the Filesystem API needs.
// @termsOfService http://swagger.io/terms/
// @contact.name Dimas Bagus Susilo
// @contact.url http://www.linkedin.com/in/dimasbagussusilo
// @contact.email dimasbagussusilo@gmail.com
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func V1(router *gin.Engine, c *config.Config, db *gorm.DB, s *service.Services) *gin.Engine {

	return router
}
