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
	linkTableName := appConfig.GetLinkTableName()
	statsTableName := appConfig.GetStatsTableName()

	cache := cache.NewRedisCache(redisAddress, redisPassword, redisDB)

	linkRepo := repository.NewLinkRepository(context.TODO(), linkTableName)
	statsRepo := repository.NewStatsRepository(context.TODO(), statsTableName)

	linkService := services.NewLinkService(linkRepo, cache)
	statsService := services.NewStatsService(statsRepo, cache)

	handler := handlers.NewStatsFunctionHandler(linkService, statsService)

	lambda.Start(handler.Stats)
}
