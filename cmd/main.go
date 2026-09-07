package main

import (
	"database/sql"
	"log"
	"myproject/cmd/api"
	"myproject/config"
	"os"

	"myproject/db"

	"github.com/go-sql-driver/mysql"
)

// main → entry point of the application, initializes the API server and starts it
func main() {
	// Debug: log the configuration being used
	log.Printf("DB Config - Host: %s, User: %s, Name: %s", config.Envs.DBAdress, config.Envs.DBUser, config.Envs.DBName)
	// Debug: log the environment variables being used
	log.Printf("DEBUG env: DOCKER=%q DBHost=%q DBPort=%q DBUser=%q DBName=%q", os.Getenv("DOCKER"), os.Getenv("DBHost"), os.Getenv("DBPort"), os.Getenv("DBUser"), os.Getenv("DBName"))
	// Initialize the database connection using the configuration values
	db, err := db.NewMySqlStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAdress,
		Net:                  "tcp",
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

	// Initialize the API server with the configured address (from env) and database connection
	server := api.NewAPIServer(":"+config.Envs.Port, db)
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
