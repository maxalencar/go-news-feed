package model

const (
	SourceBBC = "bbc.co.uk/news"
	SourceSky = "news.sky.com"
)

var Sources []Source

type Source struct {
	Category string `json:"category,omitempty" bson:"category,omitempty"`
	FeedURL  string `json:"feedUrl,omitempty" bson:"feedUrl,omitempty"`
	Provider string `json:"provider,omitempty" bson:"provider,omitempty"`
}

var DefaultSources = map[string]Source{
	"bbc.co.uk/news/uk": {
		Category: CategoryUK,
		FeedURL:  "https://feeds.bbci.co.uk/news/uk/rss.xml",
		Provider: ProviderBBC,
	},
	"bbc.co.uk/news/technology": {
		Category: CategoryTechnology,
		FeedURL:  "https://feeds.bbci.co.uk/news/technology/rss.xml",
		Provider: ProviderBBC,
	},
	"news.sky.com/uk": {
		Category: CategoryUK,
		FeedURL:  "https://feeds.skynews.com/feeds/rss/uk.xml",
		Provider: ProviderSky,
	},
	"news.sky.com/technology": {
		Category: CategoryTechnology,
		FeedURL:  "https://feeds.skynews.com/feeds/rss/technology.xml",
		Provider: ProviderSky,
	},
}
