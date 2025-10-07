package routes

import (
	"open-illustrations-go/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")

	api.GET("/illustrations", controllers.GetIllustrations)
	api.POST("/illustrations/upload", controllers.UploadIllustration)
	api.GET("/illustrations/:id", controllers.GetIllustration)
	api.POST("/illustrations", controllers.CreateIllustration)
	api.DELETE("/illustrations/:id", controllers.DeleteIllustration)
	api.GET("/illustrations/:id/download", controllers.Download)
	api.GET("/illustrations/:id/url", controllers.GetIllustrationURL)

	api.POST("/category", controllers.CreateCategory)
	api.GET("/categories", controllers.GetCategories)
	api.GET("/categories/:id", controllers.GetCategory)
	api.PUT("/categories/:id", controllers.DeleteCategory)

	api.POST("/pack", controllers.CreatePack)
	api.GET("/packs", controllers.GetPacks)
	api.GET("/packs/:id", controllers.GetPack)
	api.PUT("/packs/:id", controllers.DeletePack)
	api.GET("/packs/:id/download", controllers.DownloadPacks)

	api.GET("/info/about", controllers.About)
	api.GET("/info/license", controllers.License)
}
