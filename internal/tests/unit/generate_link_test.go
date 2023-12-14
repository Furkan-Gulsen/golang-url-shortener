package unit

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func GenerateLinkUnitTests(t *testing.T) {
	apiHandler := SetupTest()

	tests := []struct {
		longURL            string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			longURL:            "https://example.com/abcdefg",
			expectedStatusCode: 200,
			expectedBody:       "",
		},
		{
			longURL:            "",
			expectedStatusCode: 400,
			expectedBody:       "URL cannot be empty",
		},
		{
			longURL:            "invalid",
			expectedStatusCode: 400,
			expectedBody:       "URL must be at least 15 characters long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.longURL, func(t *testing.T) {
			body := `{"long": "` + tt.longURL + `"}`
			request := events.APIGatewayV2HTTPRequest{Body: body}
			response, err := apiHandler.CreateShortLink(context.Background(), request)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatusCode, response.StatusCode)

			if tt.expectedStatusCode != 200 {
				assert.Equal(t, tt.expectedBody, response.Body)
			}
		})
	}

}
