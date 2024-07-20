package middleware

import (
	"net/http"

	apimodels "github.com/akrck02/valhalla-core-sdk/models/api"
	"github.com/akrck02/valhalla-core-sdk/valerror"
)

const USER_AGENT_HEADER = "User-Agent"

func Request(r *http.Request, context *apimodels.ApiContext) *apimodels.Error {

	parserError := r.ParseForm()

	if parserError != nil {
		return &apimodels.Error{
			Status:  http.StatusBadRequest,
			Error:   valerror.INVALID_REQUEST,
			Message: parserError.Error(),
		}
	}

	context.Request = apimodels.Request{
		Authorization: r.Header.Get(AUTHORITATION_HEADER),
		IP:            r.Host,
		UserAgent:     r.Header.Get(USER_AGENT_HEADER),
		Headers:       r.Header,
		Body:          r.Body,
		Params:        r.Form,
	}

	// If files are present, add them to the context
	if r.MultipartForm != nil {
		context.Request.Files = r.MultipartForm.File
	}

	return nil
}
