package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/BuddhiLW/DL-backend/internal/db"
	"github.com/gin-gonic/gin"
)

// UploadBook handles POST /books to upload a new book.
func UploadBook(queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		// Extract form data
		title := c.PostForm("title")
		author := c.PostForm("author")
		category := c.PostForm("category")

		// Validate required fields
		if title == "" || author == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Title and Author are required"})
			return
		}

		// Read the PDF content
		pdfBuffer, err := readFileFromRequest(c, "file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Read the cover image
		coverBuffer, coverType, err := readFileWithMimeFromRequest(c, "cover")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create the book record using sqlc
		params := db.CreateBookParams{
			Title:          title,
			Author:         author,
			Category:       sqlNullString(category),
			Content:        pdfBuffer.Bytes(),
			CoverImage:     coverBuffer.Bytes(),
			CoverImageType: sqlNullString(coverType),
		}

		if err := queries.CreateBook(ctx, params); err != nil {
			log.Printf("Error creating book: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save book"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Book uploaded successfully"})
	}
}

// DownloadBook handles GET /books/:id/download to retrieve a book's content.
func DownloadBook(queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		id, err := parseIDParam(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Fetch the book
		book, err := queries.GetBook(ctx, int32(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}

		c.Header("Content-Disposition", "attachment; filename=\""+book.Title+".pdf\"")
		c.Data(http.StatusOK, "application/pdf", book.Content)
	}
}

// GetCoverImage handles GET /books/:id/cover to fetch the cover image.
func GetCoverImage(queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		id, err := parseIDParam(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Fetch the book
		book, err := queries.GetBook(ctx, int32(id))
		if err != nil {
			log.Printf("Error fetching book: %v", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}

		c.Header("Content-Type", book.CoverImageType.String)
		c.Data(http.StatusOK, book.CoverImageType.String, book.CoverImage)
	}
}

// parseIDParam parses the :id parameter from the URL
func parseIDParam(c *gin.Context) (int64, error) {
	id := c.Param("id")
	if id == "" {
		return 0, errors.New("missing id parameter")
	}
	var parsedID int64
	if _, err := fmt.Sscanf(id, "%d", &parsedID); err != nil {
		return 0, errors.New("invalid id parameter")
	}
	return parsedID, nil
}

// readFileFromRequest reads a file from the request and returns its content as a buffer
func readFileFromRequest(c *gin.Context, field string) (*bytes.Buffer, error) {
	file, err := c.FormFile(field)
	if err != nil {
		return nil, errors.New("file upload failed: " + field)
	}

	fileContent, err := file.Open()
	if err != nil {
		return nil, errors.New("failed to open uploaded file")
	}
	defer fileContent.Close()

	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, fileContent); err != nil {
		return nil, errors.New("failed to read file content")
	}

	return &buffer, nil
}

// readFileWithMimeFromRequest reads a file from the request and returns its content and MIME type
func readFileWithMimeFromRequest(c *gin.Context, field string) (*bytes.Buffer, string, error) {
	file, err := c.FormFile(field)
	if err != nil {
		return nil, "", errors.New("file upload failed: " + field)
	}

	fileContent, err := file.Open()
	if err != nil {
		return nil, "", errors.New("failed to open uploaded file")
	}
	defer fileContent.Close()

	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, fileContent); err != nil {
		return nil, "", errors.New("failed to read file content")
	}

	mimeType := file.Header.Get("Content-Type")
	return &buffer, mimeType, nil
}

// sqlNullString converts a string to sql.NullString
func sqlNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}
