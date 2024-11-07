package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter initializes the Gin router with all routes.
func SetupRouter(database *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.POST("/books", UploadBook(database))
	router.GET("/books/:id/download", DownloadBook(database))

	return router
}
