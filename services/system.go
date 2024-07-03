package services

import (
	"os"
	"runtime"
	"time"

	"github.com/akrck02/valhalla-core-sdk/http"
	systemmodels "github.com/akrck02/valhalla-core-sdk/models/system"
	"github.com/gin-gonic/gin"
)

type ValhallaCoreInfo struct {
	Version    string   `json:"version"`
	License    string   `json:"license"`
	Authors    []string `json:"authors"`
	Copyleft   string   `json:"copyleft"`
	Repository string   `json:"repository"`
	GoVersion  string   `json:"go-version"`
}

func ValhallaCoreInfoHttp(context systemmodels.ValhallaContext) (*systemmodels.Response, *systemmodels.Error) {

	// get go version
	goVersion := runtime.Version()

	return &systemmodels.Response{
		Code: http.HTTP_STATUS_OK,
		Response: ValhallaCoreInfo{
			Version: os.Getenv("VERSION"),
			License: "GNU GPLv3",
			Authors: []string{
				"akrck02",
				"AlejandroMacazaga",
				"AnderRod01",
				"Itros97",
			},
			Copyleft:   time.Now().Format("2006"),
			Repository: os.Getenv("REPOSITORY"),
			GoVersion:  goVersion,
		},
	}, nil
}

func EmptyCheck(context systemmodels.ValhallaContext, gin *gin.Context) (*systemmodels.Response, *systemmodels.Error) {
	return nil, nil
}
