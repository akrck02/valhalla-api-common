package middleware

import (
	"net/http"

	apierror "github.com/akrck02/valhalla-core-sdk/error"
	apimodels "github.com/akrck02/valhalla-core-sdk/models/api"
)

const USER_AGENT_HEADER = "User-Agent"

func Request(r *http.Request, context *apimodels.ApiContext) *apimodels.Error {

	parserError := r.ParseForm()

	if parserError != nil {
		return &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   apierror.InvalidRequest,
			Message: parserError.Error(),
		}
	}

	context.Request = apimodels.Request{
		Authorization: r.Header.Get(AUTHORITATION_HEADER),
		IP:            r.Host,
		UserAgent:     r.Header.Get(USER_AGENT_HEADER),
		Headers:       map[string]string{},
		Body:          r.Body,
		Params:        map[string]string{},
	}

	// If files are present, add them to the context
	if r.MultipartForm != nil {
		context.Request.Files = r.MultipartForm.File
	}

	// Add headers to the context
	for key, value := range r.Header {
		for _, v := range value {
			context.Request.Headers[key] = v
		}
	}

	// Add params to the context
	for key, value := range r.Form {
		for _, v := range value {
			context.Request.Params[key] = v
		}
	}

	return nil
}
