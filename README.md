# Crypto API (Go)

A lightweight crypto data API written in Go, focused on simplicity, performance, and clarity.

The API is read-only and does not require authentication.

This project aggregates crypto-related data (prices, chains, protocols) from external sources and exposes a clean, versioned HTTP API with in-memory caching and test coverage.

Built as part of **Skeletor Labs** to demonstrate backend engineering, API design, and Web3 infrastructure skills.

### Features

- 🚀 Fast HTTP API written in Go
- 🧠 In-memory caching with TTL
- 🔁 External data aggregation (e.g. DefiLlama-style sources)
- 🧪 Unit-tested handlers
- 🧩 Clean separation of concerns
- 🔢 API versioning (`/v1`)
- 🧼 Normalized inputs and outputs
- 🧾 Standard response metadata (`meta.updatedAt`, `meta.cached`)

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
<br />Returns live Bitcoin network metrics and derived intelligence (Trend & Halving state).

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

Notes:

- Data sourced from **mempool.space**
- Results are cached in-memory to reduce upstream load
- Network metrics are aggregated from multiple endpoints:
  - Block height
  - Rolling hashrate (TH/s)
  - Current difficulty
  - Average block time
- All responses include standardized metadata for cache state and freshness

<hr style="margin: 40px 0;" />

### Project Structure

```
├── internal
│   ├── cache        # In-memory cache implementation
│   ├── engine       # Domain-specific logic (Halving, Trend calculations)
│   │   └── bitcoin
│   ├── filters      # Domain-level filtering logic
│   ├── httpx        # HTTP handlers
│   ├── middleware   # HTTP middleware
│   ├── models       # API response models
│   └── sources      # External data sources
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
- **Stateful Trends**: Network trends are computed using a sliding window buffer, allowing the API to detect sentiment changes rather than just showing raw snapshots.

<hr style="margin: 40px 0;" />

### Status

This project is actively used as an internal data layer and technical showcase.

**Planned improvements:**

- [x] Deployment (Render)
- [x] Derived data (Bitcoin Halving & Network Trends)
- [ ] Docker support
- [ ] Optional ticker-to-id mapping for price endpoints

<hr style="margin: 40px 0;" />

### License

MIT
