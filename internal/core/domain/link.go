package domain

import "time"

type Link struct {
	Id          string    `dynamodbav:"id" json:"id"`
	OriginalURL string    `dynamodbav:"original_url" json:"original_url"`
	CreatedAt   time.Time `dynamodbav:"created_at" json:"created_at"`
	Stats       []Stats   `dynamodbav:"-" json:"stats"`
}
