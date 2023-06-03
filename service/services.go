package service

import (
	"github.com/dbsSensei/filesystem-api/models"
	"gorm.io/gorm"
)

type Services struct {
	UserService       IRepository
	FilesystemService IRepository
}

func Init(db *gorm.DB) *Services {
	return &Services{
		UserService:       NewRepository(&models.User{}, db),
		FilesystemService: NewRepository(&models.Filesystem{}, db),
	}
}
