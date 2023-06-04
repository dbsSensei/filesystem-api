package controllers

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/dbsSensei/filesystem-api/config"
	"github.com/dbsSensei/filesystem-api/forms"
	"github.com/dbsSensei/filesystem-api/models"
	"github.com/dbsSensei/filesystem-api/service"
	"github.com/dbsSensei/filesystem-api/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type FilesystemController struct {
	config *config.Config
	db     *gorm.DB
	s      *service.Services
}

func NewFilesystemController(config *config.Config, db *gorm.DB, s *service.Services) *FilesystemController {
	return &FilesystemController{
		config: config,
		db:     db,
		s:      s,
	}
}

// Download godoc
// @Summary Download a compressed file
// @Description Downloads a tar.gz file for processing
// @Tags Files
// @Accept */*
// @Produce application/file
// @Success 200 {object} utils.Response
// @Failure 500 {object} utils.Response{data=object}
// @Param filename path string true "file you want to download"
// @Router /api/v1/filesystem/download/{filename} [get]
func (f *FilesystemController) Download(ctx *gin.Context) {
	filename := ctx.Param("filename")

	// Check if the file exists in the extracted folder
	extractedFolder := "./extracted"
	filePath := filepath.Join(extractedFolder, filename)
	_, err := os.Stat(filePath)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.ResponseData("error", "File not found", nil))
		return
	}

	// Set the appropriate headers for the file download
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.File(filePath)
}

// MyFiles godoc
// @Summary Show logged-in user files.
// @Description get all logged-in user files.
// @Tags Files
// @Accept */*
// @Produce json
// @Success 200 {object} utils.Response{data=models.Filesystem}
// @Failure 500 {object} utils.Response{data=object}
// @Security ApiKeyAuth
// @Param page query int false  "files page"
// @Param limit query int false  "limit per files"
// @Param order_by query string false  "order files by"
// @Router /api/v1/filesystem/my-files [get]
func (ac *UserController) MyFiles(ctx *gin.Context) {
	authPayload := ctx.MustGet("authorization_payload").(*utils.TokenPayload)
	pageNum, _ := strconv.Atoi(ctx.Query("page"))
	if pageNum == 0 {
		pageNum = 1
	}

	pageSize, _ := strconv.Atoi(ctx.Query("limit"))
	if pageSize == 0 {
		pageSize = 10
	}

	filesFilterAndSort := func(query *gorm.DB) *gorm.DB {
		query.Where("user_id = ?", authPayload.UserId)

		queryOrder := ctx.Query("order_by")
		switch queryOrder {
		case "newest":
			query.Order("created_at desc")
		case "oldest":
			query.Order("created_at asc")
		default:
			query.Order("created_at desc")
		}
		return query
	}
	results, pagination, err := ac.s.FilesystemService.FindAll(pageNum, pageSize, filesFilterAndSort, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}

	var files []models.Filesystem
	for _, result := range results {
		var file models.Filesystem
		resultJson, err := json.Marshal(result)
		if err != nil {
			return
		}

		err = json.Unmarshal(resultJson, &file)
		if err != nil {
			return
		}
		files = append(files, file)
	}
	ctx.JSON(http.StatusOK, utils.ResponseData("success", "success get current user files", forms.GetMyFilesResponse{
		Files:      files,
		Pagination: pagination,
	}))
}

// Upload godoc
// @Summary Upload a compressed file
// @Description Uploads a tar.gz file for processing
// @Tags Files
// @Accept multipart/form-data
// @Produce application/json
// @Param file formData file true "The tar.gz file to upload"
// @Success 200 {object} utils.Response
// @Failure 500 {object} utils.Response{data=object}
// @Security ApiKeyAuth
// @Router /api/v1/filesystem/upload [post]
func (f *FilesystemController) Upload(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ResponseData("error", "Failed to retrieve the file", nil))
		return
	}

	// Create a temporary folder to extract the file contents
	tempFolder := "./temp"
	if err := os.MkdirAll(tempFolder, os.ModePerm); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", "Failed to create temporary folder", nil))
		return
	}

	// Save the uploaded file to the temporary folder
	filePath := filepath.Join(tempFolder, file.Filename)
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", "Failed to save the file", nil))
		return
	}

	// Extract the file contents
	extractedFolder := "./extracted"
	if err := extractFile(ctx, f.s, filePath, extractedFolder); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", "Failed to extract the file", nil))
		return
	}

	// Read the contents of the extracted files
	fileContents, err := readExtractedFiles(extractedFolder)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", "Failed to read the extracted files", nil))
		return
	}

	// Clean up the temporary and extracted folders
	os.RemoveAll(tempFolder)
	//os.RemoveAll(extractedFolder)

	ctx.JSON(http.StatusOK, utils.ResponseData("success", "Success extract and upload files", map[string]any{"uploaded_file": fileContents}))
}

func extractFile(ctx *gin.Context, s *service.Services, filePath, targetFolder string) error {
	authPayload := ctx.MustGet("authorization_payload").(*utils.TokenPayload)

	// Open the compressed file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a gzip reader to read the compressed file
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	// Create a tar reader to read the contents of the compressed file
	tarReader := tar.NewReader(gzipReader)

	// Iterate over each file in the tar archive
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			// End of archive
			break
		}
		if err != nil {
			return err
		}

		// Ensure the file is a regular file (not a directory or symbolic link)
		if header.Typeflag != tar.TypeReg {
			continue
		}

		filename := fmt.Sprintf("%v-%v-%v", authPayload.UserId, time.Now().UnixMilli(), header.Name)

		// Extract the file to the target folder
		filePath := filepath.Join(targetFolder, filename)
		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
		if err != nil {
			return err
		}

		if _, err := io.Copy(file, tarReader); err != nil {
			file.Close()
			return err
		}

		file.Close()

		_, _ = s.FilesystemService.Create(&models.Filesystem{
			UserID: authPayload.UserId,
			Name:   filename,
		}, nil)

	}

	return nil
}

func readExtractedFiles(extractedFolder string) ([]string, error) {
	var fileContents []string

	// Traverse the extracted folder and read the contents of each file
	err := filepath.Walk(extractedFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			// Read the contents of the file
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			// TODO handle uploaded
			_, err = io.ReadAll(file)
			if err != nil {
				return err
			}

			// Add the file contents to the result
			fileContents = append(fileContents, info.Name())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return fileContents, nil
}
