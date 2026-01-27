# Crypto API (Go)

A lightweight crypto data API written in Go, focused on simplicity, performance, and clarity.

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

- Go 1.21+

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
