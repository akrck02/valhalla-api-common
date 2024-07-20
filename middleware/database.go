package middleware

import (
	"net/http"

	"github.com/akrck02/valhalla-core-dal/database"
	apimodels "github.com/akrck02/valhalla-core-sdk/models/api"
)

func Database(r *http.Request, context *apimodels.ApiContext) *apimodels.Error {

	if nil != context.Database.Client {
		return nil
	}

	client := database.Connect()
	context.Database.Client = client
	return nil
}
