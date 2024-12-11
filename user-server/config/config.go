package config

import (
	"blob-manager/aws"
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
	ServerPort     int
	SecretKey      string
	Msg91Config    *Msg91Config
	SendgridConfig *SendgridConfig
	MongoConfig    *mongodb.MongoConfig
	AWSConfig      *aws.AWSConfig
}

type Msg91Config struct {
	BaseUrl    string
	AuthKey    string
	TemplateId string
}

type SendgridConfig struct {
	SenderId       string
	SendgridApiKey string
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
		ServerPort:     port,
		SecretKey:      os.Getenv("SECRET_KEY"),
		Msg91Config:    getMsg91Config(),
		SendgridConfig: getSendgridConfig(),
		MongoConfig:    getMongoConfig(),
		AWSConfig:      getAWSConfig(),
	}
}

func getAWSConfig() *aws.AWSConfig {
	return &aws.AWSConfig{
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AccessKeySecret: os.Getenv("AWS_ACCESS_KEY_SECRET"),
		Region:          os.Getenv("AWS_REGION"),
		BucketName:      os.Getenv("AWS_BUCKET_NAME"),
		UploadTimeout:   1000000,
		BaseURL:         os.Getenv("AWS_BASE_URL"),
	}
}

func getSendgridConfig() *SendgridConfig {
	return &SendgridConfig{
		SenderId:       os.Getenv("SENDGRID_SENDER_ID"),
		SendgridApiKey: os.Getenv("SENDGRID_API_KEY"),
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
