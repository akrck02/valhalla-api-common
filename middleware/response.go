package middleware

import (
	"time"

	"github.com/akrck02/valhalla-api-common/models"
	"github.com/akrck02/valhalla-core-dal/database"
	"github.com/akrck02/valhalla-core-sdk/http"
	"github.com/akrck02/valhalla-core-sdk/log"
	systemmodels "github.com/akrck02/valhalla-core-sdk/models/system"
	"github.com/gin-gonic/gin"
)

type EmptyResponse struct {
}

// Manage errors in a generic way passing the function that will be executed
//
// [param] endpoint | Endpoint: endpoint
// [return] func(c *gin.Context): handler
func APIResponseManagement(endpoint models.Endpoint) func(c *gin.Context) {

	return func(ginContext *gin.Context) {

		// calculate the time of the request
		start := time.Now()

		// get the context
		var request, _ = ginContext.Get("request")
		user := request.(systemmodels.Request).User

		valhallaContext := &systemmodels.ValhallaContext{
			Database: systemmodels.Database{},
			Launcher: systemmodels.Launcher{
				Id:           user.ID,
				LauncherType: systemmodels.USER,
			},
			Trazability: systemmodels.Trazability{
				Method:    endpoint.Path,
				Timestamp: time.Now().String(),
			},
			Request: request.(systemmodels.Request),
		}

		// check parameters and return error if necessary
		error := endpoint.Checks(valhallaContext, ginContext)
		if error != nil {
			ginContext.JSON(error.Status, error)
			return
		}

		// connect to the database if necessary
		if endpoint.Database {
			client := database.Connect()
			defer client.Disconnect(database.GetDefaultContext())

			valhallaContext.Database.Client = client
			valhallaContext.Database.Name = database.CurrentDatabase
		}

		// execute the function
		result, error := endpoint.Listener(valhallaContext)

		// calculate the time of the response
		end := time.Now()
		elapsed := end.Sub(start)

		// if something went wrong, return error
		if error != nil {
			ginContext.JSON(error.Status, error)
			return
		}

		// if response is nil, return {}
		if result == nil {
			log.Logger.Warn("Response is nil")
			ginContext.JSON(http.HTTP_STATUS_OK, EmptyResponse{})
			return
		}

		// send response
		result.ResponseTime = elapsed.Nanoseconds()
		ginContext.JSON(result.Code, result)
	}

}
