# Valhalla api common

This library makes API creation easier creating a faster way to declare endpoints,
middleware support, env file support for api properties and creating a standard 
for each response.  

## Middleware 

Middleware can be added to the HTTP request to perform checks and tranform data.
Currently the following modules are added by default:


|| Module | Action |
|---|---|---|
| ✉️ | Request | Get HTTP request data and store it in context |
| 🔐 | Security | Check security and authoritation for HTTP requests, also fills current User |
| 💾 | Database | Create and store the database connection to the context if necessary |
| 📈 | Trazability | Store endpoint and launcher data inside context |
| ✅ | Checks | Execute the custom endpoint checks and return an error if exists |
| 📨 | Response | Create a standard response object if in json format if needed |


## Info
Every API using api-common library will open by default an endpoint called /info
that returns the current info of the APIs, as well as the license, api version and go version

```json
{
    "response": {
        "version": "v1",
        "license": "GNU GPLv3",
        "maintainers": [
            "akrck02",
            "Itros97"
        ],
        "copyleft": "2024",
        "repository": "",
        "go-version": "go1.22.3"
    },
    "response_time": 8646
}
```

## Configuration

The base configuration for the APIs is the following .env

```shell
# API common data
IP=0.0.0.0
PORT=2000
VERSION=v1
API_NAME=common
ENV=release
SECRET=527a54f6-eaf2-418a-a679-6d9efdcabb8c

# Mongo database configuration
MONGO_USER=admin
MONGO_PASSWORD=p4ssw0rd
MONGO_SERVER=172.20.0.30
MONGO_PORT=27017

# CORS configuration
CORS_ORIGIN=*
CORS_METHODS=GET,POST,PUT,PATCH,OPTIONS,DELETE
CORS_HEADERS=Content-Type,Authorization
CORS_MAX_AGE=3600
```