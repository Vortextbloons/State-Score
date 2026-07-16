package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/isaac/statescore/internal/database"
)

func testHandler(t *testing.T) http.Handler {
	t.Helper()
	db, err := database.Open(filepath.Join(t.TempDir(), "api.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := database.Migrate(db); err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	NewHandler(db).Mount(mux)
	return mux
}

func TestHealthContract(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	w := httptest.NewRecorder()
	testHandler(t).ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["status"] != "ok" || body["app"] == nil {
		t.Fatalf("unexpected response: %v", body)
	}
}

func TestValuesSupportsBulkAndRejectsInvalidYear(t *testing.T) {
	h := testHandler(t)
	for _, tc := range []struct {
		path   string
		status int
	}{{"/api/v1/values", 200}, {"/api/v1/values?year=nope", 400}} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, tc.path, nil))
		if w.Code != tc.status {
			t.Fatalf("%s: status = %d, want %d", tc.path, w.Code, tc.status)
		}
	}
}

func TestPublicSourceCatalogAndRefreshValidation(t *testing.T) {
	h := testHandler(t)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/public-sources", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("catalog status=%d", w.Code)
	}

	w = httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/api/v1/public-sources/refresh", strings.NewReader(`{"adapterIds":["missing"],"year":2024}`))
	r.Header.Set("Content-Type", "application/json")
	h.ServeHTTP(w, r)
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("refresh status=%d", w.Code)
	}
}
