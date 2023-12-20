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
		log.Print("Error loading .env file: ", err)
	}
}

func (c *AppConfig) GetSlackParams() (string, string) {
	slackToken, tokenOK := os.LookupEnv("SlackToken")
	slackChannelID, channelOK := os.LookupEnv("SlackChannelID")
	if !tokenOK || !channelOK {
		return os.Getenv("SlackToken"), os.Getenv("SlackChannelID")
	}
	return slackToken, slackChannelID
}

func (c *AppConfig) GetLinkTableName() string {
	tableName, ok := os.LookupEnv("LinkTableName")
	if !ok {
		fmt.Println("Need LinkTableName environment variable")
		return os.Getenv("LinkTableName")
	}
	return tableName
}

func (c *AppConfig) GetStatsTableName() string {
	tableName, ok := os.LookupEnv("StastTableName")
	if !ok {
		fmt.Println("Need STATS_TABLE environment variable")
		return os.Getenv("StastTableName")
	}
	return tableName
}

func (c *AppConfig) GetRedisParams() (string, string, int) {
	address, ok := os.LookupEnv("RedisAddress")
	if !ok {
		fmt.Println("Need RedisAddress environment variable")
		return c.redisAddress, c.redisPassword, c.redisDB
	}

	password, ok := os.LookupEnv("RedisPassword")
	if !ok {
		fmt.Println("Need RedisPassword environment variable")
		return address, c.redisPassword, c.redisDB
	}

	dbStr, ok := os.LookupEnv("RedisDB")
	if !ok {
		fmt.Println("Need RedisDB environment variable")
		return address, password, c.redisDB
	}

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		fmt.Printf("RedisDB environment variable is not a valid integer: %v\n", err)
		return address, password, c.redisDB
	}

	return address, password, db
}
