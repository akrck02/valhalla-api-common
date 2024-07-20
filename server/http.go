package server

import (
	"bytes"
	"io"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

func MultipartToBytes(c *gin.Context, key string) ([]byte, error) {

	fileheader, err := c.FormFile(key)
	if nil != err {
		return nil, err
	}

	// if empty file
	if fileheader == nil {
		return nil, nil
	}

	// open file
	var file multipart.File
	file, err = fileheader.Open()

	if nil != err {
		return nil, err
	}

	// read file bytes
	defer file.Close()
	if nil != err {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); nil != err {
		return nil, err
	}

	return buf.Bytes(), nil
}
