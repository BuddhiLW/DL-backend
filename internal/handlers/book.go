package handlers

import (
	"bytes"
	"io"
	"net/http"

	"github.com/BuddhiLW/DL-backend/internal/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UploadBook handles POST /books to upload a new book.
func UploadBook(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.PostForm("title")
		author := c.PostForm("author")
		category := c.PostForm("category")

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
			return
		}

		fileContent, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
			return
		}
		defer fileContent.Close()

		var buffer bytes.Buffer
		io.Copy(&buffer, fileContent)

		book := db.Book{
			Title:    title,
			Author:   author,
			Category: category,
			Content:  buffer.Bytes(),
		}
		database.Create(&book)

		c.JSON(http.StatusCreated, gin.H{"id": book.ID, "message": "Book uploaded successfully"})
	}
}

// DownloadBook handles GET /books/:id/download to retrieve a book's content.
func DownloadBook(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var book db.Book
		id := c.Param("id")

		if err := database.First(&book, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}

		c.Header("Content-Disposition", "attachment; filename="+book.Title+".pdf")
		c.Data(http.StatusOK, "application/pdf", book.Content)
	}
}
