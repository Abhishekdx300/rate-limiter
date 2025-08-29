package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Abhishekdx300/rate-limiter/internal/limiter"
	"github.com/redis/go-redis/v9"
)

func main() {

	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("could not connect to Redis: %v", err)
	}

	fmt.Println("Successfully connected to Redis!")
	fmt.Println("---")

	rateLimiter := limiter.NewRateLimiter(rdb)

	userKey := "user:123"
	limit := 5
	rate := 2.0
	fmt.Printf("Simulating 7 requests for key '%s' with a limit of %d and rate of %.1f/s...\n", userKey, limit, rate)
	for i := 0; i < 7; i++ {
		allowed, err := rateLimiter.Allow(ctx, userKey, limit, rate)
		if err != nil {
			log.Fatalf("Rate limiter failed: %v", err)
		}

		if allowed {
			fmt.Printf("Request %d: Allowed\n", i+1)
		} else {
			fmt.Printf("Request %d: Denied\n", i+1)
		}
	}

	fmt.Println("---")
	// refill
	fmt.Println("Waiting for 2 seconds...")
	time.Sleep(2 * time.Second)

	fmt.Println("trying one more req after wait")
	// it should allow now
	allowed, err := rateLimiter.Allow(ctx, userKey, limit, rate)
	if err != nil {
		log.Fatalf("Rate limiter failed: %v", err)
	}

	if allowed {
		fmt.Println("Request Allowed")
	} else {
		fmt.Println("Request Denied")
	}

}
