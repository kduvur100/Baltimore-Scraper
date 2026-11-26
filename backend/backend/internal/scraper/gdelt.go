package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/kduvur100/baltimore-sports-scraper/internal/models"
)

const gdeltBase = "https://api.gdeltproject.org/api/v2/doc/doc"

type gdeltResponse struct {
	Articles []gdeltArticle `json:"articles"`
}

type gdeltArticle struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Seendate    string `json:"seendate"`
	Domain      string `json:"domain"`
	Language    string `json:"language"`
	SocialImage string `json:"socialimage"`
}

type GDELTScraper struct{ client *http.Client }

func NewGDELTScraper() *GDELTScraper {
	return &GDELTScraper{client: &http.Client{Timeout: 30 * time.Second}}
}

func (g *GDELTScraper) Name() string            { return "GDELT" }
func (g *GDELTScraper) Interval() time.Duration { return 15 * time.Minute }

func (g *GDELTScraper) Fetch(ctx context.Context) ([]models.Article, error) {
	queries := []struct {
		q    string
		team models.Team
	}{
		{`"Baltimore Orioles" OR "Orioles baseball"`, models.TeamOrioles},
		{`"Baltimore Ravens" OR "Ravens NFL"`, models.TeamRavens},
	}

	var all []models.Article
	for _, q := range queries {
		items, err := g.fetchQuery(ctx, q.q, q.team)
		if err != nil {
			continue
		}
		all = append(all, items...)
	}
	return all, nil
}

func (g *GDELTScraper) fetchQuery(ctx context.Context, query string, team models.Team) ([]models.Article, error) {
	p := url.Values{}
	p.Set("query", query)
	p.Set("mode", "artlist")
	p.Set("maxrecords", "25")
	p.Set("format", "json")
	p.Set("timespan", "1d")
	p.Set("sort", "datedesc")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, gdeltBase+"?"+p.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "BaltimoreSportsHub/1.0")

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GDELT returned HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return nil, err
	}

	var result gdeltResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse GDELT response: %w", err)
	}

	var articles []models.Article
	for _, a := range result.Articles {
		if a.Language != "English" || a.Title == "" || a.URL == "" {
			continue
		}
		pub, err := time.Parse("20060102T150405Z", a.Seendate)
		if err != nil {
			pub = time.Now()
		}
		articles = append(articles, models.Article{
			Title:       a.Title,
			URL:         a.URL,
			Summary:     "via " + a.Domain,
			Source:      models.SourceGDELT,
			Team:        team,
			ImageURL:    a.SocialImage,
			PublishedAt: pub,
		})
	}
	return articles, nil
}
