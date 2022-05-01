package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/br0tchain/docker-builder/internal/logging"
	"github.com/docker/docker/api/types/network"
	"io"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	specs "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/docker/docker/api/types"
)

type DockerClient interface {
	ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error)
	ImagePull(ctx context.Context, refStr string, options types.ImagePullOptions) (io.ReadCloser, error)
	ImageBuild(ctx context.Context, buildContext io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error)
	ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *specs.Platform, containerName string) (container.ContainerCreateCreatedBody, error)
	ContainerRemove(ctx context.Context, containerID string, options types.ContainerRemoveOptions) error
	ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error
	ContainerStop(ctx context.Context, containerID string, timeout *time.Duration) error
}

type DockerService struct {
	client DockerClient
	cc     *client.Client
}

// ErrorDetail Docker Build Response
type ErrorDetail struct {
	Code    int    `json:",string"`
	Message string `json:"message"`
}

func NewDockerService() (*DockerService, error) {
	_, cli, err := dockerClientInit()
	if err != nil {
		fmt.Printf("error creating docker client: %v", err)
		return nil, errors.New("cannot create Docker client")
	} else {
		return &DockerService{client: cli, cc: cli}, nil
	}
}

// BuildImageWithContext accepts a build context path and relative Dockerfile path
func (d *DockerService) BuildImageWithContext(ctx context.Context, file io.Reader, imageTagName string, dockerfile string) (err error) {
	logger := logging.New("usecase", "BuildImageWithContext")

	buildResponse, err := d.ImageBuild(ctx, file, types.ImageBuildOptions{
		Tags:       []string{imageTagName},
		Dockerfile: dockerfile,
		Remove:     true,
	})

	if err != nil {
		logger.Print(fmt.Sprintf("unable to build docker image: \n%v\n", err))
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Print(fmt.Sprintf("error closing body reader: %v", err))
		}
	}(buildResponse.Body)

	logger.Print(buildResponse.OSType)

	return nil
}

// Helper functions

func dockerClientInit() (ctx context.Context, cli *client.Client, err error) {
	logger := logging.New("usecase", "StopContainer")
	ctx = context.Background()
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.Print(fmt.Sprintf("error creating Docker client: %v\nare you sure the client is running?\n", err))
		return nil, nil, err
	}
	return ctx, cli, nil
}

// Interface Implementation

func (d *DockerService) ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error) {
	return d.client.ContainerList(ctx, options)
}

func (d *DockerService) ImagePull(ctx context.Context, refStr string, options types.ImagePullOptions) (io.ReadCloser, error) {
	return d.client.ImagePull(ctx, refStr, options)
}

func (d *DockerService) ImageBuild(ctx context.Context, buildContext io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error) {
	return d.client.ImageBuild(ctx, buildContext, options)
}

func (d *DockerService) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *specs.Platform, containerName string) (container.ContainerCreateCreatedBody, error) {
	return d.client.ContainerCreate(ctx, config, hostConfig, networkingConfig, platform, containerName)
}

func (d *DockerService) ContainerRemove(ctx context.Context, containerID string, options types.ContainerRemoveOptions) error {
	return d.client.ContainerRemove(ctx, containerID, options)
}

func (d *DockerService) ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error {
	return d.client.ContainerStart(ctx, containerID, options)
}

func (d *DockerService) ContainerStop(ctx context.Context, containerID string, timeout *time.Duration) error {
	return d.client.ContainerStop(ctx, containerID, timeout)
}
