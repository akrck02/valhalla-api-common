package middleware

import (
	"net/http"

	apimodels "github.com/akrck02/valhalla-core-sdk/models/api"
)

func Checks(r *http.Request, context *apimodels.ApiContext) *apimodels.Error {
	// check parameters and return checkError if necessary
	checkError := context.Trazability.Endpoint.Checks(context)
	if checkError != nil {
		return checkError
	}
}
