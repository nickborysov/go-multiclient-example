package dependency

import "github.com/nickborysov/go-multiclient-example/internal/model"

type Service interface {
	GetTestResponse() (model.ExampleResponse, error)
}
