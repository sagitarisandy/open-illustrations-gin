package controllers

import (
	"net/http"

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
