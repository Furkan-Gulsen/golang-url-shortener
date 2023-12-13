package main

import (
	"context"

	"github.com/Furkan-Gulsen/golang-url-shortener/config"
	"github.com/Furkan-Gulsen/golang-url-shortener/domain"
	"github.com/Furkan-Gulsen/golang-url-shortener/handlers"
	"github.com/Furkan-Gulsen/golang-url-shortener/store"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	appConfig := config.NewConfig()
	redisAddress, redisPassword, redisDB := appConfig.GetRedisParams()
	tableName := appConfig.GetTableName()

	db := store.NewDynamoDBStore(context.TODO(), tableName)
	cache := store.NewRedisCache(redisAddress, redisPassword, redisDB)

	domain := domain.NewLinkDomain(db, cache)
	handler := handlers.NewAPIGatewayV2Handler(domain)
	lambda.Start(handler.CreateShortLink)
}
