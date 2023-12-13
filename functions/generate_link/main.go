package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Furkan-Gulsen/golang-url-shortener/domain"
	"github.com/Furkan-Gulsen/golang-url-shortener/handlers"
	"github.com/Furkan-Gulsen/golang-url-shortener/store"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	tableName, ok := os.LookupEnv("TABLE")
	if !ok {
		fmt.Println("Need TABLE environment variable")
		tableName = "UrlShortenerTable"
	}

	db := store.NewDynamoDBStore(context.TODO(), tableName)
	// cache := store.NewRedisCache(context.TODO())
	domain := domain.NewLinkDomain(db)
	handler := handlers.NewAPIGatewayV2Handler(domain)
	lambda.Start(handler.CreateShortLink)
}
