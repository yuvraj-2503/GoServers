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

var Configuration *ServerConfig

type ServerConfig struct {
	ServerPort  int
	SecretKey   string
	Msg91Config *Msg91Config
	MongoConfig *mongodb.MongoConfig
}

type Msg91Config struct {
	BaseUrl    string
	AuthKey    string
	TemplateId string
}

func init() {
	loadConfig()
}

func loadConfig() {
	_, currentFilePath, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFilePath)
	envFilePath := filepath.Join(currentDir, ".env")
	err := godotenv.Load(envFilePath)

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Panic(err)
	}

	Configuration = &ServerConfig{
		ServerPort:  port,
		SecretKey:   os.Getenv("SECRET_KEY"),
		Msg91Config: getMsg91Config(),
		MongoConfig: getMongoConfig(),
	}
}

func getMongoConfig() *mongodb.MongoConfig {
	return &mongodb.MongoConfig{
		ConnectionString: os.Getenv("MONGO_CONNECTION_STRING"),
		Database:         os.Getenv("DATABASE"),
		Username:         os.Getenv("DB_USER"),
		Password:         os.Getenv("DB_PASSWORD"),
	}
}

func getMsg91Config() *Msg91Config {
	return &Msg91Config{
		BaseUrl:    os.Getenv("MSG91_BASE_URL"),
		AuthKey:    os.Getenv("MSG91_AUTH_KEY"),
		TemplateId: os.Getenv("MSG91_TEMPLATE_ID"),
	}
}
