package db

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Book represents a book in the library.
// type Book struct {
// 	ID             uint   `json:"id" gorm:"primaryKey"`
// 	Title          string `json:"title"`
// 	Author         string `json:"author"`
// 	Category       string `json:"category"`
// 	Content        []byte `json:"-"` // PDF content stored as a BLOB
// 	CoverImage     []byte `json:"-"` // Cover image stored as a BLOB
// 	CoverImageType string `json:"-"` // MIME type for the cover image
// }

// ApplyMigrations applies all pending database migrations
func ApplyMigrations(migrationsPath, dbURL string) error {
	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		log.Fatalf("Failed to initialize migrations: %v", err)
		return err
	}

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
		return err
	}

	log.Println("Migrations applied successfully.")
	return nil
}
