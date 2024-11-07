package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Book represents a book in the library with its PDF content.
type Book struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Category string `json:"category"`
	Content  []byte `json:"-"` // PDF content stored as a BLOB
}

// Initialize the database and return a connection.
func initDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("library.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&Book{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	return db
}

func main() {
	// Initialize Gin router
	router := gin.Default()

	// Initialize database
	db := initDatabase()

	// Upload a new book (metadata + PDF file)
	router.POST("/books", func(c *gin.Context) {
		title := c.PostForm("title")
		author := c.PostForm("author")
		category := c.PostForm("category")

		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
			return
		}

		// Open the uploaded file
		fileContent, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
			return
		}
		defer fileContent.Close()

		// Read file content
		var buffer bytes.Buffer
		_, err = io.Copy(&buffer, fileContent)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process file"})
			return
		}

		// Create and save the book
		book := Book{
			Title:    title,
			Author:   author,
			Category: category,
			Content:  buffer.Bytes(),
		}
		result := db.Create(&book)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save book"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Book uploaded successfully", "id": book.ID})
	})

	// Download a book's PDF content
	router.GET("/books/:id/download", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
			return
		}

		var book Book
		result := db.First(&book, id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}

		// Serve the PDF file
		c.Header("Content-Disposition", "attachment; filename="+book.Title+".pdf")
		c.Data(http.StatusOK, "application/pdf", book.Content)
	})

	// Start the server
	log.Println("Server is running on http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
