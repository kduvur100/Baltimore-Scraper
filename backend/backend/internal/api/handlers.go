package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"

	"github.com/kduvur100/baltimore-sports-scraper/internal/models"
)

type Handler struct {
	db  *pgxpool.Pool
	rdb *goredis.Client
}

func NewHandler(db *pgxpool.Pool, rdb *goredis.Client) *Handler {
	return &Handler{db: db, rdb: rdb}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *Handler) GetArticles(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	team := q.Get("team")
	source := q.Get("source")
	limit := clampInt(q.Get("limit"), 20, 100)
	offset := clampInt(q.Get("offset"), 0, 1<<31)

	args := []any{}
	where := "WHERE 1=1"
	i := 1

	if team != "" {
		where += fmt.Sprintf(" AND (team = $%d OR team = 'both')", i)
		args = append(args, team)
		i++
	}
	if source != "" {
		where += fmt.Sprintf(" AND source = $%d", i)
		args = append(args, source)
		i++
	}
	args = append(args, limit, offset)

	rows, err := h.db.Query(r.Context(), `
		SELECT id, title, url, summary, source, team, author,
		       image_url, sentiment_score, published_at, created_at
		FROM articles `+where+`
		ORDER BY published_at DESC
		LIMIT $`+strconv.Itoa(i)+` OFFSET $`+strconv.Itoa(i+1),
		args...)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "failed to query articles")
		return
	}
	defer rows.Close()

	articles := make([]models.Article, 0)
	for rows.Next() {
		var a models.Article
		if err := rows.Scan(&a.ID, &a.Title, &a.URL, &a.Summary, &a.Source, &a.Team,
			&a.Author, &a.ImageURL, &a.SentimentScore, &a.PublishedAt, &a.CreatedAt); err != nil {
			continue
		}
		articles = append(articles, a)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"articles": articles,
		"count":    len(articles),
		"offset":   offset,
	})
}

func (h *Handler) GetArticle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var a models.Article
	err := h.db.QueryRow(r.Context(), `
		SELECT id, title, url, summary, source, team, author,
		       image_url, sentiment_score, published_at, created_at
		FROM articles WHERE id = $1`, id).
		Scan(&a.ID, &a.Title, &a.URL, &a.Summary, &a.Source, &a.Team,
			&a.Author, &a.ImageURL, &a.SentimentScore, &a.PublishedAt, &a.CreatedAt)
	if err != nil {
		writeErr(w, http.StatusNotFound, "article not found")
		return
	}
	writeJSON(w, http.StatusOK, a)
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	team := r.URL.Query().Get("team")
	limit := clampInt(r.URL.Query().Get("limit"), 20, 100)

	if query == "" {
		writeErr(w, http.StatusBadRequest, "q is required")
		return
	}

	like := "%" + query + "%"
	args := []any{like, like}
	extra := ""
	if team != "" {
		extra = " AND (team = $3 OR team = 'both')"
		args = append(args, team, limit)
	} else {
		args = append(args, limit)
	}
	limitArg := "$" + strconv.Itoa(len(args))

	rows, err := h.db.Query(r.Context(), `
		SELECT id, title, url, summary, source, team, author,
		       image_url, sentiment_score, published_at, created_at
		FROM articles
		WHERE (title ILIKE $1 OR summary ILIKE $2)`+extra+`
		ORDER BY published_at DESC
		LIMIT `+limitArg, args...)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "search failed")
		return
	}
	defer rows.Close()

	articles := make([]models.Article, 0)
	for rows.Next() {
		var a models.Article
		if err := rows.Scan(&a.ID, &a.Title, &a.URL, &a.Summary, &a.Source, &a.Team,
			&a.Author, &a.ImageURL, &a.SentimentScore, &a.PublishedAt, &a.CreatedAt); err != nil {
			continue
		}
		articles = append(articles, a)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"articles": articles,
		"query":    query,
		"count":    len(articles),
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func clampInt(s string, def, max int) int {
	if s == "" {
		return def
	}
	v, err := strconv.Atoi(s)
	if err != nil || v < 0 {
		return def
	}
	if v > max {
		return max
	}
	return v
}
