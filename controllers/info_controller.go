package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func About(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"about": "This is Open Illustrations built with Go Gin & MinIO.",
	})
}

func License(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"license": "All illustrations are free to use under the MIT license.",
	})
}
