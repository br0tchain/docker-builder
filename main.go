package main

import (
	"github.com/br0tchain/docker-builder/boot"
	"github.com/br0tchain/docker-builder/internal/logging"
)

func main() {
	logger := logging.New("main", "main")
	logger.Infof("initializing services")
	boot.LoadServices()
	logger.Infof("initializing controllers")
	boot.LoadControllers()
	boot.StartServer()
}
