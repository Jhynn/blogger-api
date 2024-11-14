package routes

import (
	"blogger/src/controllers"
	"net/http"
)

var UserRoutes = []Route{
	{
		URIPrefix:    "/api/v1",
		URI:          "/users",
		Method:       http.MethodGet,
		Action:       controllers.IndexUser,
		AuthRequired: false,
	},
	{
		URIPrefix:    "/api/v1",
		URI:          "/users",
		Method:       http.MethodPost,
		Action:       controllers.StoreUser,
		AuthRequired: false,
	},
	{
		URIPrefix:    "/api/v1",
		URI:          "/users/{id}",
		Method:       http.MethodGet,
		Action:       controllers.ShowUser,
		AuthRequired: true,
	},
	{
		URIPrefix:    "/api/v1",
		URI:          "/users/{id}",
		Method:       http.MethodPut,
		Action:       controllers.UpdateUser,
		AuthRequired: true,
	},
	{
		URIPrefix:    "/api/v1",
		URI:          "/users/{id}",
		Method:       http.MethodPatch,
		Action:       controllers.UpdateUser,
		AuthRequired: true,
	},
	{
		URIPrefix:    "/api/v1",
		URI:          "/users/{id}",
		Method:       http.MethodDelete,
		Action:       controllers.DeleteUser,
		AuthRequired: true,
	},
	{
		URIPrefix:    "/api/v1",
		URI:          "/users/me",
		Method:       http.MethodGet,
		Action:       controllers.Me,
		AuthRequired: true,
	},
	// Followers
	{
		URIPrefix:    "/api/v1",
		URI:          "/users/{id}/follow",
		Method:       http.MethodPost,
		Action:       controllers.FollowUser,
		AuthRequired: true,
	},
	{
		URIPrefix:    "/api/v1",
		URI:          "/users/{id}/unfollow",
		Method:       http.MethodPost,
		Action:       controllers.UnfollowUser,
		AuthRequired: true,
	},
	{
		URIPrefix:    "/api/v1",
		URI:          "/users/{id}/followers",
		Method:       http.MethodGet,
		Action:       controllers.GetFollowers,
		AuthRequired: false,
	},
	{
		URIPrefix:    "/api/v1",
		URI:          "/users/{id}/following",
		Method:       http.MethodGet,
		Action:       controllers.GetFollowing,
		AuthRequired: false,
	},
}
