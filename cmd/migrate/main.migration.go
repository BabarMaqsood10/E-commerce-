package main

import (
	"log"
	"myproject/config"
	"myproject/db"
	"os"

	mysqlCfg "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// here migration occurs , the other main handles the server and other stuff, this is only for migration
// it actually changes the database schema, so it should be used with caution and only when necessary
// You don't need to restart your migration code every time your API receives a request. That's why the two programs are separate.
// // main → entry point of the migration tool, connects to the database and applies or rolls back database schema changes when explicitly requested.

func main() {
	// Initialize the database connection using the configuration values
	db, err := db.NewMySqlStorage(mysqlCfg.Config{
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
	if err := db.Ping(); err != nil {
		log.Fatal("failed to ping database:", err)
	}
	// Create a new migrate instance using the MySQL database connection and the file source for migrations, and check for errors during initialization
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatal("failed to create migrate driver:", err)
	}
	// Create a new migrate instance using the MySQL database connection and the file source for migrations, and check for errors during initialization
	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"mysql",
		driver,
	)
	// Check for errors during initialization of the migrate instance and log them if any
	if err != nil {
		log.Fatal("failed to create migrate instance:", err)
	}
	// Get the last command-line argument to determine whether to apply or rollback migrations, and execute the corresponding migration action based on the command, logging the results of the migration operation (success or failure) accordingly
	cmd := os.Args[(len(os.Args) - 1)]
	// Get the last command-line argument to determine whether to apply or rollback migrations, and execute the corresponding migration action based on the command, logging the results of the migration operation (success or failure) accordingly

	switch cmd {
	case "up":
		// Apply the migrations to the database, and check for errors during the migration process, logging the results accordingly (success or failure)
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("failed to apply migrations:", err)
		}
		log.Println("migrations applied successfully")
		// Rollback the last migration, and check for errors during the rollback process, logging the results accordingly (success or failure)
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("failed to rollback migrations:", err)
		}
		log.Println("migrations rolled back successfully")
		// If the command is not recognized (neither "up" nor "down"), log an error message and exit the program
	default:
		log.Fatal("invalid command, use 'up' or 'down'")
	}
}
