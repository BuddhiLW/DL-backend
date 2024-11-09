package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/BuddhiLW/DL-backend/internal/db"
	"github.com/BuddhiLW/DL-backend/internal/handlers"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// BookDB wraps the sql.DB connection and sqlc-generated Queries for database interactions.
type BookDB struct {
	dbConn *sql.DB
	*db.Queries
}

// NewBookDB initializes a new BookDB instance.
func NewBookDB(dbConn *sql.DB) *BookDB {
	return &BookDB{
		dbConn:  dbConn,
		Queries: db.New(dbConn),
	}
}

// callTx executes a transaction with the provided function, rolling back if an error occurs.
func (b *BookDB) callTx(ctx context.Context, fn func(*db.Queries) error) error {
	tx, err := b.dbConn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	txQueries := db.New(tx)
	defer tx.Rollback()
	if err := fn(txQueries); err != nil {
		return err
	}
	return tx.Commit()
}

func main() {
	// Set up the context
	// ctx := context.Background()

	// Connect to the database
	dbURL := "postgres://user:password@localhost:5151/books?sslmode=disable"
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Apply database migrations
	migrationsPath := "file://sql/migrations"
	log.Println("Applying database migrations...")
	if err := db.ApplyMigrations(migrationsPath, dbURL); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	// Initialize BookDB with the database connection
	bookDB := NewBookDB(dbConn)

	// Set up the Gin router
	router := handlers.SetupRouter(bookDB.Queries)

	// Start the server
	log.Println("Server running on http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
