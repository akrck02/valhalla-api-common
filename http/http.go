package http

import "github.com/gin-gonic/gin"

func SendResponse(c *gin.Context, status int, response gin.H) {
	c.JSON(status, response)
}
