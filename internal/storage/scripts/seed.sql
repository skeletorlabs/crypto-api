-- ==========================================================
-- PHASE ENGINE - DEVELOPMENT SEED
-- ==========================================================

-- 1. Full reset
TRUNCATE intelligence_snapshots, macro_stats, network_stats, price_history
RESTART IDENTITY CASCADE;

-- 2. Generate synchronized historical data
WITH generated_data AS (
    SELECT 
        (55000 + (s.i * 200) + (random() * 500)) AS generated_price,
        (21000 + (s.i * 0.5)) AS generated_m2,
        DATE_TRUNC('day', NOW() - (s.i || ' days')::interval) AS target_date,
        s.i AS day_offset
    FROM generate_series(0, 90) AS s(i)
),

-- Insert price history (INCLUDING today)
inserted_prices AS (
    INSERT INTO price_history (timestamp, asset, price_usd, source)
    SELECT target_date, 'BTC', generated_price, 'dev-seed'
    FROM generated_data
    RETURNING timestamp
),

-- Insert macro history (INCLUDING today)
inserted_macro AS (
    INSERT INTO macro_stats (m2_supply, source_date)
    SELECT generated_m2, target_date
    FROM generated_data
    RETURNING source_date
)

-- 3. Insert historical intelligence snapshots (exclude today)
INSERT INTO intelligence_snapshots (
    snapshot_date,
    created_at,
    price_usd, 
    m2_supply, 
    btc_m2_ratio, 
    correlation, 
    block_height, 
    hashrate_ths, 
    difficulty, 
    avg_block_time,
    network_health_score, 
    trend_status, 
    source_attribution
)
SELECT 
    target_date::date,
    target_date,
    generated_price,
    generated_m2,
    (generated_price / NULLIF(generated_m2, 0)),
    (0.85 + (random() * 0.05)),
    (830000 - day_offset),
    (710000000.0 - (day_offset * 10000)),
    88000000000000.0,
    600.0,
    95,
    'Stable',
    'dev-seed'
FROM generated_data
WHERE day_offset > 0;

-- 4. Insert current network state
INSERT INTO network_stats (
    block_height, 
    hashrate_ths, 
    avg_block_time, 
    difficulty, 
    created_at
)
VALUES (
    830000,
    715000000.0,
    601.0,
    88000000000000.0,
    NOW()
);