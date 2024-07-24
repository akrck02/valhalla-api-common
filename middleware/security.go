package middleware

import (
	"net/http"

	"github.com/akrck02/valhalla-core-dal/database"
	userdal "github.com/akrck02/valhalla-core-dal/services/user"
	apierror "github.com/akrck02/valhalla-core-sdk/error"
	"github.com/akrck02/valhalla-core-sdk/log"
	apimodels "github.com/akrck02/valhalla-core-sdk/models/api"
)

const AUTHORITATION_HEADER = "Authorization"

func Security(context *apimodels.ApiContext) *apimodels.Error {

	// Check if endpoint is registered and secured
	if !context.Trazability.Endpoint.Secured {
		log.FormattedInfo("Endpoint ${0} is not secured", context.Trazability.Endpoint.Path)
		return nil
	}

	log.FormattedInfo("Endpoint ${0} is secured", context.Trazability.Endpoint.Path)

	// Check if token is empty
	if context.Request.Authorization == "" {
		return &apimodels.Error{
			Status:  http.StatusForbidden,
			Error:   apierror.InvalidToken,
			Message: "Missing token",
		}

	}

	// Connect to the database and add it to the context
	client := database.Connect()
	context.Database.Client = client

	// Check if token is valid
	// for now we're just checking if the token is in database
	// TODO: Enchance this to check if the token is expired
	// TODO: Enchance this to check if the token claims are valid
	user, err := userdal.IsTokenValid(client, context.Request.Authorization)

	if nil != err {
		return err
	}

	log.FormattedInfo("User ${0} is authorized", user.Username)

	// add user to request
	context.Trazability.User = user
	return nil
}
