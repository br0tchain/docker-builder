package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"mime/multipart"
)

type DockerService struct {
	mock.Mock
}

func (m *DockerService) BuildImage(ctx context.Context, dockerFile multipart.File) (float64, error) {
	args := m.Called(ctx, dockerFile)
	if args.Get(0) != nil {
		obj, _ := args.Get(0).(float64)
		return obj, nil
	}
	return 0, args.Error(1)
}
