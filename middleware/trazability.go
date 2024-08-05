package middleware

import (
	apimodels "github.com/akrck02/valhalla-core-sdk/models/api"
	"github.com/akrck02/valhalla-core-sdk/utils"
)

func Trazability(context *apimodels.ApiContext) *apimodels.Error {

	time := utils.GetCurrentMillis()

	context.Trazability = apimodels.Trazability{
		Endpoint:  context.Trazability.Endpoint,
		Timestamp: &time,
		User:      context.Trazability.User,
	}

	if context.Trazability.User != nil {

		// we're assuming that the user is the launcher for now
		// this MUST be changed in the future
		context.Trazability.Launcher = apimodels.Launcher{
			Id:           context.Trazability.User.Id,
			LauncherType: apimodels.USER,
		}

	}

	return nil
}
