package middleware

import (
	"encoding/json"
	"net/http"
	"time"

	apierror "github.com/akrck02/valhalla-core-sdk/error"
	apimodels "github.com/akrck02/valhalla-core-sdk/models/api"
)

const CONTENT_TYPE_HEADER = "Content-Type"

type EmptyResponse struct {
}

func Response(context *apimodels.ApiContext, endpoint *apimodels.Endpoint, writer http.ResponseWriter) {

	switch endpoint.ResponseMimeType {
	case apimodels.MimeApplicationJson:
		sendJson(context, writer)
	case apimodels.MimeApplicationOctetStream:
		sendOctetStream(context, writer)
	default:
		SendResponse(writer, http.StatusUnsupportedMediaType, &apimodels.Error{
			Error:   apierror.UnsupportedOperation,
			Status:  http.StatusUnsupportedMediaType,
			Message: "Unsupported response mime type",
		}, apimodels.MimeApplicationJson)
	}
}

func sendJson(context *apimodels.ApiContext, writer http.ResponseWriter) {

	// calculate the time of the request
	start := time.Now()

	// execute the function
	response, responseError := context.Trazability.Endpoint.Listener(context)

	// calculate the time of the response
	end := time.Now()
	elapsed := end.Sub(start)

	// if something went wrong, return error
	if nil != responseError {
		SendResponse(writer, responseError.Status, responseError, apimodels.MimeApplicationJson)
		return
	}

	// if response is nil, return {}
	if nil == response {
		response = &apimodels.Response{
			Code:     http.StatusNoContent,
			Response: EmptyResponse{},
		}
	}

	// send response
	response.ResponseTime = elapsed.Nanoseconds()
	context.Response = *response
	SendResponse(writer, response.Code, response.Response, apimodels.MimeApplicationJson)
}

func sendOctetStream(context *apimodels.ApiContext, writer http.ResponseWriter) {

	// execute the function
	result, responseError := context.Trazability.Endpoint.Listener(context)

	// if something went wrong, return error
	if nil != responseError {
		SendResponse(writer, responseError.Status, responseError, apimodels.MimeApplicationJson)
		return
	}

	// if response is nil, return nothing
	if nil == result {
		SendResponse(writer, http.StatusNoContent, nil, apimodels.MimeApplicationOctetStream)
		return
	}

	// send response
	context.Response = *result
	SendResponse(writer, result.Code, result.Response, apimodels.MimeApplicationOctetStream)

}

func SendResponse(w http.ResponseWriter, status int, response interface{}, contentType apimodels.MimeType) {
	w.Header().Set(CONTENT_TYPE_HEADER, string(contentType))
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
