package boot

import (
	"context"
	"github.com/br0tchain/docker-builder/infrastructure/controller"
	"github.com/br0tchain/docker-builder/internal/logging"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	router  *gin.Engine
	server  *http.Server
	closers []io.Closer
)

func LoadControllers() {
	router = gin.New()
	router.Use(gin.Recovery())

	dockerController := controller.NewDockerController(dockerService)

	router.GET("/live", func(c *gin.Context) { Live(c.Writer, c.Request) })
	router.GET("/ready", func(c *gin.Context) { Ready(c.Writer, c.Request) })

	root := router.Group("/docker")
	root.POST("/build", dockerController.Build)

}

func Live(response http.ResponseWriter, _ *http.Request) {
	response.WriteHeader(http.StatusOK)
}

func Ready(response http.ResponseWriter, _ *http.Request) {
	response.WriteHeader(http.StatusOK)
}

func StartServer() {
	logger := logging.New("boot", "StartServer")
	begin := time.Now()

	go func() {
		if err := router.Run(":8080"); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %s", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	logger.Printf("server started in %s", time.Since(begin))
	<-quit
	logger.Printf("shutting down server ...")
	for _, closer := range closers {
		closer.Close()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("server shutdown: %s", err)
	}
	logger.Printf("server exiting after %s", time.Since(begin))
}
