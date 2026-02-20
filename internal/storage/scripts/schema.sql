-- ==========================================================
-- PHASE ENGINE - DATABASE SCHEMA
-- ==========================================================

-- Drop dependent objects first
DROP VIEW IF EXISTS view_correlation_data;
DROP TABLE IF EXISTS intelligence_snapshots CASCADE;
DROP TABLE IF EXISTS price_history CASCADE;
DROP TABLE IF EXISTS macro_stats CASCADE;
DROP TABLE IF EXISTS network_stats CASCADE;

-- ==========================================================
-- 1. Price History (Raw Market Data)
-- ==========================================================
CREATE TABLE price_history (
    id SERIAL PRIMARY KEY,
    timestamp TIMESTAMPTZ NOT NULL,
    asset VARCHAR(20) NOT NULL,
    price_usd DOUBLE PRECISION NOT NULL,
    source VARCHAR(50) NOT NULL,

    CONSTRAINT unique_price_entry UNIQUE (timestamp, asset),

    -- Enforce midnight-only daily entries
    CONSTRAINT check_midnight CHECK (
        EXTRACT(HOUR FROM timestamp) = 0 AND
        EXTRACT(MINUTE FROM timestamp) = 0 AND
        EXTRACT(SECOND FROM timestamp) = 0
    )
);

-- ==========================================================
-- 2. Network Statistics (Raw Chain Data)
-- ==========================================================
CREATE TABLE network_stats (
    id SERIAL PRIMARY KEY,
    block_height BIGINT NOT NULL,
    hashrate_ths DOUBLE PRECISION NOT NULL,
    avg_block_time DOUBLE PRECISION NOT NULL,
    difficulty DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- ==========================================================
-- 3. Macro Data (Monetary Supply)
-- ==========================================================
CREATE TABLE macro_stats (
    id SERIAL PRIMARY KEY,
    m2_supply DOUBLE PRECISION NOT NULL,
    source_date TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- ==========================================================
-- 4. Intelligence Snapshots (Computed Output)
-- ==========================================================
CREATE TABLE intelligence_snapshots (
    id SERIAL PRIMARY KEY,

    snapshot_date DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    price_usd DOUBLE PRECISION NOT NULL,
    m2_supply DOUBLE PRECISION NOT NULL,
    btc_m2_ratio DOUBLE PRECISION NOT NULL,
    correlation DOUBLE PRECISION NOT NULL,

    block_height BIGINT NOT NULL,
    hashrate_ths DOUBLE PRECISION NOT NULL,
    difficulty DOUBLE PRECISION NOT NULL,
    avg_block_time DOUBLE PRECISION NOT NULL,

    network_health_score INTEGER NOT NULL,
    trend_status TEXT NOT NULL,
    source_attribution TEXT,

    CONSTRAINT unique_snapshot_date UNIQUE (snapshot_date)
);

-- ==========================================================
-- 5. Indexes
-- ==========================================================
CREATE INDEX idx_price_history_latest 
    ON price_history (asset, timestamp DESC);

CREATE INDEX idx_network_latest 
    ON network_stats (created_at DESC, block_height DESC);

CREATE INDEX idx_macro_latest 
    ON macro_stats (source_date DESC);

CREATE INDEX idx_intel_latest 
    ON intelligence_snapshots (created_at DESC);

-- ==========================================================
-- 6. Correlation View
-- ==========================================================
CREATE VIEW view_correlation_data AS
SELECT 
    p.timestamp,
    p.price_usd AS btc_price,
    m.m2_supply
FROM price_history p
JOIN macro_stats m 
    ON p.timestamp::DATE = m.source_date::DATE
WHERE p.asset = 'BTC'
ORDER BY p.timestamp DESC;