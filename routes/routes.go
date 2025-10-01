package routes

import (
	"open-illustrations-go/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")

	api.GET("/illustrations", controllers.GetIllustrations)
	api.GET("/illustrations/:id", controllers.GetIllustrations)
	api.POST("/illustrations", controllers.CreateIllustration)
	api.DELETE("/illustrations/:id", controllers.DeleteIllustration)
	api.GET("/illustrations/:id/download", controllers.Download)

	api.GET("/info/about", controllers.About)
	api.GET("/info/license", controllers.License)
}
