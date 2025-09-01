
## Problem Statement:
Design a service that enforces API rate limits across a horizontally scaled application. 

## Technologies:
- Rate Limiter Service (Golang)
- Redis
- gRPC

## LLD:
### Rate Limiting Algorithm:
The choice of rate-limiting algorithm is a critical design decision that balances accuracy, performance, and user experience.

#### Token Bucket:
This is often the preferred algorithm for user-facing APIs due to its flexibility in handling bursts of traffic. Each identifier is associated with a "bucket" that has a maximum capacity and is refilled with "tokens" at a constant rate. Each incoming request consumes one token. If the bucket is empty, the request is denied. This model allows a user who has been inactive to accumulate tokens, permitting them to make a burst of requests without being throttled, which generally leads to a better user experience.

[rate limiting algos](https://www.geeksforgeeks.org/system-design/rate-limiting-algorithms-system-design/)

### Redis Implementation:
To implement the Token Bucket algorithm in a distributed environment without race conditions, atomic operations are essential. A single, uninterruptible operation must fetch the current state, apply the logic, and write the new state back.
A Lua script executed via the EVAL command runs atomically on the Redis server.
The entire sequence, executed as a single Lua script, guarantees that concurrent requests for the same user are serialized at the Redis level, providing strong consistency.

### Go Service Implementation:

The service will use gRPC for communication with the API Gateway. gRPC is built on HTTP/2 and uses Protocol Buffers, offering lower latency and higher throughput than traditional REST/JSON, which is critical for a service on the hot path of every API request.

Within the Go service, an incoming gRPC request can be handled by a goroutine. To prevent overwhelming the service with too many concurrent Redis commands, a worker pool pattern can be implemented. A fixed number of worker goroutines can process requests from a channel, thereby controlling the level of concurrency and ensuring predictable performance.

## Implementation:

### gRPC:
```
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    api/proto/ratelimit.proto
```



## TODO:
- learn a bit of gRPC
- Redis datatypes and storage

## Test:

### without sharding
command:
```
$ ghz --insecure \
  --proto ./api/proto/ratelimit.proto \
  -c 150 -n 2000 \
  -d '{"key": "user_{{mod .RequestNumber 25}}", "limit": 5, "rate": 2}' \
--call ratelimit.RateLimiterService.ShouldAllow \
0.0.0.0:50051

```
output:
```
Summary:
  Count:        2000
  Total:        143.66 ms
  Slowest:      23.25 ms
  Fastest:      1.25 ms
  Average:      7.92 ms
  Requests/sec: 13921.64

Response time histogram:
  1.252  [2]   |
  3.452  [77]  |∎∎∎∎
  5.652  [338] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  7.852  [745] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  10.051 [377] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  12.251 [307] |∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎∎
  14.451 [92]  |∎∎∎∎∎
  16.651 [52]  |∎∎∎
  18.850 [7]   |
  21.050 [2]   |
  23.250 [1]   |

Latency distribution:
  10 % in 4.73 ms
  25 % in 5.87 ms
  50 % in 7.34 ms
  75 % in 9.86 ms
  90 % in 12.00 ms
  95 % in 13.14 ms
  99 % in 16.18 ms

Status code distribution:
  [OK]   2000 responses
```

### with sharding

command:
```
$ ghz --insecure \
  --proto ./api/proto/ratelimit.proto \
  -c 150 -n 2000 \
  -d '{"key": "user_{{mod .RequestNumber 25}}", "limit": 5, "rate": 2}' \
--call ratelimit.RateLimiterService.ShouldAllow \
0.0.0.0:50051

```
output:
```
```

