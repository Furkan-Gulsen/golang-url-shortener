package mock

import "github.com/Furkan-Gulsen/golang-url-shortener/internal/core/domain"

var MockLinkData []domain.Link = []domain.Link{
	{Id: "testid1", OriginalURL: "https://example.com/link1"},
	{Id: "testid2", OriginalURL: "https://example.com/link2"},
	{Id: "testid3", OriginalURL: "https://example.com/link3"},
}

var MockStatsData []domain.Stats = []domain.Stats{
	{Id: "abcdefg1", Platform: domain.PlatformUnknown, LinkID: "testid1"},
	{Id: "abcdefg2", Platform: domain.PlatformInstagram, LinkID: "testid2"},
	{Id: "abcdefg3", Platform: domain.PlatformTwitter, LinkID: "testid3"},
}
