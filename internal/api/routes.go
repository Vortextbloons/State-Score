package api

import "net/http"

// Mount registers feature routes under the versioned API prefix.
func (h *Handler) Mount(mux *http.ServeMux) {
	h.mountOperations(mux)
	h.mountCatalogs(mux)
	h.mountSourcesAndImports(mux)
	h.mountScores(mux)
}

func (h *Handler) mountOperations(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/health", h.health)
	mux.HandleFunc("GET /api/v1/status", h.status)
}

func (h *Handler) mountCatalogs(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/states", h.listStates)
	mux.HandleFunc("GET /api/v1/states/{code}", h.getState)
	mux.HandleFunc("GET /api/v1/regions", h.listRegions)
	mux.HandleFunc("GET /api/v1/categories", h.listCategories)
	mux.HandleFunc("GET /api/v1/metrics", h.listMetrics)
	mux.HandleFunc("GET /api/v1/metrics/{metricId}", h.getMetric)
	mux.HandleFunc("GET /api/v1/values", h.listMetricValues)
	mux.HandleFunc("GET /api/v1/profiles", h.listProfiles)
	mux.HandleFunc("GET /api/v1/profiles/{profileId}", h.getProfile)
	mux.HandleFunc("GET /api/v1/profiles/default", h.getDefaultProfile)
}

func (h *Handler) mountSourcesAndImports(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/sources", h.listSources)
	mux.HandleFunc("POST /api/v1/sources", h.createSource)
	mux.HandleFunc("PUT /api/v1/sources/{sourceId}", h.updateSource)
	mux.HandleFunc("GET /api/v1/imports", h.listImports)
	mux.HandleFunc("GET /api/v1/imports/{importId}", h.getImport)
	mux.HandleFunc("POST /api/v1/imports", h.createImport)
}

func (h *Handler) mountScores(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/scores", h.listScores)
	mux.HandleFunc("POST /api/v1/scores/recalculate", h.recalculateScores)
}
