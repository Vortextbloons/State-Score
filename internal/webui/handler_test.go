package webui

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
)

func TestSPAFallback(t *testing.T) {
	t.Parallel()

	dist := fstest.MapFS{
		"dist/index.html":  {Data: []byte("<html>index</html>")},
		"dist/200.html":    {Data: []byte("<html>spa</html>")},
		"dist/_app/app.js": {Data: []byte("console.log(1)")},
		"dist/robots.txt":  {Data: []byte("User-agent: *")},
	}

	h, err := New(dist)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if !h.HasAssets() {
		t.Fatal("expected assets")
	}

	t.Run("index", func(t *testing.T) {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d", rec.Code)
		}
		if body := rec.Body.String(); body != "<html>index</html>" {
			t.Fatalf("body = %q", body)
		}
	})

	t.Run("asset", func(t *testing.T) {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/_app/app.js", nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d", rec.Code)
		}
		if body := rec.Body.String(); body != "console.log(1)" {
			t.Fatalf("body = %q", body)
		}
	})

	t.Run("spa route", func(t *testing.T) {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/rankings", nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d", rec.Code)
		}
		if body := rec.Body.String(); body != "<html>spa</html>" {
			t.Fatalf("body = %q", body)
		}
	})
}
