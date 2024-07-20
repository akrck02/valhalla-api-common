package apicommon

import (
	"encoding/json"
	"net/http"

	"github.com/akrck02/valhalla-api-common/configuration"
	"github.com/akrck02/valhalla-api-common/middleware"
	"github.com/akrck02/valhalla-core-dal/database"

	"github.com/akrck02/valhalla-api-common/services"
	sdkhhttp "github.com/akrck02/valhalla-core-sdk/http"
	"github.com/akrck02/valhalla-core-sdk/log"
	apimodels "github.com/akrck02/valhalla-core-sdk/models/api"
)

const API_PATH = "/api/"

func Start(configuration configuration.APIConfiguration, endpoints []apimodels.Endpoint) {

	// set debug or release mode
	if configuration.IsDevelopment() {
		log.Logger.WithDebug()
	}

	// show log app title and start router
	log.ShowLogAppTitle("Valhalla " + configuration.ApiName + " API")

	// // CORS configuration
	http.HandleFunc("OPTIONS", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
	})

	// Add API path to endpoints
	newEndpoints := []apimodels.Endpoint{}
	for _, endpoint := range endpoints {
		endpoint.Path = API_PATH + configuration.ApiName + "/" + configuration.Version + "/" + endpoint.Path
		newEndpoints = append(newEndpoints, endpoint)
	}

	// Add core info endpoint
	newEndpoints = append(newEndpoints, apimodels.Endpoint{
		Path:     API_PATH + configuration.ApiName + "/" + configuration.Version + "/",
		Method:   sdkhhttp.HTTP_METHOD_GET,
		Listener: services.ValhallaCoreInfoHttp,
		Checks:   services.EmptyCheck,
		Secured:  false,
		Database: false,
	})

	// Register endpoints
	registerEndpoints(newEndpoints)

	// Start listening HTTP requests
	log.FormattedInfo("API started on http://${0}:${1}${2}", configuration.Ip, configuration.Port, API_PATH)
	state := http.ListenAndServe(configuration.Ip+":"+configuration.Port, nil)
	log.Error(state.Error())

}

func registerEndpoints(endpoints []apimodels.Endpoint) {

	for _, endpoint := range endpoints {

		log.FormattedInfo("Endpoint ${0} registered.", endpoint.Path)

		switch endpoint.Method {
		case sdkhhttp.HTTP_METHOD_GET:
			endpoint.Path = "GET " + endpoint.Path
		case sdkhhttp.HTTP_METHOD_POST:
			endpoint.Path = "POST " + endpoint.Path
		case sdkhhttp.HTTP_METHOD_PUT:
			endpoint.Path = "PUT " + endpoint.Path
		case sdkhhttp.HTTP_METHOD_DELETE:
			endpoint.Path = "DELETE " + endpoint.Path
		case sdkhhttp.HTTP_METHOD_PATCH:
			endpoint.Path = "PATCH " + endpoint.Path
		}

		http.HandleFunc(endpoint.Path, func(w http.ResponseWriter, r *http.Request) {

			// create basic api context
			context := &apimodels.ApiContext{}

			// Register middleware
			err := middleware.Request(r, context)

			if nil != err {
				w.WriteHeader(err.Status)
				json.NewEncoder(w).Encode(err)
				return
			}

			err = middleware.Security(r, context)

			if nil != err {
				w.WriteHeader(err.Status)
				json.NewEncoder(w).Encode(err)
				return
			}

			err = middleware.Database(r, context)

			if nil != err {
				w.WriteHeader(err.Status)
				json.NewEncoder(w).Encode(err)
				return
			}

			err = middleware.Trazability(r, context)

			if nil != err {
				w.WriteHeader(err.Status)
				json.NewEncoder(w).Encode(err)
				return
			}

			err = middleware.Response(r, context)

			if nil != err {
				w.WriteHeader(err.Status)
				json.NewEncoder(w).Encode(err)
				return
			}

			if nil != context.Database.Client {
				defer context.Database.Client.Disconnect(database.GetDefaultContext())
			}
		})

	}
}
