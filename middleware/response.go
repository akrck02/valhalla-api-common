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

	return func(c *gin.Context) {

		// calculate the time of the request
		start := time.Now()

		// get the context

		context := systemmodels.ValhallaContext{
			Launcher: systemmodels.Launcher{
				Id:           1,
				LauncherType: 1,
			},
			Database: systemmodels.Database{},
			Trazability: systemmodels.Trazability{
				Method:    "GET",
				Timestamp: time.Now().String(),
			},
		}

		// check parameters
		_, error := endpoint.Checks(context, c)

		// if wrong parameters, return error
		if error != nil {
			c.JSON(error.Status, error)
			return
		}

		// connect to the database
		client := database.Connect()
		defer client.Disconnect(database.GetDefaultContext())

		context.Database.Client = client
		context.Database.Name = database.CurrentDatabase

		// execute the function
		result, error := endpoint.Listener(context, c)

		// calculate the time of the response
		end := time.Now()
		elapsed := end.Sub(start)

		// if something went wrong, return error
		if error != nil {
			c.JSON(error.Status, error)
			return
		}

		// if response is nil, return {}
		if result == nil {
			log.Logger.Warn("Response is nil")
			c.JSON(http.HTTP_STATUS_OK, EmptyResponse{})
			return
		}

		// send response
		result.ResponseTime = elapsed.Nanoseconds()
		c.JSON(result.Code, result)
	}

}
