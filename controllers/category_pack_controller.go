package controllers

import (
	"archive/zip"
	"context"
	"io"
	"net/http"
	"time"

	"open-illustrations-go/config"
	"open-illustrations-go/models"
	"open-illustrations-go/services"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

type createNamedDTO struct {
	Name string `json:"name" binding:"required"`
}

// ---- Category Handlers ----
func CreateCategory(c *gin.Context) {
	var dto createNamedDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	cat, err := services.CreateCategory(dto.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": cat})
}

func GetCategories(c *gin.Context) {
	cats, err := services.GetCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": cats})
}

func GetCategory(c *gin.Context) {
	cat, err := services.GetCategory(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": cat})
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	cat, err := services.SoftDeleteCategory(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ts := ""
	if cat.DeletedAt.Valid {
		ts = cat.DeletedAt.Time.Format(time.RFC3339)
	}
	c.JSON(http.StatusOK, gin.H{"id": cat.ID, "deleted_at": ts})
}

// ---- Pack Handlers ----
func CreatePack(c *gin.Context) {
	var dto createNamedDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	p, err := services.CreatePack(dto.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": p})
}

func GetPacks(c *gin.Context) {
	list, err := services.GetPacks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func GetPack(c *gin.Context) {
	p, err := services.GetPack(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": p})
}

func DeletePack(c *gin.Context) {
	id := c.Param("id")
	p, err := services.SoftDeletePack(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ts := ""
	if p.DeletedAt.Valid {
		ts = p.DeletedAt.Time.Format(time.RFC3339)
	}
	c.JSON(http.StatusOK, gin.H{"id": p.ID, "deleted_at": ts})
}

// DownloadPacks: create a zip stream with all illustration SVGs in a pack.
func DownloadPacks(c *gin.Context) {
	pack, err := services.GetPack(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "pack not found"})
		return
	}
	// Load illustrations for this pack
	var ills []models.Illustration
	if err := config.DB.Where("pack_id = ? AND deleted_at IS NULL", pack.ID).Find(&ills).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	filename := services.PackArchiveFileName(pack)
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	w := zip.NewWriter(c.Writer)
	for _, ill := range ills {
		// fetch object from MinIO
		obj, err := config.MinioClient.GetObject(context.Background(), config.BucketName, ill.StorageKey, minio.GetObjectOptions{})
		if err != nil {
			continue
		}
		f, err := w.Create(ill.FileName)
		if err != nil {
			continue
		}
		io.Copy(f, obj)
		obj.Close()
	}
	w.Close()
}
