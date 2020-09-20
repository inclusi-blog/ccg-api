package init

import (
	"ccg-api/configuration"
	"ccg-api/controller"
)

var (
	healthController = controller.HealthController{}
)

func Objects(configData *configuration.ConfigData) {
}
