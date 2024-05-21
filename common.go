package apicommon

import (
	"time"

	"github.com/akrck02/valhalla-api-common/configuration"
	"github.com/akrck02/valhalla-api-common/middleware"
	"github.com/akrck02/valhalla-api-common/models"
	"github.com/akrck02/valhalla-api-common/services"
	"github.com/akrck02/valhalla-core-sdk/http"
	"github.com/akrck02/valhalla-core-sdk/log"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

const API_PATH = "/api/"

// Start API
func Start(configuration configuration.APIConfiguration, endpoints []models.Endpoint) {

	// set debug or release mode
	if configuration.IsDevelopment() {
		gin.SetMode(gin.DebugMode)
		log.Logger.WithDebug()
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	log.ShowLogAppTitle()
	router := gin.Default()
	router.NoRoute(middleware.NotFound())

	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE, OPTIONS",
		RequestHeaders:  "User-Agent, Accept, Accept-Language, Authorization, Accept-Encoding, Referer, Content-type, mode, Origin, Connection, Sec-Fetch-Dest, Sec-Fetch-Mode, Sec-Fetch-Site, Pragma, Cache-Control",
		ExposedHeaders:  "",
		MaxAge:          300 * time.Second,
		Credentials:     false,
		ValidateHeaders: false,
	}))

	newEndpoints := []models.Endpoint{}

	// Add API path to endpoints
	for _, endpoint := range endpoints {
		log.FormattedInfo("Endpoint ${0} registered", API_PATH+configuration.ApiName+"/"+configuration.Version+"/"+endpoint.Path)
		endpoint.Path = API_PATH + configuration.ApiName + "/" + configuration.Version + "/" + endpoint.Path
		newEndpoints = append(newEndpoints, endpoint)
	}

	newEndpoints = append(newEndpoints, models.Endpoint{
		Path:     API_PATH + configuration.ApiName + "/" + configuration.Version + "/",
		Method:   http.HTTP_METHOD_GET,
		Listener: services.ValhallaCoreInfoHttp,
	})

	// Register middleware
	router.Use(middleware.Request())
	router.Use(middleware.Security(newEndpoints))
	router.Use(middleware.Panic())

	registerEndpoints(router, newEndpoints)

	log.FormattedInfo("API started on https://${0}:${1}${2}", configuration.Ip, configuration.Port, API_PATH)
	state := router.Run(configuration.Ip + ":" + configuration.Port)
	log.Error(state.Error())

}

// Register endpoints
//
// [param] router | *gin.Engine: router
// [param] endpoints | []models.Endpoint: endpoints
func registerEndpoints(router *gin.Engine, endpoints []models.Endpoint) {

	for _, endpoint := range endpoints {

		log.FormattedInfo("Registering endpoint ${0}", endpoint.Path)

		switch endpoint.Method {
		case http.HTTP_METHOD_GET:
			router.GET(endpoint.Path, middleware.APIResponseManagement(endpoint))
		case http.HTTP_METHOD_POST:
			router.POST(endpoint.Path, middleware.APIResponseManagement(endpoint))
		case http.HTTP_METHOD_PUT:
			router.PUT(endpoint.Path, middleware.APIResponseManagement(endpoint))
		case http.HTTP_METHOD_DELETE:
			router.DELETE(endpoint.Path, middleware.APIResponseManagement(endpoint))
		}
	}
}
