package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/isaac/statescore/internal/config"
)

// Handler serves Phase 1 HTTP routes (placeholder UI + health).
type Handler struct {
	StartedAt time.Time
}

// NewHandler constructs the API/root handler.
func NewHandler() *Handler {
	return &Handler{StartedAt: time.Now().UTC()}
}

// Mount registers routes on mux.
func (h *Handler) Mount(mux *http.ServeMux) {
	mux.HandleFunc("GET /{$}", h.placeholder)
	mux.HandleFunc("GET /api/v1/health", h.health)
}

func (h *Handler) placeholder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(placeholderHTML))
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":    "ok",
		"version":   config.Version,
		"startedAt": h.StartedAt.Format(time.RFC3339),
		"app":       config.AppName,
	})
}

const placeholderHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>StateScore</title>
  <style>
    :root {
      --ink: #e8eef4;
      --muted: #8b9aab;
      --paper: #0d1218;
      --accent: #3dba8c;
      --line: #2a3542;
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      min-height: 100vh;
      font-family: "Segoe UI", "Helvetica Neue", sans-serif;
      color: var(--ink);
      background:
        radial-gradient(ellipse 80% 50% at 10% 0%, #143028 0%, transparent 55%),
        radial-gradient(ellipse 60% 40% at 100% 100%, #1a2230 0%, transparent 50%),
        var(--paper);
      display: grid;
      place-items: center;
      padding: 2rem;
    }
    main {
      max-width: 36rem;
      width: 100%;
    }
    h1 {
      font-size: clamp(2.5rem, 8vw, 3.75rem);
      font-weight: 700;
      letter-spacing: -0.03em;
      margin: 0 0 0.5rem;
      color: var(--accent);
    }
    p {
      margin: 0 0 1.25rem;
      font-size: 1.125rem;
      line-height: 1.55;
      color: var(--muted);
    }
    .status {
      display: inline-block;
      padding: 0.35rem 0;
      border-top: 1px solid var(--line);
      border-bottom: 1px solid var(--line);
      font-size: 0.875rem;
      letter-spacing: 0.04em;
      text-transform: uppercase;
      color: var(--ink);
    }
  </style>
</head>
<body>
  <main>
    <h1>StateScore</h1>
    <p>Local state comparison is starting up. The full interface arrives in a later phase.</p>
    <div class="status">Application shell running</div>
  </main>
</body>
</html>
`
