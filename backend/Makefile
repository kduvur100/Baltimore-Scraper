.PHONY: dev backend frontend setup tidy build docker-up docker-down clean

setup:
	@echo "→ Setting up backend..."
	cd backend && go mod tidy
	@echo "→ Setting up frontend..."
	cd frontend && npm install
	@[ -f .env ] || cp .env.example .env
	@echo "✓ Setup complete. Edit .env then run: make dev"

dev:
	@make -j2 backend frontend

backend:
	cd backend && go run ./cmd/server

frontend:
	cd frontend && npm run dev

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down -v

build-backend:
	cd backend && go build -o bin/server ./cmd/server

tidy:
	cd backend && go mod tidy

check:
	cd backend && go vet ./...
	cd frontend && npm run check

clean:
	rm -rf backend/bin frontend/build frontend/.svelte-kit
