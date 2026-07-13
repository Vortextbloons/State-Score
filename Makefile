.PHONY: dev dev-frontend dev-backend dev-backend-live build frontend-build go-build test tidy

dev:
ifeq ($(OS),Windows_NT)
	powershell -Command "$env:STATESCORE_PORT='8080'; $env:STATESCORE_NO_BROWSER='1'; Start-Process cmd -ArgumentList '/c cd /d frontend && npm run dev' -WindowStyle Minimized; air"
else
	cd frontend && npm run dev &
	STATESCORE_PORT=8080 STATESCORE_NO_BROWSER=1 air
endif

dev-frontend:
	cd frontend && npm run dev

# Dual-server development: Vite proxies /api to this process on :8080.
dev-backend:
ifeq ($(OS),Windows_NT)
	cmd /C "set STATESCORE_PORT=8080&& set STATESCORE_NO_BROWSER=1&& go run ./cmd/statescore"
else
	STATESCORE_PORT=8080 STATESCORE_NO_BROWSER=1 go run ./cmd/statescore
endif

# Hot-reload backend using air (install: go install github.com/air-verse/air@latest)
dev-backend-live:
ifeq ($(OS),Windows_NT)
	cmd /C "set STATESCORE_PORT=8080&& set STATESCORE_NO_BROWSER=1&& air"
else
	STATESCORE_PORT=8080 STATESCORE_NO_BROWSER=1 air
endif

frontend-build:
	cd frontend && npm run build

go-build:
	go build -o bin/statescore$(if $(filter Windows_NT,$(OS)),.exe,) ./cmd/statescore

build: frontend-build go-build

test:
	go test ./cmd/... ./internal/... ./web/...
	cd frontend && npm test

tidy:
	go mod tidy
	cd frontend && npm install
