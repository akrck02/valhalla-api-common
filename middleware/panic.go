package middleware

import (
	"github.com/akrck02/valhalla-core-sdk/error"
	"github.com/akrck02/valhalla-core-sdk/http"
	"github.com/akrck02/valhalla-core-sdk/models"
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

			c.JSON(http.HTTP_STATUS_INTERNAL_SERVER_ERROR, &models.Error{
				Status:  http.HTTP_STATUS_INTERNAL_SERVER_ERROR,
				Error:   error.UNEXPECTED_ERROR,
				Message: err.Error(),
			})
		}
	}
}
