package handlers

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/Furkan-Gulsen/golang-url-shortener/internal/core/domain"
	"github.com/Furkan-Gulsen/golang-url-shortener/internal/core/services"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type RequestBody struct {
	Long string `json:"long"`
}
type GenerateLinkFunctionHandler struct {
	linkService  *services.LinkService
	statsService *services.StatsService
}

func NewGenerateLinkFunctionHandler(l *services.LinkService, s *services.StatsService) *GenerateLinkFunctionHandler {
	return &GenerateLinkFunctionHandler{linkService: l, statsService: s}
}

func (h *GenerateLinkFunctionHandler) CreateShortLink(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	var requestBody RequestBody
	err := json.Unmarshal([]byte(req.Body), &requestBody)
	if err != nil {
		return ClientError(http.StatusBadRequest, "Invalid JSON")
	}

	if requestBody.Long == "" {
		return ClientError(http.StatusBadRequest, "URL cannot be empty")
	}
	if len(requestBody.Long) < 15 {
		return ClientError(http.StatusBadRequest, "URL must be at least 15 characters long")
	}
	if !IsValidLink(requestBody.Long) {
		return ClientError(http.StatusBadRequest, "Invalid URL format")
	}

	link := domain.Link{
		Id:          GenerateShortURLID(8),
		OriginalURL: requestBody.Long,
		CreatedAt:   time.Now(),
	}

	err = h.linkService.Create(ctx, link)
	if err != nil {
		return ServerError(err)
	}

	js, err := json.Marshal(link)
	if err != nil {
		return ServerError(err)
	}

	err = h.statsService.Create(ctx, domain.Stats{
		Id:        uuid.NewString(),
		LinkID:    link.Id,
		Platform:  domain.PlatformTwitter,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("failed to create stats: ", err)
	}

	sendMessageToQueue(ctx, link)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(js),
	}, nil
}

func sendMessageToQueue(ctx context.Context, link domain.Link) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err.Error())
		return
	}

	sqsClient := sqs.NewFromConfig(cfg)
	queueUrl := os.Getenv("QueueUrl")

	if queueUrl == "" {
		log.Println("QueueUrl is not set")
		return
	}

	_, err = sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &queueUrl,
		MessageBody: aws.String("The system generated a short URL with the ID " + link.Id),
	})

	if err != nil {
		fmt.Printf("Failed to send message to SQS, %v", err.Error())
	}
}

func GenerateShortURLID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		charIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[charIndex.Int64()]
	}
	return string(result)
}
