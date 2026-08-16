package main

import (
	"database/sql"
	"log"
	"myproject/cmd/api"
	"myproject/config"

	"myproject/db"

	"github.com/go-sql-driver/mysql"
)

// main → entry point of the application, initializes the API server and starts it
func main() {
	// Initialize the database connection using the configuration values
	db, err := db.NewMySqlStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAdress,
		DBName:               config.Envs.DBName,
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	// Check for errors during database initialization and log them if any
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	// If the database connection is successful, check the connection by pinging the database and log the result
	initStorage(db)

	// If the database connection is successful, initialize the API server with the address and database connection, and start the server
	server := api.NewAPIServer(":8080", db)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

// initStorage → function that checks the database connection by pinging it, and logs the result (success or failure)
func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	log.Println("Database connection successful")
}
