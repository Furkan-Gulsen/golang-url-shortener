package unit

import (
	"context"
	"testing"

	"github.com/Furkan-Gulsen/golang-url-shortener/internal/adapters/handlers"
	"github.com/stretchr/testify/assert"
)

func TestSlack(t *testing.T) {

	t.Run("Send Message to Slack", func(t *testing.T) {
		err := handlers.PostMessageToSlack(context.Background(), "Hello world! API Gateway message.")
		assert.Nil(t, err)
	})
}
