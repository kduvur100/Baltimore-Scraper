package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kduvur100/baltimore-sports-scraper/internal/models"
	redisClient "github.com/kduvur100/baltimore-sports-scraper/internal/redis"
)

func (h *Handler) SSEStream(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	teamFilter := r.URL.Query().Get("team")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	fmt.Fprintf(w, "event: connected\ndata: {\"status\":\"connected\"}\n\n")
	flusher.Flush()

	ch := make(chan models.Article, 16)
	ctx := r.Context()

	go redisClient.Subscribe(ctx, h.rdb, "$", ch)

	keepAlive := time.NewTicker(20 * time.Second)
	defer keepAlive.Stop()

	log.Printf("SSE client connected %s (filter=%q)", r.RemoteAddr, teamFilter)

	for {
		select {
		case <-ctx.Done():
			log.Printf("SSE client disconnected %s", r.RemoteAddr)
			return

		case a := <-ch:
			if teamFilter != "" && string(a.Team) != teamFilter && a.Team != models.TeamBoth {
				continue
			}
			data, err := json.Marshal(a)
			if err != nil {
				continue
			}
			fmt.Fprintf(w, "event: article\ndata: %s\n\n", data)
			flusher.Flush()

		case <-keepAlive.C:
			fmt.Fprintf(w, ": keepalive\n\n")
			flusher.Flush()
		}
	}
}
