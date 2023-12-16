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
	cache := cache.NewRedisCache(redisAddress, redisPassword, redisDB)

	linkRepo := repository.NewDynamoDBStore(context.TODO(), tableName)
	statisticsRepo := repository.NewStatisticsRepository(context.TODO(), tableName)

	linkService := services.NewLinkService(linkRepo, cache)
	statisticsService := services.NewStatisticsService(statisticsRepo, cache)

	handler := handlers.NewStatisticsFunctionHandler(linkService, statisticsService)

	lambda.Start(handler.Statistics)
}
