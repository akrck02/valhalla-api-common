package middleware

import (
	"net/http"

	"github.com/akrck02/valhalla-core-dal/database"
	apierror "github.com/akrck02/valhalla-core-sdk/error"
	apimodels "github.com/akrck02/valhalla-core-sdk/models/api"
)

func Database(context *apimodels.ApiContext) *apimodels.Error {

	if !context.Trazability.Endpoint.Database || nil != context.Database.Client {
		return nil
	}

	client := database.Connect()
	err := PingDatabase(context)
	context.Database.Client = client
	context.Database.Name = database.CurrentDatabase
	return err
}

func PingDatabase(context *apimodels.ApiContext) *apimodels.Error {

	if nil == context.Database.Client {
		return nil
	}

	err := context.Database.Client.Ping(database.GetDefaultContext(), nil)
	if nil != err {
		return &apimodels.Error{
			Error:   apierror.DatabaseError,
			Message: "Database connection error",
			Status:  http.StatusInternalServerError,
		}
	}

	return nil
}
