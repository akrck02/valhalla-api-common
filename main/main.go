package main

import (
	apicommon "github.com/akrck02/valhalla-api-common"
	"github.com/akrck02/valhalla-api-common/configuration"
	"github.com/akrck02/valhalla-api-common/services"
	databaseConfig "github.com/akrck02/valhalla-core-dal/configuration"
	sdkhttp "github.com/akrck02/valhalla-core-sdk/http"
	apimodels "github.com/akrck02/valhalla-core-sdk/models/api"
)

const ENV_FILE_PATH = "../.env"

func main() {

	config := configuration.LoadConfiguration(ENV_FILE_PATH)
	databaseConfig.LoadConfiguration(ENV_FILE_PATH)

	apicommon.Start(
		config,
		[]apimodels.Endpoint{
			{
				Path:     "health",
				Method:   sdkhttp.HTTP_METHOD_POST,
				Listener: services.Health,
				Checks:   services.EmptyCheck,
				Secured:  false,
				Database: true,
			},
		},
	)
}
