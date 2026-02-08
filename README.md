# Crypto API (Go)

A lightweight crypto data API written in Go, focused on simplicity, performance, and clarity.

The API is read-only and does not require authentication.

This project aggregates crypto-related data (prices, chains, protocols) from external sources and exposes a clean, versioned HTTP API with in-memory caching and test coverage.

Built as part of **Skeletor Labs** to demonstrate backend engineering, API design, and Web3 infrastructure skills.

### Features

- 🚀 Fast HTTP API: Built with Go's standard library for maximum performance.
- 🧠 Generic Cache: Type-safe in-memory caching using Go Generics (Get[T]/Set[T]).
- 📈 Intelligence Engine: Derived metrics including BTC/M2 Valuation and Statistical Correlation.
- 🧪 Test Driven: High coverage for handlers, engines, and data sources.
- 🔢 API Versioning: Clean routing under /v1.
- 🧼 Standardized Responses: Consistent metadata envelope (meta.updatedAt, meta.cached).

### API Endpoints

All endpoints are versioned under `/v1`.

**Response Metadata**
<br />Most endpoints return a standard response envelope containing metadata.

```json
{
  "meta": {
    "updatedAt": "2026-01-31T21:14:00Z",
    "cached": false
  },
  "data": { ... }
}
```

<hr style="margin: 40px 0;" />

**Health Check**

`GET /v1/health`
<br />Intentionally excluded from the standard response format.

Example:

```
curl http://localhost:8080/v1/health
```

Response:

```json
{
  "status": "ok"
}
```

Status code:

```
200 OK
```

<hr style="margin: 40px 0;" />

**Token Price**

`GET /v1/price/<token>`

- `token` is case-insensitive, token identifier (eg: bitcoin), not the ticker symbol
- Internally normalized to lowercase
- Returned token-name is always lowercase

Example:

```
curl http://localhost:8080/v1/price/bitcoin
```

Response:

```json
{
  "meta": {
    "updatedAt": "2026-01-31T21:14:00Z",
    "cached": false
  },
  "token": "bitcoin",
  "usd": 87723
}
```

<hr style="margin: 40px 0;" />

**Chains**

`GET /v1/chains`
<br />Returns a list of supported chains with basic metrics.

Example:

```
curl http://localhost:8080/v1/chains
```

Response:

```json
{
  "meta": {
    "updatedAt": "2026-01-31T21:14:00Z",
    "cached": false
  },
  "data": [
    {
      "name": "Ethereum",
      "tvl": 123456789,
      "symbol": "ETH"
    }
  ]
}
```

<hr style="margin: 40px 0;" />

**Protocols**

`GET /v1/protocols`

Optional query parameters:

- `chain`
- `category`

Example:

```
curl "http://localhost:8080/v1/protocols?chain=Ethereum"
```

Response:

```json
{
  "meta": {
    "updatedAt": "2026-01-31T21:14:00Z",
    "cached": false
  },
  "data": [
    {
      "name": "Lido",
      "slug": "lido",
      "tvl": 123456,
      "chain": "Ethereum",
      "category": "Liquid Staking"
    }
  ]
}
```

<hr style="margin: 40px 0;" />

**Bitcoin Fees**

`GET /v1/bitcoin/fees`
<br />Returns recommended Bitcoin transaction fees based on current mempool conditions.

Example:

```
curl http://localhost:8080/v1/bitcoin/fees
```

Response:

```json
{
  "meta": {
    "updatedAt": "2026-01-31T21:14:00Z",
    "cached": false
  },
  "fastestFee": 42,
  "halfHourFee": 28,
  "hourFee": 18
}
```

<hr style="margin: 40px 0;" />

**Bitcoin Network**

`GET /v1/bitcoin/network`
<br />Returns live network metrics, including Trend analysis and Halving projections.

Example:

```
curl http://localhost:8080/v1/bitcoin/network
```

Response:

```json
{
  "meta": {
    "updatedAt": "2026-02-04T23:40:21Z",
    "cached": false
  },
  "blockHeight": 935051,
  "hashrateTHs": 920496599.65,
  "difficulty": 141668107417558.2,
  "avgBlockTimeSeconds": 719.77,
  "trend": "Stable",
  "halving": {
    "currentBlock": 935051,
    "nextHalvingBlock": 1050000,
    "blocksRemaining": 114949,
    "progressPercent": 45.26,
    "estimatedDate": "2028-09-19T14:22:21Z"
  }
}
```

`trend` indicates short-term network behavior:

- `Improving`
- `Stable`
- `Worsening`

<hr style="margin: 40px 0;" />

**Bitcoin Mempool**

`GET /v1/bitcoin/mempool`
<br />Returns current Bitcoin mempool congestion metrics.

Example:

```
curl http://localhost:8080/v1/bitcoin/mempool
```

Response:

```json
{
  "meta": {
    "updatedAt": "2026-01-31T21:14:00Z",
    "cached": false
  },
  "count": 31468,
  "vsize": 15577951,
  "totalFee": 2587922
}
```

<hr style="margin: 40px 0;" />

**Bitcoin Valuation**

`GET /v1/bitcoin/valuation`
<br />Analyzes Bitcoin price relative to global M2 money supply.

Example:

```
curl http://localhost:8080/v1/bitcoin/valuation
```

Response:

```json
{
  "meta": {
    "updatedAt": "2026-02-08T10:30:00Z",
    "cached": false
  },
  "btcPrice": 69133,
  "m2SupplyBillions": 21050.5,
  "ratio": 3.28,
  "description": "Bitcoin is trading at 3.28x relative to M2 liquidity..."
}
```

<hr style="margin: 40px 0;" />

**Bitcoin Correlation**

`GET /v1/bitcoin/correlation`
<br />Calculates the Pearson correlation coefficient between BTC and M2 liquidity over a historical window.

Example:

```
curl http://localhost:8080/v1/bitcoin/correlation
```

Response:

```json
{
  "meta": {
    "updatedAt": "2026-02-08T10:30:00Z",
    "cached": true
  },
  "coefficient": 0.87,
  "sampleCount": 730,
  "startDate": "2024-02-08T00:00:00Z",
  "endDate": "2026-02-08T00:00:00Z"
}
```

<hr style="margin: 40px 0;" />

Notes:

- Data sourced from:
  - **mempool.space** (Network)
  - **CoinGecko** (Prices)
  - **FRED - Federal Reserve Economic Data** (Macro/M2)
- Results are managed via Type-safe Generic In-memory caching.
- Network metrics are aggregated from multiple endpoints:
  - Block height
  - Rolling hashrate (TH/s)
  - Current difficulty
  - Average block time
- Intelligence Layer provides advanced analytics (Valuation & Correlation).

<hr style="margin: 40px 0;" />

### Project Structure

```
├── internal
│   ├── cache        # In-memory cache implementation
│   ├── engine       # Domain-specific logic (Halving, Trend calculations)
│   │   └── bitcoin
│   │       ├── correlation  # Pearson statistical correlation analysis
│   │       ├── valuation    # Fair value metrics (M2/BTC ratio)
│   │       ├── halving      # Scarcity and epoch projections
│   │       └── trend        # Network behavior and sentiment analysis
│   ├── filters      # Domain-level filtering logic
│   ├── httpx        # HTTP handlers
│   ├── middleware   # HTTP middleware
│   ├── models       # API response models
│   └── sources      # Upstream providers
│       ├── macro    # FRED API (M2 Money Supply)
│       ├── market   # CoinGecko API (Historical & Spot Prices)
│       └── bitcoin  # Mempool.space API (Fees, Blocks, Mempool)
```

### Running Locally

Requirements

```bash
Go 1.21+
```

Run the server

```
go run .
```

Server will start on:

```
http://localhost:8080
```

### Running with Docker

1. Ensure you have a `.env` file with your `FRED_API_KEY`.
2. Run using Docker Compose:

```bash
docker-compose up --build
```

<hr style="margin: 40px 0;" />

### Running Tests

Run all tests:

```

go test ./...

```

Run tests with **gotestsum**:

```bash
# Run tests with professional formatting (requires gotestsum)
gotestsum --format pkgname-and-test-fails
```

<hr style="margin: 40px 0;" />

### Design Notes

- Handlers are kept small and focused
- External sources are isolated behind a `sources` layer
- Cache keys and API outputs are normalized for consistency
- API versioning is handled at the routing layer
- No frameworks, minimal dependencies
- **Domain-Driven Engines**: Complex business logic (like Halving projections and Trend analysis) is isolated in `internal/engine`. This keeps handlers slim and makes the core logic easily testable without HTTP concerns.
- **Stateful Trends**: Network trends are computed using a sliding window buffer, allowing - the API to detect sentiment changes rather than just showing raw snapshots.
- Type Safety: Leverage Go Generics in the cache layer to eliminate manual type assertions and runtime panics.
- Separation of Concerns: Analytics are decoupled from HTTP concerns via the engine layer, ensuring core logic is easily testable.
- Resilience: Upstream failures are handled gracefully with cached fallbacks and standardized error codes.
- **Statistical Analysis**: The correlation engine implements the Pearson Correlation Coefficient to measure the linear relationship between BTC price and M2 supply over a 730-day rolling window.

<hr style="margin: 40px 0;" />

### Status

This project is actively used as an internal data layer and technical showcase.

**Planned improvements:**

- [x] Derived data (Halving & Trends)
- [x] Intelligence Engine (Valuation & Correlation)
- [x] Type-safe Cache implementation
- [x] Docker support
- [ ] Persistence layer for historical snapshots (próximo passo natural?)

<hr style="margin: 40px 0;" />

### License

MIT
