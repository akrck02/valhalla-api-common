package http

import (
	"bytes"
	"io"
	"mime/multipart"

	systemmodels "github.com/akrck02/valhalla-core-sdk/models/system"
	"github.com/gin-gonic/gin"
)

func SendResponse(c *gin.Context, status int, response gin.H) {
	c.JSON(status, response)
}

// Read body json as the given object
//
// [param] c | interface: gin context
// [param] obj | interface{}: object to read
//
// [return] error: error
func ReadBodyJson(params interface{}, obj interface{}) error {

	return nil
}

// MultipartToBytes reads a file from a multipart request
//
// [param] c | *gin.Context: gin context
// [param] key | string: key of the file
//
// [return] []byte: file bytes
func MultipartToBytes(c *gin.Context, key string) ([]byte, error) {

	fileheader, err := c.FormFile(key)
	if err != nil {
		return nil, err
	}

	// if empty file
	if fileheader == nil {
		return nil, nil
	}

	// open file
	var file multipart.File
	file, err = fileheader.Open()

	if err != nil {
		return nil, err
	}

	// read file bytes
	defer file.Close()
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GetRequestMetadata returns the request metadata
//
// [param] c | *gin.Context: gin context
//
// [return] *models.Request: request metadata
func GetRequestMetadata(c *gin.Context) *systemmodels.Request {
	var request, exists = c.Get("request")

	if !exists {
		return nil
	}

	var casted systemmodels.Request = request.(systemmodels.Request)
	return &casted
}
