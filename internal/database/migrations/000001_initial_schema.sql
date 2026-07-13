-- Initial StateScore schema (Phase 1).

CREATE TABLE IF NOT EXISTS states (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    code       TEXT NOT NULL UNIQUE,
    name       TEXT NOT NULL,
    region     TEXT,
    division   TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS categories (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    slug           TEXT NOT NULL UNIQUE,
    name           TEXT NOT NULL,
    description    TEXT,
    default_weight REAL NOT NULL DEFAULT 1.0,
    display_order  INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS data_sources (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT NOT NULL,
    publisher   TEXT,
    source_url  TEXT,
    license     TEXT,
    format      TEXT,
    description TEXT,
    created_at  TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at  TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS metrics (
    id                    INTEGER PRIMARY KEY AUTOINCREMENT,
    category_id           INTEGER NOT NULL REFERENCES categories(id),
    slug                  TEXT NOT NULL UNIQUE,
    name                  TEXT NOT NULL,
    description           TEXT,
    unit                  TEXT,
    higher_is_better      INTEGER NOT NULL DEFAULT 1,
    normalization_method  TEXT NOT NULL DEFAULT 'minmax',
    default_weight        REAL NOT NULL DEFAULT 1.0,
    source_id             INTEGER REFERENCES data_sources(id),
    active                INTEGER NOT NULL DEFAULT 1,
    created_at            TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at            TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS imports (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    source_id         INTEGER REFERENCES data_sources(id),
    status            TEXT NOT NULL DEFAULT 'pending',
    started_at        TEXT,
    completed_at      TEXT,
    records_read      INTEGER NOT NULL DEFAULT 0,
    records_inserted  INTEGER NOT NULL DEFAULT 0,
    records_rejected  INTEGER NOT NULL DEFAULT 0,
    checksum          TEXT,
    error_summary     TEXT
);

CREATE TABLE IF NOT EXISTS metric_values (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    state_id         INTEGER NOT NULL REFERENCES states(id),
    metric_id        INTEGER NOT NULL REFERENCES metrics(id),
    year             INTEGER NOT NULL,
    value            REAL NOT NULL,
    source_record_id TEXT,
    import_id        INTEGER REFERENCES imports(id),
    created_at       TEXT NOT NULL DEFAULT (datetime('now')),
    UNIQUE (state_id, metric_id, year, import_id)
);

CREATE TABLE IF NOT EXISTS import_errors (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    import_id     INTEGER NOT NULL REFERENCES imports(id) ON DELETE CASCADE,
    row_number    INTEGER,
    field_name    TEXT,
    raw_value     TEXT,
    error_message TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS scoring_profiles (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT NOT NULL UNIQUE,
    description TEXT,
    is_default  INTEGER NOT NULL DEFAULT 0,
    is_system   INTEGER NOT NULL DEFAULT 0,
    created_at  TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at  TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS profile_category_weights (
    profile_id  INTEGER NOT NULL REFERENCES scoring_profiles(id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL REFERENCES categories(id),
    weight      REAL NOT NULL,
    PRIMARY KEY (profile_id, category_id)
);

CREATE TABLE IF NOT EXISTS profile_metric_weights (
    profile_id INTEGER NOT NULL REFERENCES scoring_profiles(id) ON DELETE CASCADE,
    metric_id  INTEGER NOT NULL REFERENCES metrics(id),
    weight     REAL NOT NULL,
    PRIMARY KEY (profile_id, metric_id)
);

CREATE TABLE IF NOT EXISTS score_snapshots (
    id                   INTEGER PRIMARY KEY AUTOINCREMENT,
    profile_id           INTEGER NOT NULL REFERENCES scoring_profiles(id),
    state_id             INTEGER NOT NULL REFERENCES states(id),
    year                 INTEGER NOT NULL,
    overall_score        REAL NOT NULL,
    completeness         REAL NOT NULL DEFAULT 1.0,
    calculated_at        TEXT NOT NULL DEFAULT (datetime('now')),
    calculation_version  TEXT NOT NULL DEFAULT '1',
    UNIQUE (profile_id, state_id, year, calculation_version)
);

CREATE TABLE IF NOT EXISTS category_score_snapshots (
    score_snapshot_id INTEGER NOT NULL REFERENCES score_snapshots(id) ON DELETE CASCADE,
    category_id       INTEGER NOT NULL REFERENCES categories(id),
    score             REAL NOT NULL,
    completeness      REAL NOT NULL DEFAULT 1.0,
    PRIMARY KEY (score_snapshot_id, category_id)
);

CREATE TABLE IF NOT EXISTS application_settings (
    key        TEXT PRIMARY KEY,
    value      TEXT NOT NULL,
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);
