package controller

import (
	"github.com/br0tchain/docker-builder/domain"
	"github.com/br0tchain/docker-builder/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

//Docker controller interface
type Docker interface {
	Build(gtx *gin.Context)
}

type docker struct {
	DockerService usecase.Docker
}

func NewDockerController(dockerService usecase.Docker) Docker {
	return &docker{DockerService: dockerService}
}

// Build
// @Summary Create docker image
// @Description Retrieve the dockerfile provided to build a docker image based on it
// @Success 200 "Success"
// @Failure 400 {object} response.ErrorResponse "BadRequest"
// @Failure 403 {object} response.ErrorResponse "Forbidden"
// @Failure 500 {object} response.ErrorResponse "InternalServerError"
// @Router /docker/build [POST]
func (c docker) Build(gtx *gin.Context) {
	logger, ctx := extractContext(gtx, "Create")

	dockerfile, _, err := gtx.Request.FormFile("dockerfile")
	if err != nil {
		handlingError(gtx, logger, err)
		return
	}

	imageID, perf, err := c.DockerService.BuildImage(ctx, dockerfile)
	if err != nil {
		handlingError(gtx, logger, err)
		return
	}
	resp := &domain.BuildResponse{
		Performance: perf,
		JobID:       imageID,
	}
	gtx.JSON(http.StatusOK, resp)
}
