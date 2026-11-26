# Baltimore Sports Hub 

Real-time Baltimore Orioles & Ravens news aggregator. Scrapes GDELT, ESPN, MLB, NFL, and Baltimore Sun RSS feeds and pushes live updates to the browser via Server-Sent Events.

**Stack:** Go (chi router) · PostgreSQL · Redis Streams · SvelteKit · Vercel · Fly.io

---

## Project Structure

    Baltimore_Sports_Scraper/
    ├── backend/
    │   ├── cmd/server/main.go
    │   ├── internal/
    │   │   ├── api/
    │   │   │   ├── handlers.go        REST endpoints
    │   │   │   └── sse.go             Server-Sent Events stream
    │   │   ├── db/
    │   │   │   └── postgres.go        pgxpool + auto-migrations
    │   │   ├── models/
    │   │   │   └── article.go         Shared Article struct
    │   │   ├── redis/
    │   │   │   └── streams.go         Publish / Subscribe
    │   │   └── scraper/
    │   │       ├── manager.go         Runs all scrapers concurrently
    │   │       ├── rss.go             ESPN, MLB, NFL, Baltimore Sun
    │   │       ├── gdelt.go           GDELT DOC 2.0 API (free, no key)
    │   │       ├── reddit.go          Stub (activate when ready)
    │   │       └── bluesky.go         Stub (activate when ready)
    │   ├── go.mod
    │   └── Dockerfile
    ├── frontend/
    │   ├── src/
    │   │   ├── routes/
    │   │   │   ├── +layout.svelte     Topbar + live indicator
    │   │   │   ├── +layout.ts         SSR config
    │   │   │   ├── +page.svelte       Article grid
    │   │   │   ├── +page.server.ts    SSR data load
    │   │   │   └── +error.svelte      Error page
    │   │   └── lib/
    │   │       ├── stores/
    │   │       │   └── articles.ts    Svelte stores + derived filters
    │   │       └── components/
    │   │           ├── ArticleCard.svelte
    │   │           ├── SearchBar.svelte
    │   │           └── LiveFeed.svelte  Headless SSE client
    │   ├── static/favicon.svg
    │   ├── package.json
    │   ├── svelte.config.js
    │   ├── tsconfig.json
    │   └── vite.config.ts
    ├── docker-compose.yml
    ├── fly.toml
    ├── Makefile
    ├── setup.sh
    └── .env.example

---

## Quick Start (Docker)

Requires Docker Desktop.

    git clone https://github.com/kduvur100/Baltimore_Sports_Scraper.git
    cd Baltimore_Sports_Scraper
    cp .env.example .env
    docker-compose up --build

- Frontend: http://localhost:5173
- Backend: http://localhost:8080
- Health: http://localhost:8080/health

---

## Local Dev (without Docker)

Prerequisites: Go 1.22+, Node 20+, PostgreSQL 16, Redis 7

    bash setup.sh   # installs deps, copies .env
    nano .env       # fill in your DB/Redis URLs
    make dev        # starts backend :8080 + frontend :5173

---

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | /health | Health check |
| GET | /api/articles | List articles (paginated) |
| GET | /api/articles/{id} | Single article |
| GET | /api/search?q=lamar | Full-text search |
| GET | /api/stream | SSE live feed |

Query params for /api/articles: team=orioles|ravens, source=rss|gdelt, limit=20, offset=0

---

## Environment Variables

    PORT=8080
    DATABASE_URL=postgres://postgres:postgres@localhost:5432/baltimore_sports?sslmode=disable
    REDIS_URL=redis://localhost:6379
    FRONTEND_URL=http://localhost:5173

    # Production (set in Vercel dashboard):
    # API_URL=https://baltimore-sports-backend.fly.dev

---

## Deployment

### Backend to Fly.io

    fly auth login
    fly launch --no-deploy
    fly secrets set DATABASE_URL="your-prod-postgres-url"
    fly secrets set REDIS_URL="your-prod-redis-url"
    fly secrets set FRONTEND_URL="https://your-app.vercel.app"
    fly deploy

### Frontend to Vercel

    cd frontend
    npx vercel --prod
    # Then in Vercel dashboard > Environment Variables:
    # API_URL = https://baltimore-sports-backend.fly.dev

---

## Enabling Optional Scrapers

**Reddit:**
1. Create app at https://www.reddit.com/prefs/apps
2. Add to .env: REDDIT_CLIENT_ID, REDDIT_CLIENT_SECRET, REDDIT_USER_AGENT
3. Uncomment NewRedditScraper() in backend/internal/scraper/manager.go

**Bluesky:**
1. Create app password at https://bsky.app > Settings > App Passwords
2. Add to .env: BLUESKY_IDENTIFIER, BLUESKY_APP_PASSWORD
3. Uncomment NewBlueskyScraper() in backend/internal/scraper/manager.go

---

## Makefile Commands

    make setup        # Install all deps + copy .env
    make dev          # Run backend + frontend concurrently
    make docker-up    # Start everything in Docker
    make docker-down  # Stop Docker + remove volumes
    make tidy         # go mod tidy
    make check        # go vet + svelte-check
    make clean        # Remove build artifacts

---

## How It Works

1. Scrapers run as goroutines on a ticker (RSS every 10 min, GDELT every 15 min)
2. New articles are upserted into PostgreSQL — deduplicated by URL
3. Each new article is published to a Redis Stream
4. The SSE endpoint subscribes to that stream and pushes JSON to browsers
5. SvelteKit frontend receives events and prepends them to the article store live
6. Articles are SSR-seeded from the database on first load for instant first paint

---

## Cleanup

    docker-compose down -v    # stop containers + wipe DB data
    make clean                # remove build artifacts
    rm -rf frontend/node_modules
    go clean -cache
READMEEOF
