package routes

import (
	"blogger/src/controllers"
	"net/http"
)

var PostRoutes = []Route{
	{
		URIPrefix:    "/api/v1",
		URI:          "/posts",
		Method:       http.MethodGet,
		Action:       controllers.IndexPost,
		AuthRequired: false,
	},
	{
		URIPrefix:    "/api/v1",
		URI:          "/posts",
		Method:       http.MethodPost,
		Action:       controllers.StorePost,
		AuthRequired: true,
	},
	{
		URIPrefix:    "/api/v1",
		URI:          "/posts/{id}",
		Method:       http.MethodGet,
		Action:       controllers.ShowPost,
		AuthRequired: false,
	},
	{
		URIPrefix:    "/api/v1",
		URI:          "/posts/{id}",
		Method:       http.MethodPut,
		Action:       controllers.UpdatePost,
		AuthRequired: true,
	},
	{
		URIPrefix:    "/api/v1",
		URI:          "/posts/{id}",
		Method:       http.MethodPatch,
		Action:       controllers.UpdatePost,
		AuthRequired: true,
	},
	{
		URIPrefix:    "/api/v1",
		URI:          "/posts/{id}",
		Method:       http.MethodDelete,
		Action:       controllers.DeletePost,
		AuthRequired: true,
	},
	//
	{
		URIPrefix:    "/api/v1",
		URI:          "/posts/user/{id}",
		Method:       http.MethodGet,
		Action:       controllers.ShowUserPosts,
		AuthRequired: false,
	},
	{
		URIPrefix:    "/api/v1",
		URI:          "/posts/{id}/like",
		Method:       http.MethodPost,
		Action:       controllers.LikePost,
		AuthRequired: true,
	},
	{
		URIPrefix:    "/api/v1",
		URI:          "/posts/{id}/unlike",
		Method:       http.MethodPost,
		Action:       controllers.UnlikePost,
		AuthRequired: true,
	},
}
