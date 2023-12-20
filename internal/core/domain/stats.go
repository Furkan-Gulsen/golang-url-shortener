package domain

import "time"

type Platform int

const (
	PlatformUnknown Platform = iota
	PlatformInstagram
	PlatformTwitter
	PlatformYouTube
)

func (p Platform) String() string {
	switch p {
	case PlatformInstagram:
		return "Instagram"
	case PlatformTwitter:
		return "Twitter"
	case PlatformYouTube:
		return "YouTube"
	default:
		return "Unknown"
	}
}

type Stats struct {
	Id        string    `dynamodbav:"id" json:"id"`
	Platform  Platform  `dynamodbav:"platform" json:"platform"`
	LinkID    string    `dynamodbav:"link_id" json:"link_id"`
	CreatedAt time.Time `dynamodbav:"created_at" json:"created_at"`
}
