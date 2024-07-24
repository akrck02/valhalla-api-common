package configuration

import (
	"os"

	"github.com/akrck02/valhalla-core-sdk/log"
	"github.com/joho/godotenv"
)

type APIConfiguration struct {
	Ip         string
	Port       string
	Version    string
	ApiName    string
	Repository string

	CorsAccessControlAllowOrigin  string
	CorsAccessControlAllowMethods string
	CorsAccessControlAllowHeaders string
	CorsAccessControlMaxAge       string
}

func LoadConfiguration(path string) APIConfiguration {

	err := godotenv.Load(path)

	if nil != err {
		log.Fatal("Error loading .env file")
	}

	var configuration = APIConfiguration{
		Ip:         os.Getenv("IP"),
		Port:       os.Getenv("PORT"),
		Version:    os.Getenv("VERSION"),
		ApiName:    os.Getenv("API_NAME"),
		Repository: os.Getenv("REPOSITORY"),

		CorsAccessControlAllowOrigin:  os.Getenv("CORS_ORIGIN"),
		CorsAccessControlAllowMethods: os.Getenv("CORS_METHODS"),
		CorsAccessControlAllowHeaders: os.Getenv("CORS_HEADERS"),
		CorsAccessControlMaxAge:       os.Getenv("CORS_MAX_AGE"),
	}

	checkCompulsoryVariables(configuration)
	return configuration
}

func checkCompulsoryVariables(Configuration APIConfiguration) {
	log.Jump()
	log.Line()
	log.Info("Configuration variables")
	log.Line()
	log.Info("IP: " + Configuration.Ip)
	log.Info("PORT: " + Configuration.Port)
	log.Info("VERSION: " + Configuration.Version)
	log.Info("API_NAME: " + Configuration.ApiName)
	log.Info("REPOSITORY: " + Configuration.Repository)

	log.Line()
	log.Info("CORS_ORIGIN: " + Configuration.CorsAccessControlAllowOrigin)
	log.Info("CORS_METHODS: " + Configuration.CorsAccessControlAllowMethods)
	log.Info("CORS_HEADERS: " + Configuration.CorsAccessControlAllowHeaders)
	log.Info("CORS_MAX_AGE: " + Configuration.CorsAccessControlMaxAge)
}

func (APIConfiguration) IsDevelopment() bool {
	return os.Getenv("ENV") == "development"
}
