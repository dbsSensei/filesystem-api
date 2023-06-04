package controllers

import (
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
// @Security ApiKeyAuth
// @Router /api/v1/users/whoami [get]
func (ac *UserController) WhoAmI(ctx *gin.Context) {
	authPayload := ctx.MustGet("authorization_payload").(*utils.TokenPayload)
	result, err := ac.s.UserService.FindOne(authPayload.UserId, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}
	user := result.(*models.User)

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "success get current user", forms.WhoAmIResponse{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Status: string(user.Status),
	}))
}
