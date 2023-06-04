package forms

import (
	"github.com/dbsSensei/filesystem-api/models"
	"github.com/dbsSensei/filesystem-api/utils"
)

type GetMyFilesResponse struct {
	Files      []models.Filesystem `json:"files"`
	Pagination utils.Pagination    `json:"pagination"`
}
