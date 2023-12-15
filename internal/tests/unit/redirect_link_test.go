package unit

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestRedirectLinkUnit(t *testing.T) {
	apiHandler := SetupTest()

	tests := []struct {
		shortLink        string
		expectStatusCode int
		expectLocation   string
		expectBody       string
	}{
		{
			shortLink:        "testid1",
			expectStatusCode: 301,
			expectLocation:   "https://example.com/link1",
			expectBody:       "",
		},
		{
			shortLink:        "testid2",
			expectStatusCode: 301,
			expectLocation:   "https://example.com/link2",
			expectBody:       "",
		},
		{
			shortLink:        "testid3",
			expectStatusCode: 301,
			expectLocation:   "https://example.com/link3",
			expectBody:       "",
		},
		{
			shortLink:        "nonexistentid",
			expectStatusCode: 404,
			expectLocation:   "",
			expectBody:       "Link not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.shortLink, func(t *testing.T) {
			request := events.APIGatewayV2HTTPRequest{
				RawPath: "/" + tt.shortLink,
			}

			response, err := apiHandler.Redirect(context.Background(), request)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectStatusCode, response.StatusCode)

			location := response.Headers["Location"]
			assert.Equal(t, tt.expectLocation, location)

			if tt.expectStatusCode == 404 {
				assert.Equal(t, tt.expectBody, response.Body)
			}
		})
	}
}
