package news

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/mmcdole/gofeed"
	"go.mongodb.org/mongo-driver/mongo"

	"go-news-feed/pkg/model"
)

// Service - interface
//
//go:generate mockgen -source=service.go -destination=service_mock.go --package=news
type Service interface {
	Find(ctx context.Context, sr model.FindRequest) (model.FindResponse, error)
	Load(ctx context.Context, feedURL string) ([]model.Article, error)
}

type service struct {
	feedParser *gofeed.Parser
	repository Repository
}

// newService - constructor
func newService(repository Repository) Service {
	return &service{
		feedParser: gofeed.NewParser(),
		repository: repository,
	}
}

func (s *service) Find(ctx context.Context, sr model.FindRequest) (model.FindResponse, error) {
	return s.repository.Find(ctx, sr)
}

func (s *service) Load(ctx context.Context, feedURL string) ([]model.Article, error) {
	articles, err := s.loadArticlesFromFeed(ctx, feedURL)
	if err != nil {
		return nil, err
	}

	if err := s.saveArticles(ctx, articles); err != nil {
		return nil, err
	}

	return articles, nil
}

// loadArticlesFromFeed and convert to an article slice ordered by published time (asc)
func (s *service) loadArticlesFromFeed(ctx context.Context, feedURL string) ([]model.Article, error) {
	var articles = make(model.Articles, 0)
	sources := s.getSources(feedURL)

	for _, source := range sources {
		feed, err := s.feedParser.ParseURLWithContext(source.FeedURL, ctx)
		if err != nil {
			return nil, err
		}

		result, err := s.parseFeed(feed)
		if err != nil {
			return nil, err
		}

		articles = append(articles, result...)
	}

	// could use sort from gofeed.Feed model
	// but adding in the article
	// just for the sake of an example
	// of how to sort a custom slice
	sort.Sort(articles)

	return articles, nil
}

// saveArticles persists new articles
func (s *service) saveArticles(ctx context.Context, articles []model.Article) error {
	for _, article := range articles {
		if err := s.validateArticle(ctx, article); err != nil {
			continue
		}

		if err := s.repository.Create(ctx, article); err != nil {
			return err
		}
	}

	return nil
}

// validateArticle checks if article exists already
func (s *service) validateArticle(ctx context.Context, article model.Article) error {
	tempArticle, err := s.repository.FindByID(ctx, article.ID)
	if err != nil {
		// if error is no documents found it means
		// document is valid to be created
		if err == mongo.ErrNoDocuments {
			return nil
		}

		return err
	}

	if tempArticle.ID != "" {
		return fmt.Errorf("article id %s exists already", tempArticle.ID)
	}

	return nil
}

// getSources from a feedURL
// it returns default sources if feedURL is not provided or not found
func (s *service) getSources(feedURL string) []model.Source {
	var sources = make([]model.Source, 0)

	if _, ok := model.DefaultSources[feedURL]; !ok {
		for _, source := range model.DefaultSources {
			sources = append(sources, source)
		}

		return sources
	}

	sources = append(sources, model.DefaultSources[feedURL])

	return sources
}

// parseFeed and returns the slice of articles
func (s *service) parseFeed(feed *gofeed.Feed) ([]model.Article, error) {
	if feed == nil || feed.Items == nil {
		return nil, errors.New("no feed or articles found")
	}

	var articles = make(model.Articles, len(feed.Items))
	for i, item := range feed.Items {
		article := model.Article{
			ID:                item.GUID,
			Title:             item.Title,
			Descriptiopn:      item.Description,
			Link:              item.Link,
			Source:            s.getSourceByLink(item.Link),
			PublishedDateTime: item.PublishedParsed,
			UpdatedDateTime:   item.UpdatedParsed,
		}

		articles[i] = article
	}

	return articles, nil
}

// getSourceByLink-  workaround to find the source based on the article link as bbc doesn't have
// feedLink to easily being identified and the feed link provided is the same for News and Tech
func (s *service) getSourceByLink(link string) model.Source {
	// it means it's bbc source
	if strings.Contains(link, model.SourceBBC) {
		if strings.Contains(link, model.CategoryTechnology) {
			return model.DefaultSources[fmt.Sprintf("%s/%s", model.SourceBBC, model.CategoryTechnology)]
		}

		return model.DefaultSources[fmt.Sprintf("%s/%s", model.SourceBBC, model.CategoryUK)]
	}

	// doing the same for sky to be consitent, but didn't really need to
	if strings.Contains(link, model.SourceSky) {
		if strings.Contains(link, model.CategoryTechnology) {
			return model.DefaultSources[fmt.Sprintf("%s/%s", model.SourceSky, model.CategoryTechnology)]
		}

		return model.DefaultSources[fmt.Sprintf("%s/%s", model.SourceSky, model.CategoryUK)]
	}

	return model.Source{}
}
