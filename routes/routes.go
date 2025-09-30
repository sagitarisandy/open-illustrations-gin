package routes

import (
	"open-illustrations-go/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")

	api.GET("/illustrations", controllers.GetIllustrations)
	api.GET("/illustrations/:id", controllers.GetIllustrations)
}
