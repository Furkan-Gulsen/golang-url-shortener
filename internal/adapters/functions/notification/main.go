package main

import (
	"github.com/Furkan-Gulsen/golang-url-shortener/internal/adapters/handlers"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handlers.PostMessageToSlack)
}
