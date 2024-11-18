// pkg/utils/db.go
package utils

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// InitDB initializes and returns a PostgreSQL database connection.
func InitDB(connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	
	

	return db
}
