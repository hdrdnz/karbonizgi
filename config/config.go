package config

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/sashabaranov/go-openai"
)

type Database struct {
	Host     string `json:"host"`
	DbPort   string `json:"db_port"`
	User     string `json:"user"`
	Password string `json:"db_password"`
	Name     string `json:"name"`
}
type Docker struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
	Repo     string `json:"repo"`
}
type Container struct {
	ImageName     string `json:"imageName"`
	ContainerName string `json:"containerName"`
	Port          string `json:"port"`
}
type Server struct {
	User string `json:"remoteUser"`
	Host string `json:"remoteHost"`
}
type Custom struct {
	Header string `json:"header"`
}
type OpenAI struct {
	ModelName string `json:"modelName"`
	Key       string `json:"key"`
}
type Config struct {
	Database  Database  `json:"database"`
	Docker    Docker    `json:"docker"`
	Container Container `json:"container"`
	Server    Server    `json:"server"`
	Custom    Custom    `json:"custom"`
	OpenAI    OpenAI    `json:"openai"`
}

var config *Config
var redisClient *redis.Client
var ctx = context.Background()
var client *openai.Client

func LoadConfig(filename string) (*Config, error) {
	if config != nil {
		return config, nil
	}
	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()
	config = &Config{}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}

func GetConfig() *Config {
	return config
}
func GetRedis() (*redis.Client, context.Context) {
	return redisClient, ctx
}
func LoadRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.Server.Host + ":6379",
		Password: "",
		DB:       0,
	})
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis bağlantı hatası: %v", err)
	}
}

func LoadClient() {
	client = openai.NewClient(config.OpenAI.Key)
}

func GetModulName() string {
	return config.OpenAI.ModelName
}

func GetClient() *openai.Client {
	return client
}

func GetEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
