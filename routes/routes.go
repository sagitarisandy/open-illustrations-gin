package routes

import (
	"open-illustrations-go/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")

	api.GET("/illustrations", controllers.GetIllustrations)
	api.POST("/illustrations/upload", controllers.UploadIllustration)

	// Penting: letakkan sebelum /illustrations/:id agar tidak tertutup wildcard
	// api.GET("/illustrations/file/:key", controllers.GetIllustrationFileURL)

	api.GET("/illustrations/:id/file", controllers.GetIllustrationFileURLByID)

	api.GET("/illustrations/:id", controllers.GetIllustration)
	api.POST("/illustrations", controllers.CreateIllustration)
	api.DELETE("/illustrations/:id", controllers.DeleteIllustration)
	api.GET("/illustrations/:id/download", controllers.Download)

	// Public stream for non-premium assets by ID
	api.GET("/illustrations/:id/public", controllers.StreamPublic)

	// Asset streaming via signed token path
	api.GET("/i/:token", controllers.StreamSigned)

	api.POST("/category", controllers.CreateCategory)
	api.GET("/categories", controllers.GetCategories)
	api.GET("/categories/:id", controllers.GetCategory)
	api.GET("/categories/:id/illustrations", controllers.GetIllustrationsByCategory)
	api.PUT("/categories/:id", controllers.DeleteCategory)

	api.POST("/pack", controllers.CreatePack)
	api.GET("/packs", controllers.GetPacks)
	api.GET("/packs/:id", controllers.GetPack)
	api.GET("/packs/:id/illustrations", controllers.GetIllustrationsByPack)
	api.PUT("/packs/:id", controllers.DeletePack)
	api.GET("/packs/:id/download", controllers.DownloadPacks)

	api.POST("/styles", controllers.CreateStyle)
	api.GET("/styles", controllers.GetStyles)
	api.GET("/styles/:id", controllers.GetStyle)
	api.GET("/styles/:id/illustrations", controllers.GetIllustrationsByStyle)
	api.PUT("/styles/:id", controllers.UpdateStyle)
	api.DELETE("/styles/:id", controllers.DeleteStyle)

	api.GET("/info/about", controllers.About)
	api.GET("/info/license", controllers.License)
}
