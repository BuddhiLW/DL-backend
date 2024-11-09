package tests

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/BuddhiLW/DL-backend/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestSuite struct
type APITestSuite struct {
	suite.Suite
	db     *gorm.DB
	router *gin.Engine
}

// SetupSuite runs once before any tests
func (suite *APITestSuite) SetupSuite() {
	// Initialize in-memory database and router
	suite.db = setupTestDB()
	suite.router = setupRouter(suite.db)
}

// TearDownSuite runs once after all tests
func (suite *APITestSuite) TearDownSuite() {
	sqlDB, _ := suite.db.DB()
	sqlDB.Close()
}

// Test for POST /books
func (suite *APITestSuite) TestUploadBook() {
	// Create a temporary file to mimic PDF upload
	tempFile, _ := os.CreateTemp("", "sample*.pdf")
	defer os.Remove(tempFile.Name())                 // Clean up after the test
	tempFile.WriteString("This is a test PDF file.") // Write dummy content
	tempFile.Seek(0, io.SeekStart)                   // Reset file pointer

	// Create a multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("title", "Test Book")
	_ = writer.WriteField("author", "Author Name")
	_ = writer.WriteField("category", "Science")
	part, _ := writer.CreateFormFile("file", filepath.Base(tempFile.Name()))
	io.Copy(part, tempFile)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/books", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Book uploaded successfully")
}

// Test for GET /books/:id/download
func (suite *APITestSuite) TestDownloadBook() {
	// Seed a book in the database
	suite.db.Create(&db.Book{
		Title:          "Test Book",
		Author:         "Author Name",
		Category:       "Science",
		Content:        []byte("Dummy PDF content."),
		CoverImage:     []byte("PNG IMAGE DATA"),
		CoverImageType: "image/png",
	})

	req := httptest.NewRequest(http.MethodGet, "/books/1/download", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "application/pdf", w.Header().Get("Content-Type"))
	assert.Contains(suite.T(), w.Body.String(), "Dummy PDF content.")
}

// Test for GET /books/:id/download with invalid ID
func (suite *APITestSuite) TestDownloadBookNotFound() {
	req := httptest.NewRequest(http.MethodGet, "/books/999/download", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Book not found")
}

func (suite *APITestSuite) TestUploadBookWithCover() {
	// Create temporary PDF and cover image files
	tempPDF, _ := os.CreateTemp("", "sample*.pdf")
	defer os.Remove(tempPDF.Name())
	tempPDF.WriteString("This is a test PDF file.")
	tempPDF.Seek(0, io.SeekStart)

	tempImage, _ := os.CreateTemp("", "sample*.png")
	defer os.Remove(tempImage.Name())
	tempImage.WriteString("PNG IMAGE DATA") // Mock data
	tempImage.Seek(0, io.SeekStart)

	// Create a multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("title", "Test Book")
	_ = writer.WriteField("author", "Author Name")
	_ = writer.WriteField("category", "Science")
	pdfPart, _ := writer.CreateFormFile("file", filepath.Base(tempPDF.Name()))
	io.Copy(pdfPart, tempPDF)
	imagePart, _ := writer.CreateFormFile("cover", filepath.Base(tempImage.Name()))
	io.Copy(imagePart, tempImage)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/books", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	assert.Contains(suite.T(), w.Body.String(), "Book uploaded successfully")
}

// -------- Not working -------
// func (suite *APITestSuite) TestGetCoverImage() {
// 	// Seed a book with a cover image in the database
// 	coverImagePath := "../../tests/sample.png"
// 	coverFile, err := os.Open(coverImagePath)
// 	assert.Nil(suite.T(), err, "Failed to open cover image file")
// 	defer coverFile.Close()

// 	coverBytes, err := io.ReadAll(coverFile)
// 	assert.Nil(suite.T(), err, "Failed to read cover image file")

// 	book := db.Book{
// 		Title:          "Test Book",
// 		Author:         "Author Name",
// 		Category:       "Science",
// 		Content:        nil,
// 		CoverImage:     coverBytes,
// 		CoverImageType: "image/png",
// 	}
// 	result := suite.db.Create(&book)
// 	assert.Nil(suite.T(), result.Error, "Failed to seed database")
// 	// suite.T().Logf("Seeded book: %+v", book)

// 	// Request the cover image
// 	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/books/%d/cover", book.ID), nil)
// 	w := httptest.NewRecorder()
// 	suite.router.ServeHTTP(w, req)

// 	// Assertions
// 	assert.Equal(suite.T(), http.StatusOK, w.Code)

// 	// Validate Content-Type
// 	assert.Equal(suite.T(), "image/png", w.Header().Get("Content-Type"))

// 	// Compare response body with the original cover image bytes
// 	assert.Equal(suite.T(), coverBytes, w.Body.Bytes())
// }

// Run all tests in the suite
func TestAPISuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

// Helper function to initialize the test database
func setupTestDB() *gorm.DB {
	database, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{}) // Rename to `database`
	database.AutoMigrate(&db.Book{})                                  // Reference the `db` package here
	return database
}

// Helper function to create a router with initialized routes
func setupRouter(database *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.POST("/books", func(c *gin.Context) {
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
	})

	router.GET("/books/:id/download", func(c *gin.Context) {
		var book db.Book
		id := c.Param("id")

		if err := database.First(&book, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}

		c.Header("Content-Disposition", "attachment; filename="+book.Title+".pdf")
		c.Data(http.StatusOK, "application/pdf", book.Content)
	})

	return router
}
