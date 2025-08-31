package main

import (
	"context"
	"log"
	"time"

	pb "github.com/Abhishekdx300/rate-limiter/api/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// client stub
	c := pb.NewRateLimiterServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	//
	log.Println("Simulating 7 requests...")
	for i := 0; i < 7; i++ {
		r, err := c.ShouldAllow(ctx, &pb.ShouldAllowRequest{
			Key:   "api_key:test123",
			Limit: 5,
			Rate:  2.0,
		})
		if err != nil {
			log.Fatalf("could not check rate limit: %v", err)
		}

		log.Printf("Request %d -> Allowed: %t", i+1, r.GetAllowed())
	}

}
