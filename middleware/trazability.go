package middleware

import (
	"net/http"
	"time"

	apimodels "github.com/akrck02/valhalla-core-sdk/models/api"
)

func Trazability(r *http.Request, context *apimodels.ApiContext) *apimodels.Error {

	if context.Request.user != nil {
		context.Launcher = apimodels.Launcher{
			Id:           "", //user.ID,
			LauncherType: apimodels.USER,
		}
	}

	context.Trazability = apimodels.Trazability{
		Endpoint:  nil,
		Timestamp: time.Now().String(),
	}

	return nil
}
