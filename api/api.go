package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/davidamey/omnitureproxy/logs"
	"github.com/gorilla/mux"
)

type apiError struct {
	Error   error
	Message string
	Code    int
}

type apiHandler func(http.ResponseWriter, *http.Request) *apiError

func (fn apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if e := fn(w, r); e != nil { // e is *appError, not os.Error.
		msg := fmt.Sprintf("%s\n%v", e.Message, e.Error)
		http.Error(w, msg, e.Code)
	}
}

var fetcher logs.Fetcher

func NewApi() http.Handler {
	fetcher = logs.NewFetcher("_logs")

	api := mux.NewRouter().PathPrefix("/api").Subrouter()
	api.Handle("/", apiHandler(handler))

	lr := api.PathPrefix("/logs").Subrouter()
	lr.Handle("/", apiHandler(getDates))
	lr.Handle("/{date:\\d{4}-\\d{2}-\\d{2}}/", apiHandler(getVisitors))
	lr.Handle("/{date:\\d{4}-\\d{2}-\\d{2}}/{vid:\\d+}", apiHandler(getLog))

	// api.PathPrefix("/visitors/").HandlerFunc(visitorHandler)
	// api.PathPrefix("/csv/").HandlerFunc(csvHandler)

	return api
}

func handler(w http.ResponseWriter, r *http.Request) *apiError {
	return &apiError{nil, "Can't display record", 500}
}

func getDates(w http.ResponseWriter, r *http.Request) *apiError {
	enc := json.NewEncoder(w)

	w.Header().Set("Content-Type", "application/json")
	if err := enc.Encode(fetcher.GetLogDates()); err != nil {
		return &apiError{Error: err, Message: "JSON error", Code: 500}
	}

	return nil
}

func getVisitors(w http.ResponseWriter, r *http.Request) *apiError {
	vars := mux.Vars(r)
	enc := json.NewEncoder(w)

	w.Header().Set("Content-Type", "application/json")
	if err := enc.Encode(fetcher.GetVisitorsForDate(vars["date"])); err != nil {
		return &apiError{Error: err, Message: "JSON error", Code: 500}
	}

	return nil
}

func getLog(w http.ResponseWriter, r *http.Request) *apiError {
	vars := mux.Vars(r)

	w.Header().Set("Content-Type", "application/json")
	w.Write(fetcher.GetLog(vars["date"], vars["vid"]))
	return nil
}
