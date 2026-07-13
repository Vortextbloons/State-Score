// Package importer validates and atomically imports StateScore CSV datasets.
package importer

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

const MaxUploadBytes = 10 << 20

type ValidationError struct {
	Row                   int
	Field, Value, Message string
}
type Result struct {
	RecordsRead, RecordsInserted, RecordsRejected int
	Checksum                                      string
	Errors                                        []ValidationError
}
type row struct {
	stateID, metricID int64
	year              int
	value             float64
	recordID          string
}

// CSV imports a long-form dataset. Validation is completed before the transaction
// starts, so malformed files never partially replace usable data.
func CSV(ctx context.Context, db *sql.DB, importID int64, content []byte) (Result, error) {
	result := Result{Checksum: fmt.Sprintf("%x", sha256.Sum256(content))}
	r := csv.NewReader(bytes.NewReader(content))
	r.TrimLeadingSpace = true
	r.ReuseRecord = false
	header, err := r.Read()
	if err != nil {
		return result, fmt.Errorf("read CSV header: %w", err)
	}
	columns := map[string]int{}
	for i, name := range header {
		columns[strings.ToLower(strings.TrimSpace(name))] = i
	}
	for _, required := range []string{"state_code", "metric_slug", "year", "value"} {
		if _, ok := columns[required]; !ok {
			return result, fmt.Errorf("missing required column %q", required)
		}
	}
	states, err := lookup(ctx, db, `SELECT upper(code), id FROM states`)
	if err != nil {
		return result, err
	}
	metrics, err := lookup(ctx, db, `SELECT slug, id FROM metrics WHERE active=1`)
	if err != nil {
		return result, err
	}
	var valid []row
	seen := map[string]bool{}
	for line := 2; ; line++ {
		record, readErr := r.Read()
		if errors.Is(readErr, io.EOF) {
			break
		}
		result.RecordsRead++
		if readErr != nil {
			result.Errors = append(result.Errors, ValidationError{line, "row", "", readErr.Error()})
			continue
		}
		get := func(name string) string {
			i := columns[name]
			if i >= len(record) {
				return ""
			}
			return strings.TrimSpace(record[i])
		}
		code, slug, yearRaw, valueRaw := strings.ToUpper(get("state_code")), get("metric_slug"), get("year"), get("value")
		stateID, stateOK := states[code]
		metricID, metricOK := metrics[slug]
		year, yearErr := strconv.Atoi(yearRaw)
		value, valueErr := strconv.ParseFloat(valueRaw, 64)
		var errs []ValidationError
		if !stateOK {
			errs = append(errs, ValidationError{line, "state_code", code, "Unknown two-letter state code"})
		}
		if !metricOK {
			errs = append(errs, ValidationError{line, "metric_slug", slug, "Unknown or inactive metric"})
		}
		if yearErr != nil || year < 1900 || year > time.Now().Year() {
			errs = append(errs, ValidationError{line, "year", yearRaw, "Year must be between 1900 and the current year"})
		}
		if valueErr != nil {
			errs = append(errs, ValidationError{line, "value", valueRaw, "Value must be numeric"})
		}
		key := fmt.Sprintf("%s|%s|%d", code, slug, year)
		if seen[key] {
			errs = append(errs, ValidationError{line, "row", key, "Duplicate state, metric, and year in this file"})
		}
		seen[key] = true
		if len(errs) > 0 {
			result.Errors = append(result.Errors, errs...)
			continue
		}
		valid = append(valid, row{stateID, metricID, year, value, get("source_record_id")})
	}
	result.RecordsRejected = result.RecordsRead - len(valid)
	if len(valid) == 0 {
		return result, nil
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return result, err
	}
	defer tx.Rollback()
	for _, v := range valid {
		_, err = tx.ExecContext(ctx, `INSERT INTO metric_values(state_id,metric_id,year,value,source_record_id,import_id) VALUES(?,?,?,?,?,?)`, v.stateID, v.metricID, v.year, v.value, v.recordID, importID)
		if err != nil {
			return result, fmt.Errorf("insert metric value: %w", err)
		}
	}
	if err = tx.Commit(); err != nil {
		return result, err
	}
	result.RecordsInserted = len(valid)
	return result, nil
}

func lookup(ctx context.Context, db *sql.DB, query string) (map[string]int64, error) {
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]int64{}
	for rows.Next() {
		var key string
		var id int64
		if err := rows.Scan(&key, &id); err != nil {
			return nil, err
		}
		out[key] = id
	}
	return out, rows.Err()
}
