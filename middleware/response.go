package middleware

import (
	"encoding/json"
	"time"

	"github.com/akrck02/valhalla-core-sdk/http"
	"github.com/akrck02/valhalla-core-sdk/log"
	apimodels "github.com/akrck02/valhalla-core-sdk/models/api"
)

type EmptyResponse struct {
}

func Response(r *http.Request, context *apimodels.ApiContext) *apimodels.Error {

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
		log.Logger.Warn("Response is nil")

		//FIXME set http status code
		json.NewEncoder(w).Encode(EmptyResponse{})

		return
	}

	// send response
	result.ResponseTime = elapsed.Nanoseconds()

	json.NewEncoder(w).Encode(result)

	//FIXME set http status code
	ginContext.JSON(result.Code, result)

}
