package routers

import (
	"github.com/dbsSensei/filesystem-api/config"
	"github.com/dbsSensei/filesystem-api/controllers"
	"github.com/dbsSensei/filesystem-api/middlewares"
	"github.com/dbsSensei/filesystem-api/service"
	"github.com/dbsSensei/filesystem-api/utils"
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
	//////////
	// Public
	v1 := router.Group("api/v1")

	// Auth
	authEndpoint := "/auth"
	auth := controllers.NewAuthController(c, db, s)
	v1.POST(authEndpoint+"/signin", auth.Signin)
	v1.POST(authEndpoint+"/signup", auth.Signup)
	//v1.POST(authEndpoint+"/verify", auth.Verify)

	// User
	usersEndpoint := "/users"
	users := controllers.NewUserController(c, db, s)

	//////////////
	// Authorized
	tokenMaker, _ := utils.NewJWTMaker(c.TokenSymmetricKey)
	authorizedV1 := router.Group("api/v1").Use(middlewares.AuthMiddleware(tokenMaker))
	authorizedV1.GET(usersEndpoint+"/whoami", users.WhoAmI)

	return router
}
