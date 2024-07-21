package middleware

import (
	"github.com/akrck02/valhalla-core-dal/database"
	apimodels "github.com/akrck02/valhalla-core-sdk/models/api"
)

func Database(context *apimodels.ApiContext) *apimodels.Error {

	if nil != context.Database.Client {
		return nil
	}

	client := database.Connect()
	context.Database.Client = client
	context.Database.Name = database.CurrentDatabase
	return nil
}
