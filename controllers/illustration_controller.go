package controllers

import (
	"net/http"
	"time"

	"open-illustrations-go/models"
	"open-illustrations-go/services"

	"github.com/gin-gonic/gin"
)

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

func CreateIllustration(c *gin.Context) {
	var input models.Illustration
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create"})
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
