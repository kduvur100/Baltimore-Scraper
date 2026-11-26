package scraper

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"

	"github.com/kduvur100/baltimore-sports-scraper/internal/models"
	redisClient "github.com/kduvur100/baltimore-sports-scraper/internal/redis"
)

type Scraper interface {
	Name() string
	Fetch(ctx context.Context) ([]models.Article, error)
	Interval() time.Duration
}

type Manager struct {
	pool     *pgxpool.Pool
	rdb      *goredis.Client
	scrapers []Scraper
}

func NewManager(pool *pgxpool.Pool, rdb *goredis.Client) *Manager {
	return &Manager{
		pool: pool,
		rdb:  rdb,
		scrapers: []Scraper{
			NewRSSScraper(),
			NewGDELTScraper(),
		},
	}
}

func (m *Manager) Start(ctx context.Context) {
	log.Printf("Scraper manager starting %d scrapers", len(m.scrapers))
	for _, s := range m.scrapers {
		go m.run(ctx, s)
	}
}

func (m *Manager) run(ctx context.Context, s Scraper) {
	log.Printf("[%s] started (interval: %s)", s.Name(), s.Interval())
	m.execute(ctx, s)

	t := time.NewTicker(s.Interval())
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Printf("[%s] stopped", s.Name())
			return
		case <-t.C:
			m.execute(ctx, s)
		}
	}
}

func (m *Manager) execute(ctx context.Context, s Scraper) {
	articles, err := s.Fetch(ctx)
	if err != nil {
		log.Printf("[%s] fetch error: %v", s.Name(), err)
		return
	}
	newCount := 0
	for _, a := range articles {
		inserted, err := m.upsert(ctx, a)
		if err != nil {
			log.Printf("[%s] upsert error: %v", s.Name(), err)
			continue
		}
		if inserted {
			newCount++
			if err := redisClient.Publish(ctx, m.rdb, a); err != nil {
				log.Printf("[%s] redis publish: %v", s.Name(), err)
			}
		}
	}
	if newCount > 0 {
		log.Printf("[%s] +%d new articles", s.Name(), newCount)
	}
}

func (m *Manager) upsert(ctx context.Context, a models.Article) (bool, error) {
	tag, err := m.pool.Exec(ctx, `
		INSERT INTO articles (title, url, summary, source, team, author, image_url, sentiment_score, published_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT (url) DO NOTHING
	`, a.Title, a.URL, a.Summary, string(a.Source), string(a.Team),
		a.Author, a.ImageURL, a.SentimentScore, a.PublishedAt)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}
