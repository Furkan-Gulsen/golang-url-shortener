package handlers

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Furkan-Gulsen/golang-url-shortener/domain"
	"github.com/Furkan-Gulsen/golang-url-shortener/types"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
)

type ApiGatewayV2Handler struct {
	link *domain.Link
}

func NewAPIGatewayV2Handler(l *domain.Link) *ApiGatewayV2Handler {
	return &ApiGatewayV2Handler{link: l}
}

// func (h *ApiGatewayV2Handler) HandleRequest(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
// 	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
// 	defer cancel()

// 	switch req.RequestContext.HTTP.Method {
// 	case "POST":
// 		return h.createShortLink(timeoutCtx, req)
// 	case "GET":
// 		return h.getOriginalLink(timeoutCtx, req)
// 	default:
// 		return clientError(http.StatusMethodNotAllowed)
// 	}
// }

type RequestBody struct {
	Long string `json:"long"`
}

func (h *ApiGatewayV2Handler) CreateShortLink(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	var requestBody RequestBody
	err := json.Unmarshal([]byte(req.Body), &requestBody)
	if err != nil {
		return clientError(http.StatusBadRequest)
	}

	link := types.Link{
		Id:        uuid.New().String(),
		Long:      requestBody.Long,
		Short:     generateShortURL(requestBody.Long),
		CreatedAt: time.Now(),
	}

	err = h.link.Create(ctx, link)
	if err != nil {
		return serverError(err)
	}

	js, err := json.Marshal(link)
	if err != nil {
		return serverError(err)
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

func (h *ApiGatewayV2Handler) GetOriginalLink(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	shortLink := req.PathParameters["short"]

	longLink, err := h.link.Get(ctx, shortLink)
	if err != nil {
		return clientError(http.StatusNotFound)
	}

	js, err := json.Marshal(longLink)
	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(js),
	}, nil
}

func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func serverError(err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       err.Error(),
	}, nil
}
