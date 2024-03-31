package middleware

import (
	"time"

	"github.com/akrck02/valhalla-api-common/models"
	"github.com/akrck02/valhalla-core-sdk/http"
	"github.com/akrck02/valhalla-core-sdk/log"
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

		//calculate the time of the request
		start := time.Now()
		result, error := endpoint.Listener(c)
		end := time.Now()
		elapsed := end.Sub(start)

		if error != nil {
			c.JSON(error.Status, error)
			return
		}

		if result == nil {
			log.Logger.Warn("Response is nil")
			c.JSON(http.HTTP_STATUS_OK, EmptyResponse{})
			return
		}

		result.ResponseTime = elapsed.Nanoseconds()
		c.JSON(result.Code, result)
	}

}
