package handlers

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Furkan-Gulsen/golang-url-shortener/internal/core/domain"
	"github.com/aws/aws-lambda-go/events"
)

type RequestBody struct {
	Long string `json:"long"`
}

func (h *ApiGatewayV2Handler) CreateShortLink(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
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
		Id:          generateShortURL(requestBody.Long),
		OriginalURL: requestBody.Long,
		CreatedAt:   time.Now(),
	}

	err = h.link.Create(ctx, link)
	if err != nil {
		return ServerError(err)
	}

	js, err := json.Marshal(link)
	if err != nil {
		return ServerError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(js),
	}, nil
}

func generateShortURL(longURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(longURL))
	shortURL := hex.EncodeToString(hasher.Sum(nil))
	return shortURL[:8]
}
