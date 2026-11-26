package scraper

import (
	"context"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/kduvur100/baltimore-sports-scraper/internal/models"
)

type rssSource struct {
	URL  string
	Team models.Team
	Name string
}

var rssSources = []rssSource{
	{URL: "https://www.mlb.com/orioles/rss/news", Team: models.TeamOrioles, Name: "MLB Orioles"},
	{URL: "https://www.espn.com/espn/rss/mlb/news", Team: models.TeamOrioles, Name: "ESPN MLB"},
	{URL: "https://www.espn.com/espn/rss/nfl/news", Team: models.TeamRavens, Name: "ESPN NFL"},
	{URL: "https://feeds.baltimoresun.com/baltimoresun/sports/ravens", Team: models.TeamRavens, Name: "Baltimore Sun Ravens"},
	{URL: "https://feeds.baltimoresun.com/baltimoresun/sports/orioles", Team: models.TeamOrioles, Name: "Baltimore Sun Orioles"},
	{URL: "https://rss.app/feeds/Baltimore-Sports.xml", Team: "", Name: "Baltimore Sports"},
}

var oriolesKW = []string{"orioles", "o's", "baltimore baseball", "camden yards"}
var ravensKW = []string{"ravens", "lamar", "harbaugh", "m&t bank", "nfl ravens"}

type RSSScraper struct{ parser *gofeed.Parser }

func NewRSSScraper() *RSSScraper { return &RSSScraper{parser: gofeed.NewParser()} }

func (s *RSSScraper) Name() string            { return "RSS" }
func (s *RSSScraper) Interval() time.Duration { return 10 * time.Minute }

func (s *RSSScraper) Fetch(ctx context.Context) ([]models.Article, error) {
	type result struct {
		items []models.Article
		err   error
	}
	ch := make(chan result, len(rssSources))

	for _, src := range rssSources {
		src := src
		go func() {
			items, err := s.fetchFeed(ctx, src)
			ch <- result{items, err}
		}()
	}

	var all []models.Article
	for range rssSources {
		r := <-ch
		if r.err == nil {
			all = append(all, r.items...)
		}
	}
	return all, nil
}

func (s *RSSScraper) fetchFeed(ctx context.Context, src rssSource) ([]models.Article, error) {
	feed, err := s.parser.ParseURLWithContext(src.URL, ctx)
	if err != nil {
		return nil, err
	}

	var articles []models.Article
	for _, item := range feed.Items {
		pub := time.Now()
		if item.PublishedParsed != nil {
			pub = *item.PublishedParsed
		}

		team := src.Team
		if team == "" {
			team = detectTeam(item.Title + " " + item.Description)
		}
		if team == "" {
			continue
		}

		imgURL := ""
		if item.Image != nil {
			imgURL = item.Image.URL
		}

		articles = append(articles, models.Article{
			Title:       item.Title,
			URL:         item.Link,
			Summary:     stripHTML(item.Description),
			Source:      models.SourceRSS,
			Team:        team,
			Author:      strings.Join(personNames(item.Authors), ", "),
			ImageURL:    imgURL,
			PublishedAt: pub,
		})
	}
	return articles, nil
}

func detectTeam(text string) models.Team {
	lower := strings.ToLower(text)
	o := containsAny(lower, oriolesKW)
	r := containsAny(lower, ravensKW)
	switch {
	case o && r:
		return models.TeamBoth
	case o:
		return models.TeamOrioles
	case r:
		return models.TeamRavens
	}
	return ""
}

func containsAny(s string, kws []string) bool {
	for _, kw := range kws {
		if strings.Contains(s, kw) {
			return true
		}
	}
	return false
}

func personNames(persons []*gofeed.Person) []string {
	names := make([]string, 0, len(persons))
	for _, p := range persons {
		if p != nil && p.Name != "" {
			names = append(names, p.Name)
		}
	}
	return names
}

func stripHTML(s string) string {
	var b strings.Builder
	inTag := false
	for _, c := range s {
		switch {
		case c == '<':
			inTag = true
		case c == '>':
			inTag = false
		case !inTag:
			b.WriteRune(c)
		}
	}
	return strings.TrimSpace(b.String())
}
