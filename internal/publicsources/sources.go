// Package publicsources provides extensible adapters for refreshing official datasets.
package publicsources

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Spec struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Publisher         string   `json:"publisher"`
	MetricSlugs       []string `json:"metricSlugs"`
	DefaultYear       int      `json:"defaultYear"`
	Available         bool     `json:"available"`
	UnavailableReason string   `json:"unavailableReason,omitempty"`
}

type Quality struct {
	ReportingCoverage     *float64
	ParticipatingAgencies *int
	PopulationCovered     *int64
	DataRevision          string
	ScoringEligible       bool
	ExclusionReason       string
}

type Observation struct {
	StateCode, MetricSlug, SourceRecordID string
	Year                                  int
	Value                                 float64
	Quality                               *Quality
}

type Batch struct{ Observations []Observation }

// Adapter is the only extension point required for a new public source.
type Adapter interface {
	Spec() Spec
	SourceName() string
	Fetch(context.Context, int) (Batch, error)
}

type Registry struct{ adapters map[string]Adapter }

func NewRegistry(adapters ...Adapter) *Registry {
	r := &Registry{adapters: map[string]Adapter{}}
	for _, adapter := range adapters {
		r.adapters[adapter.Spec().ID] = adapter
	}
	return r
}

func DefaultRegistry() *Registry {
	client := &http.Client{Timeout: 45 * time.Second}
	key := censusKey()
	return NewRegistry(
		&censusAdapter{client, key, "census-college-enrollment", "ACS 2024 Subject Table S1401", "young-adult-college-enrollment", "subject"},
		&censusAdapter{client, key, "census-renter-burden", "ACS 2024 Data Profile DP04", "renter-housing-cost-burden", "profile"},
		&cdcAdapter{client}, &blsAdapter{client}, &fbiAdapter{client},
	)
}

func (r *Registry) List() []Spec {
	out := make([]Spec, 0, len(r.adapters))
	for _, adapter := range r.adapters {
		out = append(out, adapter.Spec())
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func (r *Registry) Get(id string) (Adapter, bool) { a, ok := r.adapters[id]; return a, ok }

type Service struct {
	db       *sql.DB
	registry *Registry
}

func NewService(db *sql.DB, registry *Registry) *Service { return &Service{db, registry} }
func (s *Service) Specs() []Spec                         { return s.registry.List() }

func (s *Service) Prepare(adapterID string) (int64, error) {
	a, ok := s.registry.Get(adapterID)
	if !ok {
		return 0, fmt.Errorf("unknown source adapter %q", adapterID)
	}
	var sourceID int64
	if err := s.db.QueryRow(`SELECT id FROM data_sources WHERE name=?`, a.SourceName()).Scan(&sourceID); err != nil {
		return 0, fmt.Errorf("find adapter source: %w", err)
	}
	result, err := s.db.Exec(`INSERT INTO imports(source_id,status) VALUES(?,'pending')`, sourceID)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (s *Service) Run(ctx context.Context, adapterID string, year int, importID int64) (runErr error) {
	a, ok := s.registry.Get(adapterID)
	if !ok {
		return fmt.Errorf("unknown source adapter %q", adapterID)
	}
	spec := a.Spec()
	if !spec.Available {
		return s.fail(importID, spec.UnavailableReason)
	}
	if year == 0 {
		year = spec.DefaultYear
	}
	started := time.Now().UTC().Format(time.RFC3339)
	_, _ = s.db.Exec(`UPDATE imports SET status='running',started_at=? WHERE id=?`, started, importID)
	defer func() {
		if runErr != nil {
			_ = s.fail(importID, runErr.Error())
		}
	}()
	batch, err := a.Fetch(ctx, year)
	if err != nil {
		_ = s.fail(importID, err.Error())
		return err
	}
	if len(batch.Observations) == 0 {
		err = fmt.Errorf("source returned no observations")
		_ = s.fail(importID, err.Error())
		return err
	}
	encoded, _ := json.Marshal(batch.Observations)
	checksum := fmt.Sprintf("%x", sha256.Sum256(encoded))
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	inserted, rejected := 0, 0
	for index, obs := range batch.Observations {
		var stateID, metricID int64
		if err = tx.QueryRowContext(ctx, `SELECT id FROM states WHERE code=?`, strings.ToUpper(obs.StateCode)).Scan(&stateID); err != nil {
			rejected++
			s.addErrorTx(tx, importID, index+1, "state_code", obs.StateCode, "Unknown state")
			continue
		}
		if err = tx.QueryRowContext(ctx, `SELECT id FROM metrics WHERE slug=? AND active=1`, obs.MetricSlug).Scan(&metricID); err != nil {
			rejected++
			s.addErrorTx(tx, importID, index+1, "metric_slug", obs.MetricSlug, "Unknown or inactive metric")
			continue
		}
		result, insertErr := tx.ExecContext(ctx, `INSERT INTO metric_values(state_id,metric_id,year,value,source_record_id,import_id) VALUES(?,?,?,?,?,?)`, stateID, metricID, obs.Year, obs.Value, obs.SourceRecordID, importID)
		if insertErr != nil {
			return insertErr
		}
		valueID, _ := result.LastInsertId()
		inserted++
		if q := obs.Quality; q != nil {
			_, err = tx.ExecContext(ctx, `INSERT INTO metric_value_quality(metric_value_id,reporting_coverage,participating_agencies,population_covered,data_revision,scoring_eligible,exclusion_reason) VALUES(?,?,?,?,?,?,?)`, valueID, q.ReportingCoverage, q.ParticipatingAgencies, q.PopulationCovered, q.DataRevision, q.ScoringEligible, q.ExclusionReason)
			if err != nil {
				return err
			}
		}
	}
	completed := time.Now().UTC().Format(time.RFC3339)
	status := "completed"
	summary := ""
	if rejected > 0 {
		status = "completed_with_errors"
		summary = fmt.Sprintf("%d observation(s) rejected", rejected)
	}
	_, err = tx.ExecContext(ctx, `UPDATE imports SET status=?,completed_at=?,records_read=?,records_inserted=?,records_rejected=?,checksum=?,error_summary=? WHERE id=?`, status, completed, len(batch.Observations), inserted, rejected, checksum, summary, importID)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Service) fail(importID int64, message string) error {
	done := time.Now().UTC().Format(time.RFC3339)
	_, err := s.db.Exec(`UPDATE imports SET status='failed',completed_at=?,error_summary=? WHERE id=?`, done, message, importID)
	return err
}
func (s *Service) addErrorTx(tx *sql.Tx, importID int64, row int, field, raw, message string) {
	_, _ = tx.Exec(`INSERT INTO import_errors(import_id,row_number,field_name,raw_value,error_message) VALUES(?,?,?,?,?)`, importID, row, field, raw, message)
}

type censusAdapter struct {
	client                           *http.Client
	key, id, source, metric, dataset string
}

func (a *censusAdapter) SourceName() string { return a.source }
func (a *censusAdapter) Spec() Spec {
	s := Spec{a.id, map[string]string{"census-college-enrollment": "ACS college enrollment", "census-renter-burden": "ACS renter cost burden"}[a.id], "U.S. Census Bureau", []string{a.metric}, 2024, a.key != "", ""}
	if a.key == "" {
		s.UnavailableReason = "CENSUS_API_KEY is not configured"
	}
	return s
}
func (a *censusAdapter) Fetch(ctx context.Context, year int) (Batch, error) {
	vars := "NAME,S1401_C01_029E,S1401_C01_030E"
	if a.dataset == "profile" {
		vars = "NAME,DP04_0141PE,DP04_0142PE"
	}
	u := fmt.Sprintf("https://api.census.gov/data/%d/acs/acs1/%s?get=%s&for=state:*&key=%s", year, a.dataset, url.QueryEscape(vars), url.QueryEscape(a.key))
	var rows [][]string
	if err := getJSON(ctx, a.client, u, &rows); err != nil {
		return Batch{}, err
	}
	var out []Observation
	for _, row := range rows[1:] {
		if len(row) < 4 || row[3] == "11" || row[3] == "72" {
			continue
		}
		n, _ := strconv.ParseFloat(row[1], 64)
		d, _ := strconv.ParseFloat(row[2], 64)
		value := d
		if a.dataset == "subject" {
			if n == 0 {
				continue
			}
			value = d / n * 100
		}
		out = append(out, Observation{stateCodeByFIPS[row[3]], a.metric, row[0], year, math.Round(value*10000) / 10000, nil})
	}
	return Batch{out}, nil
}

type cdcAdapter struct{ client *http.Client }

func (a *cdcAdapter) SourceName() string {
	return "CDC BRFSS Nutrition, Physical Activity and Obesity 2024"
}
func (a *cdcAdapter) Spec() Spec {
	return Spec{"cdc-adult-obesity", "BRFSS adult obesity", "Centers for Disease Control and Prevention", []string{"adult-obesity-prevalence"}, 2024, true, ""}
}
func (a *cdcAdapter) Fetch(ctx context.Context, year int) (Batch, error) {
	q := url.Values{}
	q.Set("$select", "locationabbr,locationdesc,data_value")
	q.Set("$where", fmt.Sprintf("yearstart='%d' and questionid='Q036' and stratificationcategory1='Total'", year))
	q.Set("$limit", "100")
	var rows []struct {
		Code  string `json:"locationabbr"`
		Name  string `json:"locationdesc"`
		Value string `json:"data_value"`
	}
	if err := getJSON(ctx, a.client, "https://data.cdc.gov/resource/hn4x-zwk7.json?"+q.Encode(), &rows); err != nil {
		return Batch{}, err
	}
	out := []Observation{}
	for _, r := range rows {
		v, e := strconv.ParseFloat(r.Value, 64)
		if e == nil && stateFIPSByCode[r.Code] != "" {
			out = append(out, Observation{r.Code, "adult-obesity-prevalence", r.Name, year, v, nil})
		}
	}
	return Batch{out}, nil
}

type blsAdapter struct{ client *http.Client }

func (a *blsAdapter) SourceName() string { return "CES State and Metro Area 2024" }
func (a *blsAdapter) Spec() Spec {
	return Spec{"bls-employment-growth", "BLS annual employment growth", "U.S. Bureau of Labor Statistics", []string{"annual-employment-growth"}, 2024, true, ""}
}
func (a *blsAdapter) Fetch(ctx context.Context, year int) (Batch, error) {
	ids := []string{}
	byID := map[string]string{}
	for code, fips := range stateFIPSByCode {
		id := "SMS" + fips + "000000000000001"
		ids = append(ids, id)
		byID[id] = code
	}
	sort.Strings(ids)
	out := []Observation{}
	for start := 0; start < len(ids); start += 25 {
		end := start + 25
		if end > len(ids) {
			end = len(ids)
		}
		payload, _ := json.Marshal(map[string]any{"seriesid": ids[start:end], "startyear": strconv.Itoa(year - 1), "endyear": strconv.Itoa(year)})
		req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.bls.gov/publicAPI/v2/timeseries/data/", bytes.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		resp, err := a.client.Do(req)
		if err != nil {
			return Batch{}, err
		}
		var body struct {
			Results struct {
				Series []struct {
					ID   string `json:"seriesID"`
					Data []struct{ Year, Period, Value string }
				}
			}
		}
		err = json.NewDecoder(io.LimitReader(resp.Body, 4<<20)).Decode(&body)
		resp.Body.Close()
		if err != nil {
			return Batch{}, err
		}
		for _, series := range body.Results.Series {
			sum := map[int]float64{}
			count := map[int]int{}
			for _, d := range series.Data {
				y, _ := strconv.Atoi(d.Year)
				if len(d.Period) == 3 && d.Period >= "M01" && d.Period <= "M12" {
					v, _ := strconv.ParseFloat(d.Value, 64)
					sum[y] += v
					count[y]++
				}
			}
			if count[year] == 12 && count[year-1] == 12 {
				current := sum[year] / 12
				previous := sum[year-1] / 12
				out = append(out, Observation{byID[series.ID], "annual-employment-growth", series.ID, year, math.Round(((current/previous)-1)*1000000) / 10000, nil})
			}
		}
	}
	return Batch{out}, nil
}

type fbiAdapter struct{ client *http.Client }

func (a *fbiAdapter) SourceName() string { return "FBI Crime Data Explorer 2024 Property Crime" }
func (a *fbiAdapter) Spec() Spec {
	return Spec{"fbi-property-crime", "FBI property crime", "Federal Bureau of Investigation", []string{"property-crime-rate"}, 2024, true, ""}
}
func (a *fbiAdapter) Fetch(ctx context.Context, year int) (Batch, error) {
	codes := make([]string, 0, len(stateFIPSByCode))
	for c := range stateFIPSByCode {
		codes = append(codes, c)
	}
	sort.Strings(codes)
	out := []Observation{}
	for _, code := range codes {
		var body struct {
			Offenses struct {
				Actuals map[string]map[string]float64 `json:"actuals"`
			} `json:"offenses"`
			Populations struct {
				Population   map[string]map[string]float64 `json:"population"`
				Participated map[string]map[string]float64 `json:"participated_population"`
			} `json:"populations"`
			Tooltips   map[string]map[string]map[string]float64 `json:"tooltips"`
			Properties struct {
				Refresh map[string]string `json:"last_refresh_date"`
			} `json:"cde_properties"`
		}
		u := fmt.Sprintf("https://cde.ucr.cjis.gov/LATEST/summarized/state/%s/property-crime?from=01-%d&to=12-%d", code, year, year)
		if err := getJSON(ctx, a.client, u, &body); err != nil {
			return Batch{}, err
		}
		name := ""
		for k := range body.Populations.Population {
			if k != "United States" {
				name = k
				break
			}
		}
		actual := body.Offenses.Actuals[name+" Offenses"]
		populationMap := body.Populations.Population[name]
		population := firstFloat(populationMap)
		sum := 0.0
		for _, v := range actual {
			sum += v
		}
		coverageMap := body.Tooltips["Percent of Population Coverage"][name]
		coverage := 100.0
		for _, v := range coverageMap {
			if v < coverage {
				coverage = v
			}
		}
		covered := int64(math.Round(averageMap(body.Populations.Participated[name])))
		eligible := coverage >= 90
		reason := ""
		if !eligible {
			reason = "Minimum monthly FBI population coverage below 90%"
		}
		q := &Quality{ReportingCoverage: &coverage, PopulationCovered: &covered, DataRevision: body.Properties.Refresh["UCR"], ScoringEligible: eligible, ExclusionReason: reason}
		out = append(out, Observation{code, "property-crime-rate", code, year, math.Round(sum/population*10000000) / 100, q})
	}
	return Batch{out}, nil
}

func getJSON(ctx context.Context, client *http.Client, address string, dst any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, address, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("source returned %s", resp.Status)
	}
	return json.NewDecoder(io.LimitReader(resp.Body, 8<<20)).Decode(dst)
}
func firstFloat(m map[string]float64) float64 {
	for _, v := range m {
		return v
	}
	return 0
}
func averageMap(m map[string]float64) float64 {
	if len(m) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range m {
		sum += v
	}
	return sum / float64(len(m))
}
func censusKey() string {
	if key := os.Getenv("CENSUS_API_KEY"); key != "" {
		return key
	}
	file, err := os.Open(".env")
	if err != nil {
		return ""
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "CENSUS_API_KEY=") {
			return strings.Trim(strings.TrimPrefix(line, "CENSUS_API_KEY="), "\"'")
		}
	}
	return ""
}

var stateFIPSByCode = map[string]string{"AL": "01", "AK": "02", "AZ": "04", "AR": "05", "CA": "06", "CO": "08", "CT": "09", "DE": "10", "FL": "12", "GA": "13", "HI": "15", "ID": "16", "IL": "17", "IN": "18", "IA": "19", "KS": "20", "KY": "21", "LA": "22", "ME": "23", "MD": "24", "MA": "25", "MI": "26", "MN": "27", "MS": "28", "MO": "29", "MT": "30", "NE": "31", "NV": "32", "NH": "33", "NJ": "34", "NM": "35", "NY": "36", "NC": "37", "ND": "38", "OH": "39", "OK": "40", "OR": "41", "PA": "42", "RI": "44", "SC": "45", "SD": "46", "TN": "47", "TX": "48", "UT": "49", "VT": "50", "VA": "51", "WA": "53", "WV": "54", "WI": "55", "WY": "56"}
var stateCodeByFIPS = func() map[string]string {
	m := map[string]string{}
	for code, fips := range stateFIPSByCode {
		m[fips] = code
	}
	return m
}()
