package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"github.com/kduvur100/baltimore-sports-scraper/internal/api"
	"github.com/kduvur100/baltimore-sports-scraper/internal/db"
	redisClient "github.com/kduvur100/baltimore-sports-scraper/internal/redis"
	"github.com/kduvur100/baltimore-sports-scraper/internal/scraper"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	pool, err := db.NewPool(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("DB connect: %v", err)
	}
	defer pool.Close()

	if err := db.RunMigrations(pool); err != nil {
		log.Fatalf("DB migrations: %v", err)
	}

	rdb := redisClient.NewClient(os.Getenv("REDIS_URL"))
	defer rdb.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go scraper.NewManager(pool, rdb).Start(ctx)

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{frontendURL, "http://localhost:5173"},
		AllowedMethods: []string{"GET", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Content-Type"},
	}))

	h := api.NewHandler(pool, rdb)
	r.Get("/health", h.Health)
	r.Get("/api/articles", h.GetArticles)
	r.Get("/api/articles/{id}", h.GetArticle)
	r.Get("/api/search", h.Search)
	r.Get("/api/stream", h.SSEStream)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 0,
		IdleTimeout:  60 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Server listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down…")
	cancel()

	sdCtx, sdCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer sdCancel()
	_ = srv.Shutdown(sdCtx)
	log.Println("Done.")
}
