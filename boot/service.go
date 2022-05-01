package boot

import (
	"github.com/br0tchain/docker-builder/usecase"
)

var (
	dockerService usecase.Docker
)

func LoadServices() {
	dockerService = usecase.NewDocker()
}
