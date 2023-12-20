package main

import (
	"log"

	"github.com/Furkan-Gulsen/golang-url-shortener/internal/adapters/handlers"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	log.Print("Starting Lambda")
	lambda.Start(handlers.SlackHandler)
}
