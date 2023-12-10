package handlers

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/Furkan-Gulsen/golang-url-shortener/domain"
	"github.com/Furkan-Gulsen/golang-url-shortener/types"
	"github.com/aws/aws-lambda-go/events"
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
		return clientError(http.StatusBadRequest, "Invalid JSON")
	}

	if requestBody.Long == "" {
		return clientError(http.StatusBadRequest, "URL cannot be empty")
	}
	if len(requestBody.Long) < 15 {
		return clientError(http.StatusBadRequest, "URL must be at least 15 characters long")
	}
	if !isValidLink(requestBody.Long) {
		return clientError(http.StatusBadRequest, "Invalid URL format")
	}

	link := types.Link{
		Id:          generateShortURL(requestBody.Long),
		OriginalURL: requestBody.Long,
		CreatedAt:   time.Now(),
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
	pathSegments := strings.Split(req.RawPath, "/")
	if len(pathSegments) < 2 {
		return clientError(http.StatusBadRequest, "Invalid URL path")
	}

	shortLink := pathSegments[len(pathSegments)-1]
	longLink, err := h.link.Get(ctx, shortLink)
	if err != nil {
		return clientError(http.StatusNotFound, "Link not found")
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

func clientError(status int, message string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       message,
	}, nil
}

func serverError(err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       err.Error(),
	}, nil
}

func isValidLink(u string) bool {
	re := regexp.MustCompile(`^(http|https)://`)
	if !re.MatchString(u) {
		return false
	}

	parsedURL, err := url.ParseRequestURI(u)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return false
	}

	return true
}
