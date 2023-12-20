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
	cache := cache.NewRedisCache(redisAddress, redisPassword, redisDB)
	linkTableName := appConfig.GetLinkTableName()
	statsTableName := appConfig.GetStatsTableName()

	linkRepo := repository.NewLinkRepository(context.TODO(), linkTableName)
	linkService := services.NewLinkService(linkRepo, cache)

	statsRepo := repository.NewStatsRepository(context.TODO(), statsTableName)
	statsService := services.NewStatsService(statsRepo, cache)

	handler := handlers.NewRedirectFunctionHandler(linkService, statsService)

	lambda.Start(handler.Redirect)
}
