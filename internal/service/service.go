package service

import "github.com/nickborysov/go-multiclient-example/internal/model"

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) GetTestResponse() (model.ExampleResponse, error) {
	return model.ExampleResponse{
		Message: "Hello World!",
		Success: true,
	}, nil
}
