package server

import (
	"context"

	pb "github.com/Abhishekdx300/rate-limiter/api/proto"
	"github.com/Abhishekdx300/rate-limiter/internal/limiter"
)

type GrpcServer struct {
	pb.UnimplementedRateLimiterServiceServer
	limiter *limiter.RateLimiter
}

func NewGrpcServer(limiter *limiter.RateLimiter) *GrpcServer {
	return &GrpcServer{limiter: limiter}
}

func (s *GrpcServer) ShouldAllow(ctx context.Context, req *pb.ShouldAllowRequest) (*pb.ShouldAllowResponse, error) {
	allowed, err := s.limiter.Allow(ctx, req.GetKey(), int(req.GetLimit()), req.GetRate())
	if err != nil {
		return nil, err
	}

	return &pb.ShouldAllowResponse{Allowed: allowed}, nil
}
