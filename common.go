package apicommon

import (
	"net/http"
	"os"

	"github.com/akrck02/valhalla-api-common/configuration"
	"github.com/akrck02/valhalla-api-common/middleware"
	"github.com/akrck02/valhalla-api-common/services"
	"github.com/akrck02/valhalla-core-dal/database"

	"github.com/akrck02/valhalla-core-sdk/log"
	apimodels "github.com/akrck02/valhalla-core-sdk/models/api"
)

const API_PATH = "/api/"

// ApiMiddlewares is a list of middleware functions that will be applied to all API requests
// this list can be modified to add or remove middlewares
// the order of the middlewares is important, it will be applied in the order they are listed
var ApiMiddlewares = []middleware.Middleware{
	middleware.Security,
	middleware.Database,
	middleware.Trazability,
	middleware.Checks,
}

func Start(configuration configuration.APIConfiguration, endpoints []apimodels.Endpoint) {

	// set debug or release mode
	if configuration.IsDevelopment() {
		log.Logger.WithDebug()
	}

	// show log app title and start router
	log.ShowLogAppTitle("Valhalla " + configuration.ApiName + " API")

	// register middlewares
	ApiMiddlewares = append(ApiMiddlewares, middleware.Database)

	// Add API path to endpoints
	newEndpoints := []apimodels.Endpoint{}
	for _, endpoint := range endpoints {
		endpoint.Path = API_PATH + configuration.ApiName + "/" + configuration.Version + endpoint.Path
		newEndpoints = append(newEndpoints, endpoint)
	}

	// Add core info endpoint
	newEndpoints = append(newEndpoints, apimodels.Endpoint{
		Path:             API_PATH + configuration.ApiName + "/" + configuration.Version + "/info",
		Method:           apimodels.GetMethod,
		Listener:         services.ValhallaCoreInfoHttp,
		Checks:           services.EmptyCheck,
		Secured:          false,
		Database:         false,
		ResponseMimeType: apimodels.MimeApplicationJson,
		RequestMimeType:  apimodels.MimeApplicationJson,
	})

	// Register endpoints
	registerEndpoints(newEndpoints)

	// Start listening HTTP requests
	log.FormattedInfo("API started on http://${0}:${1}${2}${3}/", configuration.Ip, configuration.Port, API_PATH, configuration.ApiName)
	log.Info("")
	state := http.ListenAndServe(configuration.Ip+":"+configuration.Port, nil)
	log.Error(state.Error())

}

func registerEndpoints(endpoints []apimodels.Endpoint) {

	for _, endpoint := range endpoints {

		switch endpoint.Method {
		case apimodels.GetMethod:
			endpoint.Path = "GET " + endpoint.Path
		case apimodels.PostMethod:
			endpoint.Path = "POST " + endpoint.Path
		case apimodels.PutMethod:
			endpoint.Path = "PUT " + endpoint.Path
		case apimodels.DeleteMethod:
			endpoint.Path = "DELETE " + endpoint.Path
		case apimodels.PatchMethod:
			endpoint.Path = "PATCH " + endpoint.Path
		}

		log.FormattedInfo("Endpoint ${0} registered.", endpoint.Path)

		// set defaults
		setEndpointDefaults(&endpoint)

		// register endpoint
		http.HandleFunc(endpoint.Path, func(writer http.ResponseWriter, request *http.Request) {

			// log the request
			log.Info("")
			log.FormattedInfo("${0}", endpoint.Path)

			// enable CORS
			writer.Header().Set("Access-Control-Allow-Origin", os.Getenv("CORS_ORIGIN"))
			writer.Header().Set("Access-Control-Allow-Methods", os.Getenv("CORS_METHODS"))
			writer.Header().Set("Access-Control-Allow-Headers", os.Getenv("CORS_HEADERS"))
			writer.Header().Set("Access-Control-Max-Age", os.Getenv("CORS_MAX_AGE"))

			// create basic api context
			context := &apimodels.ApiContext{
				Trazability: apimodels.Trazability{
					Endpoint: endpoint,
				},
			}

			// Get request data
			err := middleware.Request(request, context)
			if nil != err {
				middleware.SendResponse(writer, err.Status, err, apimodels.MimeApplicationJson)
				return
			}

			// Apply middleware to the request
			err = applyMiddleware(context)
			if nil != err {
				middleware.SendResponse(writer, err.Status, err, apimodels.MimeApplicationJson)
				return
			}

			// If a database connection was created, close it after the request
			if nil != context.Database.Client {
				defer context.Database.Client.Disconnect(database.GetDefaultContext())
			}

			// Execute the endpoint and send the response
			middleware.Response(context, writer)
		})

	}
}

func setEndpointDefaults(endpoint *apimodels.Endpoint) {

	if nil == endpoint.Checks {
		endpoint.Checks = services.EmptyCheck
	}

	if nil == endpoint.Listener {
		endpoint.Listener = services.NotImplemented
	}

	if endpoint.RequestMimeType == "" {
		endpoint.RequestMimeType = apimodels.MimeApplicationJson
	}

	if endpoint.ResponseMimeType == "" {
		endpoint.ResponseMimeType = apimodels.MimeApplicationJson
	}

}

func applyMiddleware(context *apimodels.ApiContext) *apimodels.Error {

	for _, middleware := range ApiMiddlewares {
		err := middleware(context)
		if nil != err {
			return err
		}
	}

	return nil

}
