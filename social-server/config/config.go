package config

import (
	"github.com/joho/godotenv"
	"log"
	mongodb "mongo-utils"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

type ServerConfig struct {
	ServerPort  int
	MongoConfig *mongodb.MongoConfig
	SecretKey   string
}

var Configuration *ServerConfig

func init() {
	loadConfig()
}

func loadConfig() {
	_, currentFilePath, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFilePath)
	envFilePath := filepath.Join(currentDir, ".env")
	err := godotenv.Load(envFilePath)
	if err != nil {
		log.Fatalf("Error loading .env file, reason: %v", err)
	}

	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Fatalf("Error loading server port, reason: %v", err)
	}

	Configuration = &ServerConfig{
		ServerPort: port,
		MongoConfig: &mongodb.MongoConfig{
			ConnectionString: os.Getenv("MONGO_CONNECTION_STRING"),
			Database:         os.Getenv("DATABASE"),
			Username:         os.Getenv("DB_USER"),
			Password:         os.Getenv("DB_PASSWORD"),
		},
		SecretKey: os.Getenv("SECRET_KEY"),
	}
}
