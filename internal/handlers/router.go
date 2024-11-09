package handlers

import (
	"github.com/BuddhiLW/DL-backend/internal/db"
	"github.com/gin-gonic/gin"
)

// SetupRouter initializes the Gin router with all routes.
func SetupRouter(queries *db.Queries) *gin.Engine {
	router := gin.Default()

	// Routes
	router.POST("/books", UploadBook(queries))
	router.GET("/books/:id/download", DownloadBook(queries))
	router.GET("/books/:id/cover", GetCoverImage(queries))

	return router
}
