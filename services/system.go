package services

import (
	"os"
	"runtime"
	"time"

	"github.com/akrck02/valhalla-core-sdk/http"
	apimodels "github.com/akrck02/valhalla-core-sdk/models/api"
)

type ValhallaCoreInfo struct {
	Version     string   `json:"version"`
	License     string   `json:"license"`
	Maintainers []string `json:"maintainers"`
	Copyleft    string   `json:"copyleft"`
	Repository  string   `json:"repository"`
	GoVersion   string   `json:"go-version"`
}

func ValhallaCoreInfoHttp(context *apimodels.ApiContext) (*apimodels.Response, *apimodels.Error) {

	// get go version
	goVersion := runtime.Version()

	return &apimodels.Response{
		Code: http.HTTP_STATUS_OK,
		Response: ValhallaCoreInfo{
			Version: os.Getenv("VERSION"),
			License: "GNU GPLv3",
			Maintainers: []string{
				"akrck02",
				"Itros97",
			},
			Copyleft:   time.Now().Format("2006"),
			Repository: os.Getenv("REPOSITORY"),
			GoVersion:  goVersion,
		},
	}, nil
}

func EmptyCheck(context *apimodels.ApiContext) *apimodels.Error {
	return nil
}
