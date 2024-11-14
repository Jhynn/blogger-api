package routes

import (
	"blogger/src/controllers"
	"net/http"
)

var Miscellanous = []Route{
	{
		URIPrefix:    "/api/v1",
		URI:          "/ping",
		Method:       http.MethodGet,
		Action:       controllers.Pong,
		AuthRequired: false,
	},
}
