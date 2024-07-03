package middleware

import (
	"github.com/akrck02/valhalla-core-sdk/http"
	systemmodels "github.com/akrck02/valhalla-core-sdk/models/system"
	"github.com/akrck02/valhalla-core-sdk/valerror"
	"github.com/gin-gonic/gin"
)

// Manage errors in a generic way
//
// [return] gin.HandlerFunc: handler
func Panic() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {

			var err = c.Errors[0]

			c.JSON(http.HTTP_STATUS_INTERNAL_SERVER_ERROR, &systemmodels.Error{
				Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
				Error:   valerror.UNEXPECTED_ERROR,
				Message: err.Error(),
			})
		}
	}
}
