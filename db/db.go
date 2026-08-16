package db

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
)

// NewMySqlStorage initializes a MySQL database connection using the provided configuration
// It configures basic connection pooling and verifies the connection with a Ping.
func NewMySqlStorage(cfg mysql.Config) (*sql.DB, error) {
	// Format the Data Source Name (DSN) from the provided MySQL configuration
	dsn := cfg.FormatDSN()

	// Open a new database connection using the MySQL driver and the formatted DSN
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// Set connection pool parameters to optimize database performance and resource usage
	db.SetConnMaxLifetime(5 * time.Minute)
	// Set the maximum number of open connections to the database to 25
	db.SetMaxOpenConns(25)
	// Set the maximum number of idle connections in the pool to 25
	db.SetMaxIdleConns(25)
	// Verify the database connection by sending a ping. If the ping fails, close the database connection and return the error
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
