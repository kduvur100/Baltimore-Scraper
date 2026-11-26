package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(databaseURL string) (*pgxpool.Pool, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return pool, nil
}

func RunMigrations(pool *pgxpool.Pool) error {
	stmts := []string{
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`,
		`CREATE TABLE IF NOT EXISTS articles (
			id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			title           TEXT NOT NULL,
			url             TEXT UNIQUE NOT NULL,
			summary         TEXT,
			source          TEXT NOT NULL,
			team            TEXT NOT NULL,
			author          TEXT,
			image_url       TEXT,
			sentiment_score NUMERIC(5,4),
			published_at    TIMESTAMPTZ NOT NULL,
			created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_articles_team        ON articles(team)`,
		`CREATE INDEX IF NOT EXISTS idx_articles_source      ON articles(source)`,
		`CREATE INDEX IF NOT EXISTS idx_articles_published   ON articles(published_at DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_articles_title_trgm  ON articles USING gin(title gin_trgm_ops)`,
	}
	for _, s := range stmts {
		if _, err := pool.Exec(context.Background(), s); err != nil {
			if s == stmts[len(stmts)-1] {
				pool.Exec(context.Background(), `CREATE EXTENSION IF NOT EXISTS pg_trgm`)
				pool.Exec(context.Background(), s)
				continue
			}
			return fmt.Errorf("migration: %w", err)
		}
	}
	return nil
}
