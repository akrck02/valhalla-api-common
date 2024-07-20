package middleware

import apimodels "github.com/akrck02/valhalla-core-sdk/models/api"

type Middleware func(context *apimodels.ApiContext) *apimodels.Error
