package middleware

import (
	"net/http"

	apimodels "github.com/akrck02/valhalla-core-sdk/models/api"
)

const USER_AGENT_HEADER = "User-Agent"

func Request(r *http.Request, context *apimodels.ApiContext) *apimodels.Error {

	context.Request = apimodels.Request{
		Authorization: r.Header.Get(AUTHORITATION_HEADER),
		IP:            r.Host,
		UserAgent:     r.Header.Get(USER_AGENT_HEADER),
		Body:          r.Body,
	}

	return nil
}
