package routes

import (
	"net/http"
)

type Route struct {
	URIPrefix    string
	URI          string
	Method       string
	Action       func(http.ResponseWriter, *http.Request)
	AuthRequired bool
}
