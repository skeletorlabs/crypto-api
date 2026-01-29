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

### API Endpoints

All endpoints are versioned under `/v1`.

**Health Check**

`GET /v1/health`

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

**Token Price**

- `symbol` is case-insensitive
- Internally normalized to lowercase
- Returned symbol is always lowercase

Example:

```
curl http://localhost:8080/v1/price/bitcoin
```

Response:

```json
{
  "symbol": "bitcoin",
  "usd": 87723,
  "cached": false
}
```

**Chains**

`GET /v1/chains`

Returns a list of supported chains with basic metrics.

Example:

```
curl http://localhost:8080/v1/chains
```

**Protocols**

`GET /v1/protocols`

Optional query parameters:

- `chain`
- `category`

Example:

```
curl "http://localhost:8080/v1/protocols?chain=Ethereum"
```

**Bitcoin Fees**

`GET /v1/bitcoin/fees`

Returns recommended Bitcoin transaction fees based on current mempool conditions.

Example:

```
curl http://localhost:8080/v1/bitcoin/fees
```

Response:

```json
{
  "fastestFee": 42,
  "halfHourFee": 28,
  "hourFee": 18,
  "cached": false
}
```

**Bitcoin Network**

`GET /v1/bitcoin/network`

Returns live Bitcoin network metrics aggregated from mempool.space.

Example:

```
curl http://localhost:8080/v1/bitcoin/network
```

Response:

```json
{
  "blockHeight": 934140,
  "hashrateTHs": 848845527.5486102,
  "difficulty": 141668107417558.2,
  "avgBlockTimeSeconds": 509.44,
  "cached": false
}
```

**Bitcoin Mempool**

`GET /v1/bitcoin/mempool`

Returns current Bitcoin mempool congestion metrics.

Example:

```
curl http://localhost:8080/v1/bitcoin/mempool
```

Response:

```json
{
  "count": 31468,
  "vsize": 15577951,
  "totalFee": 2587922,
  "cached": false
}
```

Notes:

- Data sourced from **mempool.space**
- Results are cached in-memory to reduce upstream load
- Network metrics are aggregated from multiple endpoints:
  - Block height
  - Rolling hashrate (TH/s)
  - Current difficulty
  - Average block time

### Project Structure

```
.
├── internal
│   ├── cache        # In-memory cache implementation
│   ├── filters      # Domain-level filtering logic
│   ├── httpx        # HTTP handlers
│   ├── middleware   # HTTP middleware
│   ├── models       # API response models
│   └── sources      # External data sources
├── main.go
└── README.md
```

### Running Locally

**Requirements**

```bash
Go 1.21+
```

**Run the server**

```
go run .
```

Server will start on:

```
http://localhost:8080
```

### Running Tests

Run all tests:

```
go test ./...
```

### Design Notes

- Handlers are kept small and focused
- External sources are isolated behind a `sources` layer
- Cache keys and API outputs are normalized for consistency
- API versioning is handled at the routing layer
- No frameworks, minimal dependencies

### Status

This project is actively used as an internal data layer and technical showcase.

Planned improvements:

- Docker support
- Deployment (Fly.io or similar)
- Extended metrics and derived data
- Optional ticker-to-id mapping for price endpoints

### License

MIT
