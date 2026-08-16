package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// These config are those which we use in tableplus , to setup our database connectio.

// Config → struct that holds all the configuration values for the application (e.g., database credentials, server port)
// database connection credentials
type Config struct {
	PublicHost             string
	Port                   string
	DBUser                 string
	DBPassword             string
	DBAdress               string
	DBName                 string
	JWTExpirationInSeconds int64
	JWTSecret              string
}

// Envs → global variable that holds the initialized configuration values, accessible throughout the application
var Envs = initConfig()

// initConfig → function that initializes the Config struct with values from environment variables or default values if not set
// database connection credentials
func initConfig() Config {
	// Load environment variables from .env file (if it exists) and handle any errors that occur during loading
	godotenv.Load()
	return Config{
		PublicHost:             getEnv("PublicHost", "http://localhost"),
		Port:                   getEnv("Port", "8080"),
		DBUser:                 getEnv("DBUser", "root"),
		DBPassword:             getEnv("DBPassword", ""),
		DBAdress:               fmt.Sprintf("%s:%s", getEnv("DBHost", "127.0.0.1"), getEnv("DBPort", "3306")),
		DBName:                 getEnv("DBName", "myapp"),
		JWTExpirationInSeconds: getEnvInt64("JWT_EXP", 3600*24*7),
		JWTSecret:              getEnv("JWT_SECRET", "not-secret-secret-anymore?"),
	}

}

// getEnv → helper function that retrieves the value of an environment variable or returns a fallback value if the variable is not set
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt64(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}
