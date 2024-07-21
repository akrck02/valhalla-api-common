package apicommon

import (
	"encoding/json"
	"net/http"

	"github.com/akrck02/valhalla-api-common/configuration"
	"github.com/akrck02/valhalla-api-common/middleware"
	"github.com/akrck02/valhalla-api-common/services"
	"github.com/akrck02/valhalla-core-dal/database"

	sdkhttp "github.com/akrck02/valhalla-core-sdk/http"
	"github.com/akrck02/valhalla-core-sdk/log"
	apimodels "github.com/akrck02/valhalla-core-sdk/models/api"
)

const API_PATH = "/api/"
const CONTENT_TYPE_HEADER = "Content-Type"

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
		endpoint.Path = API_PATH + configuration.ApiName + "/" + configuration.Version + "/" + endpoint.Path
		newEndpoints = append(newEndpoints, endpoint)
	}

	// Add core info endpoint
	newEndpoints = append(newEndpoints, apimodels.Endpoint{
		Path:     API_PATH + configuration.ApiName + "/" + configuration.Version + "/info",
		Method:   sdkhttp.HTTP_METHOD_GET,
		Listener: services.ValhallaCoreInfoHttp,
		Checks:   services.EmptyCheck,
		Secured:  false,
		Database: false,
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
		case sdkhttp.HTTP_METHOD_GET:
			endpoint.Path = "GET " + endpoint.Path
		case sdkhttp.HTTP_METHOD_POST:
			endpoint.Path = "POST " + endpoint.Path
		case sdkhttp.HTTP_METHOD_PUT:
			endpoint.Path = "PUT " + endpoint.Path
		case sdkhttp.HTTP_METHOD_DELETE:
			endpoint.Path = "DELETE " + endpoint.Path
		case sdkhttp.HTTP_METHOD_PATCH:
			endpoint.Path = "PATCH " + endpoint.Path
		}

		log.FormattedInfo("Endpoint ${0} registered.", endpoint.Path)

		http.HandleFunc(endpoint.Path, func(w http.ResponseWriter, r *http.Request) {

			// log the request
			log.Info("")
			log.FormattedInfo("${0}", endpoint.Path)

			// enable CORS
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "3600")

			// create basic api context
			context := &apimodels.ApiContext{
				Trazability: apimodels.Trazability{
					Endpoint: endpoint,
				},
			}

			// Get request data
			err := middleware.Request(r, context)
			if nil != err {
				jsonResponse(w, err.Status, err)
				return
			}

			// Apply middleware to the request
			err = applyMiddleware(context)
			if nil != err {
				jsonResponse(w, err.Status, err)
				return
			}

			// Execute the endpoint
			err = middleware.Response(context)
			if nil != err {
				jsonResponse(w, err.Status, err)
				return
			}

			// Send response
			jsonResponse(w, context.Response.Code, context.Response)

			// if a database connection was created, close it
			if nil != context.Database.Client {
				defer context.Database.Client.Disconnect(database.GetDefaultContext())
			}
		})

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

func jsonResponse(w http.ResponseWriter, status int, response interface{}) {
	w.Header().Set(CONTENT_TYPE_HEADER, "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(response)
}
