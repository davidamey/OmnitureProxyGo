package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	datePattern = `/(\d{4})/(\d{2})/(\d{2})/?$`
)

func NewApi() http.Handler {
	// rgxDate = regexp.MustCompile(datePattern)

	api := mux.NewRouter().PathPrefix("/api").Subrouter()
	api.HandleFunc("/", handler)
	// api.PathPrefix("/logs").HandlerFunc(logHandler)
	// api.PathPrefix("/visitors/").HandlerFunc(visitorHandler)
	// api.PathPrefix("/csv/").HandlerFunc(csvHandler)

	return api
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API baby")
}
