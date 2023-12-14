package main

import (
	"context"

	"github.com/Furkan-Gulsen/golang-url-shortener/internal/adapters/cache"
	"github.com/Furkan-Gulsen/golang-url-shortener/internal/adapters/handlers"
	"github.com/Furkan-Gulsen/golang-url-shortener/internal/adapters/repository"
	"github.com/Furkan-Gulsen/golang-url-shortener/internal/config"
	"github.com/Furkan-Gulsen/golang-url-shortener/internal/core/services"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	appConfig := config.NewConfig()
	redisAddress, redisPassword, redisDB := appConfig.GetRedisParams()
	tableName := appConfig.GetTableName()

	db := repository.NewDynamoDBStore(context.TODO(), tableName)
	cache := cache.NewRedisCache(redisAddress, redisPassword, redisDB)

	domain := services.NewLinkDomain(db, cache)
	handler := handlers.NewAPIGatewayV2Handler(domain)
	lambda.Start(handler.CreateShortLink)
}
