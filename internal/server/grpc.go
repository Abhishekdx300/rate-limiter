package server

import (
	"context"
	"log"

	pb "github.com/Abhishekdx300/rate-limiter/api/proto"
	"github.com/Abhishekdx300/rate-limiter/internal/limiter"
)

type job struct {
	request  *pb.ShouldAllowRequest
	response chan bool
}

type GrpcServer struct {
	pb.UnimplementedRateLimiterServiceServer
	limiter  *limiter.RateLimiter
	jobsChan chan job
}

func NewGrpcServer(limiter *limiter.RateLimiter, maxWorkers int) *GrpcServer {
	s := &GrpcServer{
		limiter:  limiter,
		jobsChan: make(chan job),
	}
	s.startWorkerPool(maxWorkers)
	return s
}

func (s *GrpcServer) startWorkerPool(maxWorkers int) {
	log.Printf("starting worker pool with %d workers", maxWorkers)
	for i := range maxWorkers {
		go func(workerId int) {
			for j := range s.jobsChan {
				allowed, err := s.limiter.Allow(

					context.Background(),
					j.request.GetKey(),
					int(j.request.GetLimit()),
					j.request.GetRate(),
				)
				if err != nil {
					log.Printf("Error processing job for key %s: %v", j.request.GetKey(), err)
					allowed = false
				}

				j.response <- allowed
			}
		}(i + 1)
	}
}

func (s *GrpcServer) ShouldAllow(ctx context.Context, req *pb.ShouldAllowRequest) (*pb.ShouldAllowResponse, error) {

	responseChan := make(chan bool)

	j := job{
		request:  req,
		response: responseChan,
	}

	s.jobsChan <- j

	allowed := <-responseChan

	return &pb.ShouldAllowResponse{Allowed: allowed}, nil
}
