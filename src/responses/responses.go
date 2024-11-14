package responses

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSON return an JSON response for the request.
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Fatalln(err)
		}
	}
}

// Error returns an error in JSON-format.
func Error(w http.ResponseWriter, statusCode int, err error) {
	JSON(w, statusCode, struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	})
}
