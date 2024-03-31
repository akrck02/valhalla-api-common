package middleware

import (
	"github.com/akrck02/valhalla-core-sdk/http"
	"github.com/gin-gonic/gin"
)

// Manage not found endpoint requests
//
// [return] func(c *gin.Context): handler
func NotFound() func(c *gin.Context) {

	return func(c *gin.Context) {
		c.JSON(http.HTTP_STATUS_NOT_FOUND, gin.H{"code": http.HTTP_STATUS_NOT_FOUND, "message": "Not found"})
	}

}
