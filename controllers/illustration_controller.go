package controllers

import (
	"log"
	"net/http"
	"strings"
	"time"

	"open-illustrations-go/models"
	"open-illustrations-go/services"

	"github.com/gin-gonic/gin"
)

type CreateIllustrationDTO struct {
	Title    string `json:"title" binding:"required"`
	Category string `json:"category" binding:"required"`
	FileName string `json:"file_name" binding:"required"`
}

func GetIllustrations(c *gin.Context) {
	data, err := services.GetIllustrations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func GetIllustration(c *gin.Context) {
	id := c.Param("id")
	ill, err := services.GetIllustration(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": ill})
}

// processUpload handles multipart form upload: fields => file, title, category, file_name(optional)
func processUpload(c *gin.Context) {
	// fHeader, err := c.FormFile("file")
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
	// 	return
	// }
	// title := c.PostForm("title")
	// category := c.PostForm("category")
	// objectName := c.PostForm("file_name")
	// if objectName == "" {
	// 	objectName = fHeader.Filename
	// }
	// if title == "" || category == "" {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "title and category are required"})
	// 	return
	// }

	// file, err := fHeader.Open()
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
	// 	return
	// }

	// defer file.Close()

	// if err := services.UploadObject(objectName, file, fHeader.Size, fHeader.Header.Get("Content-Type")); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload to MinIO"})
	// 	return
	// }

	// rec := models.Illustration{
	// 	Title:    title,
	// 	Category: category,
	// 	FileName: objectName,
	// }
	// if err := services.CreateIllustrationRecord(&rec); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save record"})
	// 	return
	// }

	// c.JSON(http.StatusCreated, gin.H{"data": rec})

	// form-data: file (File), title (Text), category (Text), file_name (Text, optional)
	fh, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "form field 'file' is required"})
		return
	}
	title := c.PostForm("title")
	category := c.PostForm("category")
	objectName := c.PostForm("file_name")
	if objectName == "" {
		objectName = fh.Filename
	}
	if title == "" || category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "form fields 'title' and 'category' are required"})
		return
	}

	f, err := fh.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open uploaded file"})
		return
	}
	defer f.Close()

	// Jika sudah ada nama object yang sama -> tolak (hindari duplikat)
	exists, err := services.MinioObjectExists(objectName)
	if err != nil {
		log.Println("check minio exists err:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "minio check failed"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "file already exists in bucket", "file_name": objectName})
		return
	}

	// Simpan metadata ke MySQL
	// rec := models.Illustration{
	// 	Title:    title,
	// 	Category: category,
	// 	FileName: objectName,
	// }
	// if err := services.CreateIllustration(&rec); err != nil {
	// 	log.Println("db insert err:", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save record"})
	// 	return
	// }

	// c.JSON(http.StatusCreated, gin.H{"data": rec})

	// upload ke MinIO
	if err := services.UploadObject(objectName, f, fh.Size, fh.Header.Get("Content-Type")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload to storage"})
		return
	}

	rec := models.Illustration{Title: title, Category: category, FileName: objectName}
	if err := services.CreateIllustration(&rec); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save record"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": rec})
}

// Deprecated path: POST /illustrations/upload (still works). Prefer using POST /illustrations with multipart form-data.
func UploadIllustration(c *gin.Context) {
	processUpload(c)
}

func CreateIllustration(c *gin.Context) {
	// If client sent multipart form (file upload), reuse upload logic here so
	// people can just POST /illustrations with form-data.
	if ct := c.GetHeader("Content-Type"); strings.Contains(ct, "multipart/form-data") {
		processUpload(c)
		return
	}

	var dto CreateIllustrationDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload (make sure Content-Type: application/json)"})
		return
	}
	// c.JSON(http.StatusCreated, gin.H{"data": input})

	input := models.Illustration{
		Title:    dto.Title,
		Category: dto.Category,
		FileName: dto.FileName,
	}

	if err := services.CreateIllustration(&input); err != nil {
		log.Println("CreateIllustration DB/MINIO err:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": input})
}

func DeleteIllustration(c *gin.Context) {
	id := c.Param("id")
	if err := services.DeleteIllustration(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete illustration"})
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func Download(c *gin.Context) {
	id := c.Param("id")
	ill, err := services.GetIllustration(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	url, err := services.GetDownloadURL(ill.FileName, time.Hour*1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate link"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"download_url": url})
}
