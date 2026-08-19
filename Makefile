.PHONY: all dev test build docker-build docker-up docker-down

VERSION ?= 0.1.0
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S' 2>/dev/null || echo "unknown")

all: build

# Run local development server
dev:
	go run ./cmd/netip server

# Run backend and frontend tests
test:
	go test -v ./...
	cd web && npm run test

# Build frontend and single Go executable
build:
	cd web && npm run build
	go build -ldflags "-s -w -X netip/internal/config.Version=$(VERSION) -X netip/internal/config.Commit=$(COMMIT) -X netip/internal/config.BuildTime=$(BUILD_TIME)" -o bin/netip ./cmd/netip

# Docker deployment helpers
docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down
