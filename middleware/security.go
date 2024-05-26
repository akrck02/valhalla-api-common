package middleware

import (
	apiModels "github.com/akrck02/valhalla-api-common/models"

	userdal "github.com/akrck02/valhalla-core-dal/services/user"
	"github.com/akrck02/valhalla-core-sdk/http"
	"github.com/akrck02/valhalla-core-sdk/log"
	"github.com/akrck02/valhalla-core-sdk/models"
	"github.com/akrck02/valhalla-core-sdk/valerror"

	"github.com/gin-gonic/gin"
)

const AUTHORITATION_HEADER = "Authorization"

// Manage security
//
// [return] gin.HandlerFunc: handler
func Security(endpoints []apiModels.Endpoint) gin.HandlerFunc {
	return func(c *gin.Context) {

		var isRegistered = false

		// Check if endpoint is registered and secured
		for _, endpoint := range endpoints {
			if endpoint.Path == c.Request.URL.Path {
				if !endpoint.Secured {
					log.FormattedInfo("Endpoint ${0} is not secured", endpoint.Path)
					return
				}

				isRegistered = true
			}
		}

		// Check if endpoint is registered
		if !isRegistered {
			log.FormattedInfo("Endpoint ${0} is not registered", c.Request.URL.Path)
			return
		}
		log.FormattedInfo("Endpoint ${0} is secured", c.Request.URL.Path)

		// Get token
		token := c.Request.Header.Get(AUTHORITATION_HEADER)

		// Check if token is valid
		if token == "" {
			c.AbortWithStatusJSON(
				http.HTTP_STATUS_FORBIDDEN,
				gin.H{"code": valerror.INVALID_TOKEN, "message": "Missing token"},
			)
			return
		}

		// Check if token is valid
		user, err := userdal.IsTokenValid(token)

		if err != nil {
			c.AbortWithStatusJSON(
				err.Status,
				err,
			)

			return
		}

		// Get request

		var request, _ = c.Get("request")
		request = request.(models.Request)

		var castedRequest = request.(models.Request)
		castedRequest.User = user

		// Set user in request
		c.Set("request", castedRequest)
	}
}
