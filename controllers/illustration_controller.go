package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"open-illustrations-go/models"
	"open-illustrations-go/services"

	"github.com/gin-gonic/gin"
)

type CreateIllustrationDTO struct {
	Title      string `json:"title" binding:"required"`
	CategoryID *uint  `json:"category_id"`
	PackID     *uint  `json:"pack_id"`
	FileName   string `json:"file_name" binding:"required"`
	StorageKey string `json:"storage_key"`
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

func GetIllustrationURL(c *gin.Context) {
	id := c.Param("id")
	ill, err := services.GetIllustration(id)
	if err != nil || ill == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "illustration not found"})
		return
	}

	//Allow custom expires in seconds with cap
	expSecStr := c.Query("expires")
	exp := 15 * time.Minute
	if expSecStr != "" {
		if n, convErr := strconv.Atoi(expSecStr); convErr == nil && n > 0 {
			if n > 3600 {
				n = 3600 // cap 1h
			}
			exp = time.Duration(n) * time.Second
		}
	}

	url, err := services.GetDownloadURL(ill.StorageKey, exp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate url"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          ill.ID,
		"file_name":   ill.FileName,
		"storage_key": ill.StorageKey,
		"url":         url,
		"expires_in":  int(exp.Seconds()),
	})

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
	// category legacy removed
	catIDStr := c.PostForm("category_id")
	packIDStr := c.PostForm("pack_id")
	objectName := c.PostForm("file_name")
	if objectName == "" {
		objectName = fh.Filename
	}
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "form field 'title' is required"})
		return
	}

	f, err := fh.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open uploaded file"})
		return
	}
	defer f.Close()

	// Validate SVG only (by extension + light content-type check)
	if !isSVGFile(fh) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only .svg files are allowed"})
		return
	}

	// Generate unique storage key (random hex) while preserving original filename separately
	storageKey := generateStorageKey(objectName)

	// Jika sudah ada nama object yang sama -> tolak (hindari duplikat)
	exists, err := services.MinioObjectExists(storageKey)
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
	if err := services.UploadObject(storageKey, f, fh.Size, fh.Header.Get("Content-Type")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload to storage", "detail": err.Error()})
		return
	}

	var catIDPtr *uint
	var packIDPtr *uint
	if catIDStr != "" {
		if v, err := strconv.ParseUint(catIDStr, 10, 64); err == nil {
			vv := uint(v)
			catIDPtr = &vv
		}
	}
	if packIDStr != "" {
		if v, err := strconv.ParseUint(packIDStr, 10, 64); err == nil {
			vv := uint(v)
			packIDPtr = &vv
		}
	}
	rec := models.Illustration{Title: title, CategoryID: catIDPtr, PackID: packIDPtr, FileName: objectName, StorageKey: storageKey}
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

	storageKey := dto.StorageKey
	if storageKey == "" {
		// For JSON-based creation we expect the storage_key (object already uploaded via another service)
		c.JSON(http.StatusBadRequest, gin.H{"error": "storage_key is required when creating via JSON"})
		return
	}
	input := models.Illustration{
		Title:      dto.Title,
		CategoryID: dto.CategoryID,
		PackID:     dto.PackID,
		FileName:   dto.FileName,
		StorageKey: storageKey,
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

	url, err := services.GetDownloadURL(ill.StorageKey, time.Hour*1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate link"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"download_url": url})
}

// --- helpers ---
func isSVGFile(fh *multipart.FileHeader) bool {
	name := strings.ToLower(fh.Filename)
	if !strings.HasSuffix(name, ".svg") {
		return false
	}
	f, err := fh.Open()
	if err != nil {
		return false
	}
	defer f.Close()
	buf := make([]byte, 512)
	n, _ := f.Read(buf)
	snippet := strings.ToLower(string(buf[:n]))
	if !strings.Contains(snippet, "<svg") {
		return false
	}
	return true
}

func generateStorageKey(original string) string {
	// random 8 bytes -> 16 hex chars + original ext
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return original // fallback
	}
	ext := filepath.Ext(original)
	if ext == "" {
		ext = ".svg" // default
	}
	return fmt.Sprintf("%s-%s%s", time.Now().Format("20060102"), hex.EncodeToString(b), ext)
}
