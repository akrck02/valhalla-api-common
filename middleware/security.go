package middleware

import (
	"net/http"

	"github.com/akrck02/valhalla-core-dal/database"
	userdal "github.com/akrck02/valhalla-core-dal/services/user"
	"github.com/akrck02/valhalla-core-sdk/log"
	apimodels "github.com/akrck02/valhalla-core-sdk/models/api"
)

const AUTHORITATION_HEADER = "Authorization"

func Security(r *http.Request, context *apimodels.ApiContext) *apimodels.Error {

	// Check if endpoint is registered and secured
	if !context.Trazability.Endpoint.Secured {
		log.FormattedInfo("Endpoint ${0} is not secured", context.Trazability.Endpoint.Path)
		return
	}

	log.FormattedInfo("Endpoint ${0} is secured", context.Trazability.Endpoint.Path)

	// Check if token is empty
	if context.Request.Authorization == "" {
		// c.AbortWithStatusJSON(
		// 	sdkhttp.HTTP_STATUS_FORBIDDEN,
		// 	gin.H{"code": valerror.INVALID_TOKEN, "message": "Missing token"},
		// )
		return
	}

	// Connect to the database if necessary and add it to the context
	client := database.Connect()
	context.Database.Client = client

	// Check if token is valid
	user, err := userdal.IsTokenValid(client, context.Request.Authorization)

	if nil != err {
		return err
	}

	// add user to request
	context.Request.User = user
	return nil
}
