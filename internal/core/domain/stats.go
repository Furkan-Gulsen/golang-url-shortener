package domain

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
	Id         string   `dynamodbav:"id" json:"id"`
	ClickCount int      `dynamodbav:"click_count" json:"click_count"`
	Platform   Platform `dynamodbav:"platform" json:"platform"`
	LinkID     string   `dynamodbav:"link_id" json:"link_id"`
}
