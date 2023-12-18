package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	dynamoTableName string // DynamoDB table name
	redisAddress    string // Redis address
	redisPassword   string // Redis password
	redisDB         int    // Redis DB
	slackToken      string // Slack token
	slackChannelID  string // Slack channel ID
}

func NewConfig() *AppConfig {
	return &AppConfig{
		dynamoTableName: "UrlShortenerTable", // default value
		redisAddress:    "localhost:6379",    // default value
		redisPassword:   "",                  // default value
		redisDB:         0,                   // default value
		slackToken:      "",                  // default value
		slackChannelID:  "",                  // default value
	}
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
}

func (c *AppConfig) GetSlackParams() (string, string) {
	slackToken, tokenOK := os.LookupEnv("SLACK_TOKEN")
	slackChannelID, channelOK := os.LookupEnv("SLACK_CHANNEL_ID")
	if !tokenOK || !channelOK {
		return os.Getenv("SLACK_TOKEN"), os.Getenv("SLACK_CHANNEL_ID")
	}
	return slackToken, slackChannelID
}

func (c *AppConfig) GetTableName() string {
	tableName, ok := os.LookupEnv("TABLE")
	if !ok {
		fmt.Println("Need TABLE environment variable")
		return os.Getenv("TABLE")
	}
	return tableName
}

func (c *AppConfig) GetRedisParams() (string, string, int) {
	address, ok := os.LookupEnv("REDIS_ADDRESS")
	if !ok {
		fmt.Println("Need REDIS_ADDRESS environment variable")
		return c.redisAddress, c.redisPassword, c.redisDB
	}

	password, ok := os.LookupEnv("REDIS_PASSWORD")
	if !ok {
		fmt.Println("Need REDIS_PASSWORD environment variable")
		return address, c.redisPassword, c.redisDB
	}

	dbStr, ok := os.LookupEnv("REDIS_DB")
	if !ok {
		fmt.Println("Need REDIS_DB environment variable")
		return address, password, c.redisDB
	}

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		fmt.Printf("REDIS_DB environment variable is not a valid integer: %v\n", err)
		return address, password, c.redisDB
	}

	return address, password, db
}
