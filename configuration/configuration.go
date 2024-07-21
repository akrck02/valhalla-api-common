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
}

func (APIConfiguration) IsDevelopment() bool {
	return os.Getenv("ENV") == "development"
}
