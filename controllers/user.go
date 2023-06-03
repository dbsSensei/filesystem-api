package controllers

import (
	"fmt"
	"github.com/dbsSensei/filesystem-api/config"
	"github.com/dbsSensei/filesystem-api/forms"
	"github.com/dbsSensei/filesystem-api/models"
	"github.com/dbsSensei/filesystem-api/service"
	"github.com/dbsSensei/filesystem-api/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type UserController struct {
	c  *config.Config
	db *gorm.DB
	s  *service.Services
}

func NewUserController(config *config.Config, db *gorm.DB, s *service.Services) *UserController {
	return &UserController{
		c:  config,
		db: db,
		s:  s,
	}
}

// WhoAmI godoc
// @Summary Show logged-in user.
// @Description get logged-in user data.
// @Tags Users
// @Accept */*
// @Produce json
// @Success 200 {object} utils.Response{data=forms.WhoAmIResponse}
// @Failure 500 {object} utils.Response{data=object}
// @Router /api/v1/users/whoami [get]
func (ac *UserController) WhoAmI(ctx *gin.Context) {
	authPayload, _ := ctx.Get("authorization_payload")

	fmt.Printf("auth payload %+v\n", authPayload)
	result, err := ac.s.UserService.FindOne(1, nil)
	if err != nil {
		ctx.JSON(http.StatusCreated, utils.ResponseData("error", "failed get current user", nil))
		return
	}

	var user models.User
	user = result.(models.User)

	ctx.JSON(http.StatusCreated, utils.ResponseData("success", "success get current user", forms.WhoAmIResponse{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Status: string(user.Status),
	}))
}
