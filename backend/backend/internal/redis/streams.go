package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"github.com/kduvur100/baltimore-sports-scraper/internal/models"
)

const (
	StreamKey    = "articles:stream"
	StreamMaxLen = 1000
)

func NewClient(redisURL string) *goredis.Client {
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}
	opts, err := goredis.ParseURL(redisURL)
	if err != nil {
		log.Printf("Redis URL parse error (%v), using localhost:6379", err)
		opts = &goredis.Options{Addr: "localhost:6379"}
	}
	return goredis.NewClient(opts)
}

func Publish(ctx context.Context, rdb *goredis.Client, a models.Article) error {
	data, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	return rdb.XAdd(ctx, &goredis.XAddArgs{
		Stream: StreamKey,
		MaxLen: StreamMaxLen,
		Approx: true,
		Values: map[string]interface{}{
			"article": string(data),
			"team":    string(a.Team),
			"source":  string(a.Source),
		},
	}).Err()
}

func Subscribe(ctx context.Context, rdb *goredis.Client, lastID string, ch chan<- models.Article) {
	if lastID == "" {
		lastID = "$"
	}
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		streams, err := rdb.XRead(ctx, &goredis.XReadArgs{
			Streams: []string{StreamKey, lastID},
			Count:   10,
			Block:   5 * time.Second,
		}).Result()

		if err == goredis.Nil {
			continue
		}
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("Redis stream read: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		for _, stream := range streams {
			for _, msg := range stream.Messages {
				lastID = msg.ID
				raw, _ := msg.Values["article"].(string)
				if raw == "" {
					continue
				}
				var a models.Article
				if err := json.Unmarshal([]byte(raw), &a); err != nil {
					continue
				}
				select {
				case ch <- a:
				case <-ctx.Done():
					return
				}
			}
		}
	}
}
