package importer

import (
	"context"
	"database/sql"
	_ "modernc.org/sqlite"
	"testing"
)

func testDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	for _, q := range []string{
		`CREATE TABLE states(id INTEGER PRIMARY KEY,code TEXT)`, `CREATE TABLE metrics(id INTEGER PRIMARY KEY,slug TEXT,active INTEGER)`, `CREATE TABLE metric_values(id INTEGER PRIMARY KEY,state_id INTEGER,metric_id INTEGER,year INTEGER,value REAL,source_record_id TEXT,import_id INTEGER,UNIQUE(state_id,metric_id,year,import_id))`,
		`INSERT INTO states VALUES(1,'UT')`, `INSERT INTO metrics VALUES(2,'unemployment-rate',1)`} {
		if _, err = db.Exec(q); err != nil {
			t.Fatal(err)
		}
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestCSVValidatesAndImportsAcceptedRows(t *testing.T) {
	db := testDB(t)
	csv := []byte("state_code,metric_slug,year,value\nUT,unemployment-rate,2025,3.4\nXX,unemployment-rate,2025,nope\n")
	result, err := CSV(context.Background(), db, 7, csv)
	if err != nil {
		t.Fatal(err)
	}
	if result.RecordsRead != 2 || result.RecordsInserted != 1 || result.RecordsRejected != 1 {
		t.Fatalf("unexpected result: %+v", result)
	}
	var count int
	if err = db.QueryRow(`SELECT count(*) FROM metric_values WHERE import_id=7`).Scan(&count); err != nil || count != 1 {
		t.Fatalf("count=%d err=%v", count, err)
	}
}

func TestCSVWithNoValidRowsLeavesValuesUntouched(t *testing.T) {
	db := testDB(t)
	csv := []byte("state_code,metric_slug,year,value\nXX,unknown,1800,bad\n")
	result, err := CSV(context.Background(), db, 8, csv)
	if err != nil {
		t.Fatal(err)
	}
	if result.RecordsInserted != 0 || result.RecordsRejected != 1 {
		t.Fatalf("unexpected result: %+v", result)
	}
	var count int
	_ = db.QueryRow(`SELECT count(*) FROM metric_values`).Scan(&count)
	if count != 0 {
		t.Fatalf("expected no rows, got %d", count)
	}
}

func TestCSVRejectsMissingRequiredColumn(t *testing.T) {
	db := testDB(t)
	_, err := CSV(context.Background(), db, 9, []byte("state_code,year,value\nUT,2025,3\n"))
	if err == nil {
		t.Fatal("expected missing-column error")
	}
}
