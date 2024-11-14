package router

import (
	"blogger/src/middlewares"
	"blogger/src/router/routes"
	"fmt"
	"net/http"
)

// Router Retunrns a router
func Router() *http.ServeMux {
	r := http.NewServeMux()

	return Config(r)
}

// Config put all the routes in the router.
func Config(r *http.ServeMux) *http.ServeMux {
	allRoutes := routes.UserRoutes
	allRoutes = append(allRoutes, routes.Miscellanous...)
	allRoutes = append(allRoutes, routes.Authentication...)
	allRoutes = append(allRoutes, routes.PostRoutes...)

	for _, route := range allRoutes {
		if route.AuthRequired {
			r.HandleFunc(
				fmt.Sprint(route.Method, " ", route.URIPrefix, route.URI),
				middlewares.Logger(
					middlewares.Authenticated(route.Action),
				),
			)
		} else {
			r.HandleFunc(
				fmt.Sprint(route.Method, " ", route.URIPrefix, route.URI),
				middlewares.Logger(route.Action),
			)
		}
	}

	return r
}
