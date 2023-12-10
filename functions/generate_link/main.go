package main

import (
	"context"
	"os"

	"github.com/Furkan-Gulsen/golang-url-shortener/domain"
	"github.com/Furkan-Gulsen/golang-url-shortener/handlers"
	"github.com/Furkan-Gulsen/golang-url-shortener/store"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	tableName, ok := os.LookupEnv("TABLE")
	if !ok {
		panic("Need TABLE environment variable")
	}

	db := store.NewDynamoDBStore(context.TODO(), tableName)
	domain := domain.NewLinkDomain(db)
	handler := handlers.NewAPIGatewayV2Handler(domain)
	lambda.Start(handler.CreateShortLink)
}
