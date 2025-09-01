# Rate Limiter

A gRPC-based rate limiter service in Go, using Redis for state management and the Token Bucket algorithm for rate limiting.

---

## Features

- **Token Bucket Algorithm** for flexible, burst-friendly rate limiting
- **gRPC API** for high-performance communication
- **Redis** as a fast, atomic backend store
- **Docker Compose** for easy local setup

---

## Prerequisites

- [Docker](https://www.docker.com/get-started)
- [Go](https://golang.org/dl/) (for local development)
- [protoc](https://grpc.io/docs/protoc-installation/) (for regenerating gRPC code)
- [ghz](https://ghz.sh/docs/) (for load testing, optional)

---

## Getting Started

### 1. Clone the Repository

```sh
git clone <repo-url>
cd rate-limiter
```

### 2. Build and Start Services

Start the rate limiter service and Redis using Docker Compose:

```sh
docker-compose up --build
```

- The gRPC server will be available at `localhost:50051`
- Redis will be available at `localhost:6379`

### 3. gRPC API

The service exposes a gRPC endpoint for rate limiting. See `api/proto/ratelimit.proto` for the API definition.

To regenerate Go code from proto:

```sh
protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	api/proto/ratelimit.proto
```

### 4. Running the Client

You can run the sample client (if implemented) from the `cmd/client` directory:

```sh
cd cmd/client
go run main.go
```

### 5. Load Testing

Example using [ghz](https://ghz.sh/docs/):

```sh
ghz --insecure \
	--proto ./api/proto/ratelimit.proto \
	-c 150 -n 2000 \
	-d '{"key": "user_{{mod .RequestNumber 25}}", "limit": 5, "rate": 2}' \
	--call ratelimit.RateLimiterService.ShouldAllow \
	0.0.0.0:50051
```

---

## Project Structure

```
cmd/
	client/         # Example gRPC client
		main.go
	rate-limiter/   # Rate limiter service entrypoint
		main.go
api/
	proto/          # gRPC protobuf definitions
		ratelimit.proto
internal/         # Core logic, Redis, etc.
docker-compose.yml
readme.md
```

---

## References

- [Token Bucket Algorithm](https://www.geeksforgeeks.org/system-design/rate-limiting-algorithms-system-design/)
- [gRPC in Go](https://grpc.io/docs/languages/go/quickstart/)
- [Redis Lua Scripting](https://redis.io/docs/manual/programmability/eval-intro/)

---

