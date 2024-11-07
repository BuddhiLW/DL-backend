package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Book represents a book in the library.
type Book struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Category string `json:"category"`
	Content  []byte `json:"-"` // PDF content stored as a BLOB
}

// InitDatabase initializes the SQLite database and applies migrations.
func InitDatabase() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("library.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	// Apply migrations
	err = db.AutoMigrate(&Book{})
	if err != nil {
		panic("failed to migrate database")
	}

	return db
}
