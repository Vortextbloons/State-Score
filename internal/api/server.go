package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/isaac/statescore/internal/config"
	"github.com/isaac/statescore/internal/importer"
	"github.com/isaac/statescore/internal/jobs"
	"github.com/isaac/statescore/internal/models"
	"github.com/isaac/statescore/internal/scoring"
)

// Handler serves JSON API routes under /api/v1.
type Handler struct {
	DB        *sql.DB
	StartedAt time.Time

	States       stateStore
	Categories   categoryStore
	Metrics      metricStore
	MetricValues metricValueStore
	Sources      sourceStore
	Imports      importStore
	Profiles     profileStore
	Scores       scoreStore
	Jobs         *jobs.Manager
}

// NewHandler constructs the API handler.
func NewHandler(db *sql.DB, managers ...*jobs.Manager) *Handler {
	return NewHandlerWithDependencies(db, repositoryDependencies(db), managers...)
}

// NewHandlerWithDependencies constructs a handler from explicit boundary dependencies.
func NewHandlerWithDependencies(db *sql.DB, deps Dependencies, managers ...*jobs.Manager) *Handler {
	manager := jobs.New(context.Background())
	if len(managers) > 0 && managers[0] != nil {
		manager = managers[0]
	}
	return &Handler{
		DB:           db,
		StartedAt:    time.Now().UTC(),
		States:       deps.States,
		Categories:   deps.Categories,
		Metrics:      deps.Metrics,
		MetricValues: deps.MetricValues,
		Sources:      deps.Sources,
		Imports:      deps.Imports,
		Profiles:     deps.Profiles,
		Scores:       deps.Scores,
		Jobs:         manager,
	}
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
		"activeImports": h.Jobs.Active(),
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

	var values []models.MetricValue
	var err error
	if stateID == 0 {
		values, err = h.MetricValues.ListAll(year)
	} else {
		values, err = h.MetricValues.ListByState(stateID, year)
	}
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
			"profile":         profile,
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
			"profile":         profile,
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

func (h *Handler) createSource(w http.ResponseWriter, r *http.Request) {
	var source models.DataSource
	if err := decodeJSON(r, &source); err != nil {
		writeError(w, 400, "INVALID_SOURCE", "Source details are invalid", err)
		return
	}
	if strings.TrimSpace(source.Name) == "" {
		writeError(w, 422, "SOURCE_NAME_REQUIRED", "Source name is required", nil)
		return
	}
	if source.Format == "" {
		source.Format = "csv"
	}
	if !strings.EqualFold(source.Format, "csv") {
		writeError(w, 422, "UNSUPPORTED_SOURCE_FORMAT", "Phase 6 sources must use CSV format", nil)
		return
	}
	if err := h.Sources.Create(&source); err != nil {
		writeError(w, 500, "CREATE_SOURCE_FAILED", "Failed to create source", err)
		return
	}
	writeJSON(w, 201, map[string]any{"data": source})
}

func (h *Handler) updateSource(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("sourceId"), 10, 64)
	if err != nil {
		writeError(w, 400, "INVALID_SOURCE_ID", "Invalid source ID", err)
		return
	}
	existing, err := h.Sources.GetByID(id)
	if err != nil {
		writeError(w, 500, "GET_SOURCE_FAILED", "Failed to get source", err)
		return
	}
	if existing == nil {
		writeError(w, 404, "SOURCE_NOT_FOUND", "Source not found", nil)
		return
	}
	var source models.DataSource
	if err := decodeJSON(r, &source); err != nil {
		writeError(w, 400, "INVALID_SOURCE", "Source details are invalid", err)
		return
	}
	source.ID = id
	if strings.TrimSpace(source.Name) == "" {
		writeError(w, 422, "SOURCE_NAME_REQUIRED", "Source name is required", nil)
		return
	}
	if source.Format == "" {
		source.Format = "csv"
	}
	if err := h.Sources.Update(&source); err != nil {
		writeError(w, 500, "UPDATE_SOURCE_FAILED", "Failed to update source", err)
		return
	}
	writeJSON(w, 200, map[string]any{"data": source})
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

func (h *Handler) createImport(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, importer.MaxUploadBytes)
	if err := r.ParseMultipartForm(importer.MaxUploadBytes); err != nil {
		writeError(w, 400, "INVALID_IMPORT", "Upload a CSV file smaller than 10 MB", err)
		return
	}
	sourceID, err := strconv.ParseInt(r.FormValue("source_id"), 10, 64)
	if err != nil {
		writeError(w, 422, "SOURCE_REQUIRED", "Select a data source", err)
		return
	}
	if source, err := h.Sources.GetByID(sourceID); err != nil || source == nil {
		writeError(w, 422, "INVALID_SOURCE", "The selected source does not exist", err)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, 422, "FILE_REQUIRED", "Choose a CSV file", err)
		return
	}
	defer file.Close()
	if !strings.HasSuffix(strings.ToLower(header.Filename), ".csv") {
		writeError(w, 422, "CSV_REQUIRED", "The import file must use the .csv extension", nil)
		return
	}
	content, err := io.ReadAll(io.LimitReader(file, importer.MaxUploadBytes+1))
	if err != nil || len(content) > importer.MaxUploadBytes {
		writeError(w, 400, "FILE_TOO_LARGE", "The CSV file exceeds 10 MB", err)
		return
	}
	record := models.Import{SourceID: &sourceID, Status: "pending"}
	if err := h.Imports.Create(&record); err != nil {
		writeError(w, 500, "CREATE_IMPORT_FAILED", "Failed to create import", err)
		return
	}
	h.Jobs.Go(func(ctx context.Context) { h.runImport(ctx, record, content) })
	writeJSON(w, 202, map[string]any{"data": record, "meta": map[string]any{"fileName": header.Filename}})
}

func (h *Handler) runImport(ctx context.Context, record models.Import, content []byte) {
	now := time.Now().UTC().Format(time.RFC3339)
	record.StartedAt = &now
	record.Status = "running"
	_ = h.Imports.Update(&record)
	result, err := importer.CSV(ctx, h.DB, record.ID, content)
	record.RecordsRead = result.RecordsRead
	record.RecordsInserted = result.RecordsInserted
	record.RecordsRejected = result.RecordsRejected
	record.Checksum = result.Checksum
	for _, e := range result.Errors {
		row := e.Row
		_ = h.Imports.AddError(&models.ImportError{ImportID: record.ID, RowNumber: &row, FieldName: e.Field, RawValue: e.Value, ErrorMessage: e.Message})
	}
	done := time.Now().UTC().Format(time.RFC3339)
	record.CompletedAt = &done
	if err != nil {
		record.Status = "failed"
		record.ErrorSummary = err.Error()
	} else if result.RecordsInserted == 0 {
		record.Status = "failed"
		record.ErrorSummary = "No valid records were found"
	} else if result.RecordsRejected > 0 {
		record.Status = "completed_with_errors"
		record.ErrorSummary = fmt.Sprintf("%d row(s) rejected", result.RecordsRejected)
	} else {
		record.Status = "completed"
	}
	_ = h.Imports.Update(&record)
	if record.Status == "completed" || record.Status == "completed_with_errors" {
		if p, e := h.Profiles.GetDefault(); e == nil && p != nil {
			years, _ := h.MetricValues.AvailableYears()
			for _, year := range years {
				if ctx.Err() != nil {
					return
				}
				_, _ = scoring.Recalculate(ctx, h.DB, p.ID, year)
			}
		}
	}
}

func (h *Handler) listScores(w http.ResponseWriter, r *http.Request) {
	profileID, year, err := h.resolveScoreQuery(r)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, "SCORE_QUERY_INVALID", err.Error(), err)
		return
	}
	if err := h.ensureSnapshots(r.Context(), profileID, year); err != nil {
		writeError(w, http.StatusInternalServerError, "SCORE_ENSURE_FAILED", "Failed to prepare score snapshots", err)
		return
	}
	rows, err := h.Scores.ListByProfileYear(profileID, year, scoring.CalculationVersion)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "LIST_SCORES_FAILED", "Failed to list scores", err)
		return
	}
	scores := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		categories := make([]map[string]any, 0, len(row.Categories))
		for _, c := range row.Categories {
			categories = append(categories, map[string]any{
				"categoryId":   c.CategoryID,
				"score":        c.Score,
				"completeness": c.Completeness,
			})
		}
		scores = append(scores, map[string]any{
			"stateId":            row.Snapshot.StateID,
			"overallScore":       row.Snapshot.OverallScore,
			"completeness":       row.Snapshot.Completeness,
			"calculatedAt":       row.Snapshot.CalculatedAt.UTC().Format(time.RFC3339),
			"calculationVersion": row.Snapshot.CalculationVersion,
			"categories":         categories,
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"data": map[string]any{
			"profileId":          profileID,
			"year":               year,
			"asOfYear":           year,
			"calculationVersion": scoring.CalculationVersion,
			"scores":             scores,
		},
	})
}

func (h *Handler) recalculateScores(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ProfileID int64 `json:"profileId"`
		Year      int   `json:"year"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, 400, "INVALID_RECALCULATION", "Invalid recalculation request", err)
		return
	}
	if input.ProfileID == 0 {
		p, err := h.Profiles.GetDefault()
		if err != nil || p == nil {
			writeError(w, 422, "PROFILE_REQUIRED", "No default scoring profile is available", err)
			return
		}
		input.ProfileID = p.ID
	}
	if input.Year == 0 {
		years, err := h.MetricValues.AvailableYears()
		if err != nil || len(years) == 0 {
			writeError(w, 422, "YEAR_REQUIRED", "No imported data year is available", err)
			return
		}
		input.Year = years[0]
	}
	count, err := scoring.Recalculate(r.Context(), h.DB, input.ProfileID, input.Year)
	if err != nil {
		writeError(w, 500, "RECALCULATION_FAILED", "Failed to recalculate scores", err)
		return
	}
	writeJSON(w, 200, map[string]any{"data": map[string]any{"profileId": input.ProfileID, "year": input.Year, "statesCalculated": count}})
}

func (h *Handler) resolveScoreQuery(r *http.Request) (profileID int64, year int, err error) {
	if v := r.URL.Query().Get("profile_id"); v != "" {
		profileID, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid profile_id")
		}
	}
	if v := r.URL.Query().Get("year"); v != "" {
		year, err = strconv.Atoi(v)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid year")
		}
	}
	if profileID == 0 {
		p, e := h.Profiles.GetDefault()
		if e != nil || p == nil {
			return 0, 0, fmt.Errorf("no default scoring profile is available")
		}
		profileID = p.ID
	}
	if year == 0 {
		years, e := h.MetricValues.AvailableYears()
		if e != nil || len(years) == 0 {
			return 0, 0, fmt.Errorf("no imported data year is available")
		}
		year = years[0]
	}
	return profileID, year, nil
}

func (h *Handler) ensureSnapshots(ctx context.Context, profileID int64, year int) error {
	ok, err := h.Scores.HasSnapshots(profileID, year, scoring.CalculationVersion)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	_, err = scoring.Recalculate(ctx, h.DB, profileID, year)
	return err
}

func decodeJSON(r *http.Request, dst any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
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
