#!/usr/bin/env bash
set -e

GREEN='\033[0;32m'; YELLOW='\033[1;33m'; NC='\033[0m'

echo -e "${GREEN}Baltimore Sports Hub — Setup${NC}\n"

if [ ! -f .env ]; then
  cp .env.example .env
  echo -e "${YELLOW}→ Created .env from .env.example. Edit it if needed.${NC}"
else
  echo "→ .env already exists, skipping."
fi

echo -e "\n→ Installing Go dependencies…"
(cd backend && go mod tidy && go build ./... && echo "  ✓ Go build OK")

echo -e "\n→ Installing Node dependencies…"
(cd frontend && npm install && echo "  ✓ npm install OK")

echo -e "\n${GREEN}✓ Setup complete!${NC}"
echo "  Docker: docker-compose up --build"
echo "  Local:  make dev"
