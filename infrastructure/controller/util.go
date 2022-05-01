package controller

import (
	"context"
	"github.com/br0tchain/docker-builder/domain/response"
	"github.com/br0tchain/docker-builder/internal/logging"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func handlingError(gtx *gin.Context, logger *logrus.Entry, err error) {
	errResp, ok := err.(response.ErrorResponse)
	if !ok {
		errResp = response.NewError(response.ErrUndefined, err.Error())
	}

	// Logging and tracing
	logger.Error(err)
	gtx.JSON(errResp.GetHTTPErrorCode(), err)
}

func extractContext(gtx *gin.Context, spanName string) (*logrus.Entry, context.Context) {
	logger := logging.New("controller", spanName)
	logger.Infof("Entering")
	defer logger.Infof("Exiting")
	return logger, gtx.Request.Context()

}
