package webui

import (
	"io/fs"
	"net/http"
	"path"
	"strconv"
	"strings"
)

// Handler serves the embedded SvelteKit static build with SPA fallback routing.
type Handler struct {
	root       fs.FS
	fileServer http.Handler
	index      []byte
	fallback   []byte
	hasAssets  bool
}

// New creates a static UI handler rooted at the embedded web/dist filesystem.
func New(distFS fs.FS) (*Handler, error) {
	root, err := fs.Sub(distFS, "dist")
	if err != nil {
		return nil, err
	}

	h := &Handler{
		root:       root,
		fileServer: http.FileServerFS(root),
	}

	index, err := fs.ReadFile(root, "index.html")
	if err == nil && len(index) > 0 {
		h.index = index
		h.hasAssets = true
	}

	if fb, err := fs.ReadFile(root, "200.html"); err == nil && len(fb) > 0 {
		h.fallback = fb
	} else {
		h.fallback = h.index
	}

	return h, nil
}

// HasAssets reports whether a built frontend was embedded.
func (h *Handler) HasAssets() bool {
	return h.hasAssets
}

// ServeHTTP serves static files and falls back to the SPA shell for client routes.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !h.HasAssets() {
		http.Error(w, "frontend assets are missing; run npm run build in frontend/", http.StatusServiceUnavailable)
		return
	}

	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	reqPath := path.Clean("/" + r.URL.Path)
	if strings.HasPrefix(reqPath, "/api/") {
		http.NotFound(w, r)
		return
	}

	if reqPath == "/" {
		serveBytes(w, r, "text/html; charset=utf-8", h.index)
		return
	}

	name := strings.TrimPrefix(reqPath, "/")
	if info, err := fs.Stat(h.root, name); err == nil && !info.IsDir() {
		h.fileServer.ServeHTTP(w, r)
		return
	}

	serveBytes(w, r, "text/html; charset=utf-8", h.fallback)
}

func serveBytes(w http.ResponseWriter, r *http.Request, contentType string, body []byte) {
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	if r.Method == http.MethodHead {
		w.WriteHeader(http.StatusOK)
		return
	}
	_, _ = w.Write(body)
}
