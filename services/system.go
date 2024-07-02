package services

import (
	"os"
	"runtime"
	"time"

	"github.com/akrck02/valhalla-core-sdk/http"
	"github.com/akrck02/valhalla-core-sdk/models"
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

func ValhallaCoreInfoHttp(context systemmodels.ValhallaContext, gin *gin.Context) (*models.Response, *models.Error) {

	// get go version
	goVersion := runtime.Version()

	return &models.Response{
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

func ValhallaCoreInfoHttpCheck(context systemmodels.ValhallaContext, gin *gin.Context) (*models.Response, *models.Error) {
	return nil, nil
}
