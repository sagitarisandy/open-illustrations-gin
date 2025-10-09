package controllers

import (
	"net/http"

	"open-illustrations-go/services"

	"github.com/gin-gonic/gin"
)

type createStyleDTO struct {
	Name string `json:"name" binding:"required"`
}

type updateStyleDTO struct {
	Name string `json:"name" binding:"required"`
}

func CreateStyle(c *gin.Context) {
	var dto createStyleDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s, err := services.CreateStyle(dto.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": s})
}

func GetStyles(c *gin.Context) {
	list, err := services.GetStyles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func GetStyle(c *gin.Context) {
	id := c.Param("id")
	s, err := services.GetStyle(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": s})
}

func UpdateStyle(c *gin.Context) {
	var dto updateStyleDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	s, err := services.UpdateStyle(id, dto.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": s})
}

func DeleteStyle(c *gin.Context) {
	id := c.Param("id")
	s, err := services.SoftDeleteStyle(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// return deleted_at through service struct json
	c.JSON(http.StatusOK, gin.H{"id": s.ID, "deleted_at": s.DeletedAt.Time})
}
