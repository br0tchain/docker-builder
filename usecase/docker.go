package usecase

import (
	"context"
	"github.com/br0tchain/docker-builder/domain/response"
	"github.com/br0tchain/docker-builder/internal/logging"
	"github.com/gofrs/uuid"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"mime/multipart"
	"regexp"
	"strconv"
)

var (
	re = regexp.MustCompile(`"perf":[0-9.]*`)
)

type Docker interface {
	BuildImage(ctx context.Context, dockerFile multipart.File) (string, float64, error)
}

type docker struct {
	DockerClient *DockerService
}

func NewDocker() Docker {
	cli, err := NewDockerService()
	if err != nil {
		panic(err)
	}
	return &docker{DockerClient: cli}
}

func (d docker) BuildImage(ctx context.Context, dockerFile multipart.File) (string, float64, error) {
	logger := logging.New("usecase", "Build image")
	logger.Print("Entering")
	defer logger.Print("Exiting")

	result, err := parser.Parse(dockerFile)
	if err != nil {
		return "", 0, err
	}
	logger.Debugf(result.AST.Dump())

	//building an image calling the docker client
	imageID := uuid.Must(uuid.NewV4()).String()
	err = d.DockerClient.BuildImageWithContext(ctx, dockerFile, imageID, result.AST.Value)
	if err != nil {
		return "", 0, response.Wrapf(err, response.ErrBuildImage, "error building image with provided dockerfile")
	}

	//scanning the image to check vulnerabilities
	err = d.scanImage(imageID)
	if err != nil {
		return "", 0, response.Wrapf(err, response.ErrSecurity, "scanning docker image detected security errors")
	}

	//retrieving the performance
	//perf := d.retrievePerf(imageID)
	perf := d.retrievePerfLocal(imageID)

	return imageID, perf, nil
}

func (d docker) scanImage(imageID string) error {
	// use snyk
	return nil
}

func (d docker) retrievePerf(imageID string) float64 {
	// retrieve image and get file from container

	return 0
}

func (d docker) retrievePerfLocal(dockerfile string) float64 {
	// retrieve perf from file
	matches := re.FindAllStringSubmatch(dockerfile, -1)
	if matches == nil || len(matches) == 0 {
		return 0
	}
	float, err := strconv.ParseFloat(matches[0][1], 32)
	if err != nil {
		return 0
	}
	return float
}
