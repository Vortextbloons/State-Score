.PHONY: dev-frontend dev-backend build frontend-build go-build test tidy

dev-frontend:
	cd frontend && npm run dev

dev-backend:
	go run ./cmd/statescore

frontend-build:
	cd frontend && npm run build

go-build:
	go build -o bin/statescore ./cmd/statescore

build: frontend-build go-build

test:
	go test ./cmd/... ./internal/...
	cd frontend && npm test

tidy:
	go mod tidy
	cd frontend && npm install
