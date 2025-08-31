package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/Abhishekdx300/rate-limiter/api/proto"
	"github.com/Abhishekdx300/rate-limiter/internal/limiter"
	"github.com/Abhishekdx300/rate-limiter/internal/server"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

func main() {

	const port = ":50051"

	// tcp listener
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("could not connect to Redis: %v", err)
	}

	fmt.Println("Successfully connected to Redis!")

	rateLimiter := limiter.NewRateLimiter(rdb)

	// grpc server

	s := grpc.NewServer()
	grpcServer := server.NewGrpcServer(rateLimiter, 50)
	pb.RegisterRateLimiterServiceServer(s, grpcServer)

	fmt.Printf("grpc server listening at %v\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
