package routes

import (
	"blogger/src/controllers"
	"net/http"
)

var Authentication = []Route{
	{
		URIPrefix:    "/api/v1",
		URI:          "/authentication/login",
		Method:       http.MethodPost,
		Action:       controllers.Login,
		AuthRequired: false,
	},
	{
		URIPrefix:    "/api/v1",
		URI:          "/authentication/change-password",
		Method:       http.MethodPost,
		Action:       controllers.ChangePassword,
		AuthRequired: true,
	},
}
