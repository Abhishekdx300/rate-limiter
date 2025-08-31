
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /ratelimiter-service cmd/rate-limiter/main.go


FROM scratch

COPY --from=builder /ratelimiter-service /ratelimiter-service

EXPOSE 50051

CMD [ "/ratelimiter-service" ]
