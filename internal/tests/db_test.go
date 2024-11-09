package tests

import (
	"database/sql"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite" // SQLite migration driver
	_ "github.com/golang-migrate/migrate/v4/source/file"     // File source for migrations
	_ "github.com/mattn/go-sqlite3"                          // SQLite driver for SQL operations
	"github.com/stretchr/testify/suite"
)

type DBTestSuite struct {
	suite.Suite
	db      *sql.DB
	migrate *migrate.Migrate
	dbFile  string
}

func (suite *DBTestSuite) SetupSuite() {
	// Set up database and migrations
	suite.dbFile = "test.db"
	dbURL := "sqlite://" + suite.dbFile
	migrationsPath := "file://../../sql/migrations"

	// Initialize migrations
	var err error
	suite.migrate, err = migrate.New(migrationsPath, dbURL)
	suite.Require().NoError(err, "Failed to initialize migrations")

	// Apply migrations
	err = suite.migrate.Up()
	if err != nil && err != migrate.ErrNoChange {
		suite.FailNow("Failed to apply migrations", err)
	}

	// Open the database
	suite.db, err = sql.Open("sqlite3", suite.dbFile)
	suite.Require().NoError(err, "Failed to open test database")
}

func (suite *DBTestSuite) TearDownSuite() {
	// Cleanup resources
	if suite.db != nil {
		suite.db.Close()
	}
	if suite.dbFile != "" {
		os.Remove(suite.dbFile)
	}
	if suite.migrate != nil {
		suite.migrate.Close()
	}
}

func (suite *DBTestSuite) TestMigrations() {
	// Validate schema
	suite.validateSchema()
}

func (suite *DBTestSuite) validateSchema() {
	// Check if the 'books' table exists
	row := suite.db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='books';`)
	var tableName string
	err := row.Scan(&tableName)
	suite.NoError(err, "Table 'books' does not exist")
	suite.Equal("books", tableName, "Unexpected table name")

	// Validate columns in 'books' table
	rows, err := suite.db.Query(`PRAGMA table_info(books);`)
	suite.Require().NoError(err, "Failed to fetch table info")
	defer rows.Close()

	columns := map[string]bool{}
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dflt sql.NullString
		err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk)
		suite.Require().NoError(err, "Failed to scan column info")
		columns[name] = true
	}

	// Ensure all required columns exist
	expectedColumns := []string{"id", "title", "author", "category", "content"}
	for _, col := range expectedColumns {
		suite.True(columns[col], "Missing column '%s'", col)
	}
}

// Run the test suite
func TestDBTestSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}
