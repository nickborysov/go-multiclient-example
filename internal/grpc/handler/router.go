package handler

import (
	"context"

	"github.com/nickborysov/go-multiclient-example/internal/dependency"
	pb "github.com/nickborysov/go-multiclient-example/internal/grpc/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ pb.ExampleServer = (*Router)(nil)

type Router struct {
	pb.UnimplementedExampleServer
	service dependency.Service
}

func NewRouter(service dependency.Service) *Router {
	return &Router{
		service: service,
	}
}

func (s *Router) GetExample(ctx context.Context, req *emptypb.Empty) (*pb.ExampleResponse, error) {
	resp, err := s.service.GetTestResponse()
	if err != nil {
		return nil, err
	}

	pbResp := &pb.ExampleResponse{
		Message: resp.Message,
		Success: resp.Success,
	}

	return pbResp, nil
}
