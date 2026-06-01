CREATE TABLE IF NOT EXISTS hashrate_history (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    hps         REAL    NOT NULL,
    recorded_at INTEGER NOT NULL  -- Unix timestamp de gravação (int64)
);

CREATE INDEX IF NOT EXISTS idx_hashrate_recorded ON hashrate_history(recorded_at DESC);
