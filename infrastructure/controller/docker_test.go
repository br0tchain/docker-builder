package controller

import (
	"bytes"
	"encoding/json"
	"github.com/br0tchain/docker-builder/domain"
	response2 "github.com/br0tchain/docker-builder/domain/response"
	"github.com/br0tchain/docker-builder/internal/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestDocker_Build(t *testing.T) {
	dockerService := new(mocks.DockerService)
	dockerController := NewDockerController(dockerService)
	defaultPerf := 0.99

	router := gin.Default()
	router.POST("/docker/build", dockerController.Build)

	body, boundary, err := newFile("resources/dockerfile_sample.txt")
	if err != nil {
		log.Print(err.Error())
		t.Fail()
	}
	rw := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/docker/build", body)
	r.Header.Add("Content-Type", "multipart/form-data")
	r.Form.Add("boundary", boundary)

	dockerService.On("BuildImage", mock.Anything, mock.Anything).Return(defaultPerf, nil)

	router.ServeHTTP(rw, r)
	response := rw.Result()
	require.NotNil(t, response.Body)
	buildResp := new(domain.BuildResponse)
	err = json.NewDecoder(r.Body).Decode(buildResp)
	require.Nil(t, err)
	require.Equal(t, defaultPerf, buildResp.Performance)
	_ = response.Body.Close()
	require.Equal(t, http.StatusOK, response.StatusCode)
	mock.AssertExpectationsForObjects(t, dockerService)
}

func TestDocker_BuildError(t *testing.T) {
	dockerService := new(mocks.DockerService)
	dockerController := NewDockerController(dockerService)

	router := gin.Default()
	router.POST("/docker/build", dockerController.Build)

	body, boundary, err := newFile("resources/dockerfile_sample.txt")
	if err != nil {
		log.Print(err.Error())
		t.Fail()
	}
	rw := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/docker/build", body)
	r.Header.Add("Content-Type", "multipart/form-data")
	r.Form.Add("boundary", boundary)

	dockerService.On("BuildImage", mock.Anything, mock.Anything).Return(nil, response2.NewError(response2.ErrRequestInvalidPayload, "error"))

	router.ServeHTTP(rw, r)
	response := rw.Result()
	_ = response.Body.Close()
	require.Equal(t, http.StatusBadRequest, response.StatusCode)
	mock.AssertExpectationsForObjects(t, dockerService)
}

func newFile(path string) (*bytes.Buffer, string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, "", err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, "", err
	}
	_ = file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fi.Name())
	if err != nil {
		return nil, "", err
	}
	_, _ = part.Write(fileContents)
	boundary := writer.Boundary()
	err = writer.Close()
	if err != nil {
		return nil, "", err
	}

	return body, boundary, err
}
