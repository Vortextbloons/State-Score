package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/isaac/statescore/internal/config"
	"github.com/isaac/statescore/internal/repositories"
)

// Handler serves JSON API routes under /api/v1.
type Handler struct {
	DB        *sql.DB
	StartedAt time.Time

	States      *repositories.StateRepository
	Categories  *repositories.CategoryRepository
	Metrics     *repositories.MetricRepository
	MetricValues *repositories.MetricValueRepository
	Sources     *repositories.DataSourceRepository
	Imports     *repositories.ImportRepository
	Profiles    *repositories.ProfileRepository
}

// NewHandler constructs the API handler.
func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		DB:           db,
		StartedAt:    time.Now().UTC(),
		States:       repositories.NewStateRepository(db),
		Categories:   repositories.NewCategoryRepository(db),
		Metrics:      repositories.NewMetricRepository(db),
		MetricValues: repositories.NewMetricValueRepository(db),
		Sources:      repositories.NewDataSourceRepository(db),
		Imports:      repositories.NewImportRepository(db),
		Profiles:     repositories.NewProfileRepository(db),
	}
}

// Mount registers API routes on mux.
func (h *Handler) Mount(mux *http.ServeMux) {
	// Status endpoints
	mux.HandleFunc("GET /api/v1/health", h.health)
	mux.HandleFunc("GET /api/v1/status", h.status)

	// State endpoints
	mux.HandleFunc("GET /api/v1/states", h.listStates)
	mux.HandleFunc("GET /api/v1/states/{code}", h.getState)
	mux.HandleFunc("GET /api/v1/regions", h.listRegions)

	// Category endpoints
	mux.HandleFunc("GET /api/v1/categories", h.listCategories)

	// Metric endpoints
	mux.HandleFunc("GET /api/v1/metrics", h.listMetrics)
	mux.HandleFunc("GET /api/v1/metrics/{metricId}", h.getMetric)

	// Metric value endpoints
	mux.HandleFunc("GET /api/v1/values", h.listMetricValues)

	// Profile endpoints
	mux.HandleFunc("GET /api/v1/profiles", h.listProfiles)
	mux.HandleFunc("GET /api/v1/profiles/{profileId}", h.getProfile)
	mux.HandleFunc("GET /api/v1/profiles/default", h.getDefaultProfile)

	// Source endpoints
	mux.HandleFunc("GET /api/v1/sources", h.listSources)

	// Import endpoints
	mux.HandleFunc("GET /api/v1/imports", h.listImports)
	mux.HandleFunc("GET /api/v1/imports/{importId}", h.getImport)
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":    "ok",
		"version":   config.Version,
		"startedAt": h.StartedAt.Format(time.RFC3339),
		"app":       config.AppName,
	})
}

func (h *Handler) status(w http.ResponseWriter, r *http.Request) {
	databaseReady := false
	if h.DB != nil {
		if err := h.DB.Ping(); err == nil {
			databaseReady = true
		}
	}

	status := "ready"
	if !databaseReady {
		status = "degraded"
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status":        status,
		"version":       config.Version,
		"databaseReady": databaseReady,
		"activeImports": 0,
		"startedAt":     h.StartedAt.Format(time.RFC3339),
	})
}

func (h *Handler) listStates(w http.ResponseWriter, r *http.Request) {
	region := r.URL.Query().Get("region")

	states, err := h.States.List(region)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "LIST_STATES_FAILED", "Failed to list states", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data": states,
	})
}

func (h *Handler) getState(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		writeError(w, http.StatusBadRequest, "MISSING_STATE_CODE", "State code is required", nil)
		return
	}

	state, err := h.States.GetByCode(code)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "GET_STATE_FAILED", "Failed to get state", err)
		return
	}
	if state == nil {
		writeError(w, http.StatusNotFound, "STATE_NOT_FOUND", "State not found", nil)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data": state,
	})
}

func (h *Handler) listRegions(w http.ResponseWriter, r *http.Request) {
	regions, err := h.States.Regions()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "LIST_REGIONS_FAILED", "Failed to list regions", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data": regions,
	})
}

func (h *Handler) listCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.Categories.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "LIST_CATEGORIES_FAILED", "Failed to list categories", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data": categories,
	})
}

func (h *Handler) listMetrics(w http.ResponseWriter, r *http.Request) {
	categoryIDStr := r.URL.Query().Get("category_id")
	var categoryID int64
	if categoryIDStr != "" {
		id, err := strconv.ParseInt(categoryIDStr, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_CATEGORY_ID", "Invalid category ID", err)
			return
		}
		categoryID = id
	}

	metrics, err := h.Metrics.List(categoryID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "LIST_METRICS_FAILED", "Failed to list metrics", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data": metrics,
	})
}

func (h *Handler) getMetric(w http.ResponseWriter, r *http.Request) {
	metricIDStr := r.PathValue("metricId")
	if metricIDStr == "" {
		writeError(w, http.StatusBadRequest, "MISSING_METRIC_ID", "Metric ID is required", nil)
		return
	}

	metricID, err := strconv.ParseInt(metricIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_METRIC_ID", "Invalid metric ID", err)
		return
	}

	metric, err := h.Metrics.GetByID(metricID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "GET_METRIC_FAILED", "Failed to get metric", err)
		return
	}
	if metric == nil {
		writeError(w, http.StatusNotFound, "METRIC_NOT_FOUND", "Metric not found", nil)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data": metric,
	})
}

func (h *Handler) listMetricValues(w http.ResponseWriter, r *http.Request) {
	stateIDStr := r.URL.Query().Get("state_id")
	yearStr := r.URL.Query().Get("year")

	var stateID int64
	var year int

	if stateIDStr != "" {
		id, err := strconv.ParseInt(stateIDStr, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_STATE_ID", "Invalid state ID", err)
			return
		}
		stateID = id
	}

	if yearStr != "" {
		y, err := strconv.Atoi(yearStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_YEAR", "Invalid year", err)
			return
		}
		year = y
	}

	// For now, require state_id to list values
	if stateID == 0 {
		writeError(w, http.StatusBadRequest, "MISSING_STATE_ID", "State ID is required", nil)
		return
	}

	values, err := h.MetricValues.ListByState(stateID, year)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "LIST_VALUES_FAILED", "Failed to list metric values", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data": values,
	})
}

func (h *Handler) listProfiles(w http.ResponseWriter, r *http.Request) {
	profiles, err := h.Profiles.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "LIST_PROFILES_FAILED", "Failed to list profiles", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data": profiles,
	})
}

func (h *Handler) getProfile(w http.ResponseWriter, r *http.Request) {
	profileIDStr := r.PathValue("profileId")
	if profileIDStr == "" {
		writeError(w, http.StatusBadRequest, "MISSING_PROFILE_ID", "Profile ID is required", nil)
		return
	}

	profileID, err := strconv.ParseInt(profileIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_PROFILE_ID", "Invalid profile ID", err)
		return
	}

	profile, err := h.Profiles.GetByID(profileID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "GET_PROFILE_FAILED", "Failed to get profile", err)
		return
	}
	if profile == nil {
		writeError(w, http.StatusNotFound, "PROFILE_NOT_FOUND", "Profile not found", nil)
		return
	}

	// Get category weights
	categoryWeights, err := h.Profiles.GetCategoryWeights(profileID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "GET_WEIGHTS_FAILED", "Failed to get category weights", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data": map[string]any{
			"profile":          profile,
			"categoryWeights": categoryWeights,
		},
	})
}

func (h *Handler) getDefaultProfile(w http.ResponseWriter, r *http.Request) {
	profile, err := h.Profiles.GetDefault()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "GET_DEFAULT_PROFILE_FAILED", "Failed to get default profile", err)
		return
	}
	if profile == nil {
		writeError(w, http.StatusNotFound, "DEFAULT_PROFILE_NOT_FOUND", "Default profile not found", nil)
		return
	}

	// Get category weights
	categoryWeights, err := h.Profiles.GetCategoryWeights(profile.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "GET_WEIGHTS_FAILED", "Failed to get category weights", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data": map[string]any{
			"profile":          profile,
			"categoryWeights": categoryWeights,
		},
	})
}

func (h *Handler) listSources(w http.ResponseWriter, r *http.Request) {
	sources, err := h.Sources.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "LIST_SOURCES_FAILED", "Failed to list sources", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data": sources,
	})
}

func (h *Handler) listImports(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 50 // Default limit
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err == nil && l > 0 {
			limit = l
		}
	}

	imports, err := h.Imports.List(limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "LIST_IMPORTS_FAILED", "Failed to list imports", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data": imports,
	})
}

func (h *Handler) getImport(w http.ResponseWriter, r *http.Request) {
	importIDStr := r.PathValue("importId")
	if importIDStr == "" {
		writeError(w, http.StatusBadRequest, "MISSING_IMPORT_ID", "Import ID is required", nil)
		return
	}

	importID, err := strconv.ParseInt(importIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_IMPORT_ID", "Invalid import ID", err)
		return
	}

	importRecord, err := h.Imports.GetByID(importID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "GET_IMPORT_FAILED", "Failed to get import", err)
		return
	}
	if importRecord == nil {
		writeError(w, http.StatusNotFound, "IMPORT_NOT_FOUND", "Import not found", nil)
		return
	}

	// Get import errors
	importErrors, err := h.Imports.ListErrors(importID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "GET_IMPORT_ERRORS_FAILED", "Failed to get import errors", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"data": map[string]any{
			"import": importRecord,
			"errors": importErrors,
		},
	})
}

func writeJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, code int, errCode, message string, err error) {
	details := map[string]any{}
	if err != nil {
		details["error"] = err.Error()
	}
	writeJSON(w, code, map[string]any{
		"error": map[string]any{
			"code":    errCode,
			"message": message,
			"details": details,
		},
	})
}
