package controllers

import (
	"encoding/json"
	"github.com/dbsSensei/filesystem-api/service"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"

	"github.com/dbsSensei/filesystem-api/config"
	"github.com/dbsSensei/filesystem-api/database"
	"github.com/dbsSensei/filesystem-api/forms"
	"github.com/dbsSensei/filesystem-api/models"
	"github.com/dbsSensei/filesystem-api/utils"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	c  *config.Config
	db *gorm.DB
	s  *service.Services
}

func NewAuthController(config *config.Config, db *gorm.DB, s *service.Services) *AuthController {
	return &AuthController{
		c:  config,
		db: db,
		s:  s,
	}
}

// Signup godoc
// @Summary Signup user.
// @Description register user.
// @Tags Auth
// @Accept application/json
// @Param request body forms.SignupRequest true "request body"
// @Produce json
// @Success 200 {object} utils.Response{data=object}
// @Failure 400 {object} utils.Response{data=object}
// @Failure 500 {object} utils.Response{data=object}
// @Router /api/v1/auth/signup [post]
func (ac *AuthController) Signup(ctx *gin.Context) {
	var input forms.SignupRequest
	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	//if validRole != t	rue {
	//	ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", "invalid user role", nil))
	//	return
	//}

	createUserTransaction := func(tx *gorm.DB) error {
		hashedPassword, _ := utils.HashPassword(input.Password)
		input.Password = hashedPassword

		_, err := ac.s.UserService.Create(&models.User{
			Name:     input.Name,
			Email:    input.Email,
			Password: input.Password,
			Status:   models.UserStatusPending,
		}, tx)
		if err != nil {
			return err
		}

		return nil
	}

	if err := utils.Transaction(database.GetDB(), createUserTransaction); err != nil {
		if err == gorm.ErrDuplicatedKey {
			ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", "email already use", nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}

	ctx.JSON(http.StatusCreated, utils.ResponseData("success", "success create user", nil))
}

// Signin godoc
// @Summary Login user.
// @Description login user with credentials.
// @Tags Auth
// @Accept application/json
// @Param request body forms.SigninRequest true "request body"
// @Produce json
// @Success 200 {object} utils.Response{data=forms.SigninResponse}
// @Failure 400 {object} utils.Response{data=object}
// @Failure 500 {object} utils.Response{data=object}
// @Router /api/v1/auth/signin [post]
func (ac *AuthController) Signin(c *gin.Context) {
	cfg, _ := config.LoadConfig(".")
	tokenMaker, _ := utils.NewJWTMaker(cfg.TokenSymmetricKey)

	var input forms.SigninRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseData("error", err.Error(), nil))
		return
	}

	findUserWithEmailQuery := func(query *gorm.DB) *gorm.DB {
		pageNum := 1
		pageSize := 1
		query.Where("LOWER(email) = ?", strings.ToLower(input.Email))
		query = query.Limit(pageSize).Offset((pageNum - 1) * pageSize)
		return query
	}

	results, err := ac.s.UserService.FindAll(findUserWithEmailQuery, nil)
	if err != nil {
		return
	}

	result, err := json.Marshal(results[0])
	if err != nil {
		return
	}

	var user models.User
	err = json.Unmarshal(result, &user)
	if err != nil {
		return
	}

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, utils.ResponseData("error", "invalid email or password", nil))
		return
	}

	//if user.Status == models.UserStatusPending {
	//	c.JSON(http.StatusBadRequest, utils.ResponseData("error", "please verify your account", nil))
	//	return
	//}

	err = utils.CheckPassword(input.Password, user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ResponseData("error", "invalid email or password", nil))
		return
	}

	accessToken, accessPayload, _ := tokenMaker.CreateToken(
		user.ID,
		cfg.AccessTokenDuration,
	)

	refreshToken, refreshPayload, err := tokenMaker.CreateToken(
		user.ID,
		cfg.RefreshTokenDuration,
	)

	_, err = ac.s.TokenService.Create(&models.Token{
		ID:           refreshPayload.Id,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		IsBlocked:    false,
		ExpiresAt:    time.Time{},
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	}, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
	}

	rsp := forms.SigninResponse{
		SessionID:             refreshPayload.Id,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
	}

	c.JSON(http.StatusCreated, utils.ResponseData("success", "success signin user", rsp))
}
