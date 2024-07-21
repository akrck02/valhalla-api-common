package middleware

import (
	"net/http"
	"time"

	apimodels "github.com/akrck02/valhalla-core-sdk/models/api"
)

type EmptyResponse struct {
}

func Response(context *apimodels.ApiContext) *apimodels.Error {

	// calculate the time of the request
	start := time.Now()

	// execute the function
	result, responseError := context.Trazability.Endpoint.Listener(context)

	// calculate the time of the response
	end := time.Now()
	elapsed := end.Sub(start)

	// if something went wrong, return error
	if nil != responseError {
		return responseError
	}

	// if response is nil, return {}
	if nil == result {
		context.Response = apimodels.Response{
			Code:     http.StatusNoContent,
			Response: EmptyResponse{},
		}

		return nil
	}

	// send response
	result.ResponseTime = elapsed.Nanoseconds()
	context.Response = *result
	return nil

}
