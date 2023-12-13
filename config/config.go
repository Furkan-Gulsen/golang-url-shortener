package config

import (
	"fmt"
	"os"
	"strconv"
)

type AppConfig struct {
	dynamoTableName string // DynamoDB table name
	redisAddress    string // Redis address
	redisPassword   string // Redis password
	redisDB         int    // Redis DB
}

func NewConfig() *AppConfig {
	return &AppConfig{
		dynamoTableName: "UrlShortenerTable",
		redisAddress:    "localhost:6379",
		redisPassword:   "",
		redisDB:         0,
	}
}

func (c *AppConfig) GetTableName() string {
	tableName, ok := os.LookupEnv("TABLE")
	if !ok {
		fmt.Println("Need TABLE environment variable")
		return c.dynamoTableName
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
