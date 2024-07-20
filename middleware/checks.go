package middleware

import (
	apimodels "github.com/akrck02/valhalla-core-sdk/models/api"
)

func Checks(context *apimodels.ApiContext) *apimodels.Error {

	checkError := context.Trazability.Endpoint.Checks(context)
	if checkError != nil {
		return checkError
	}

	return nil
}
