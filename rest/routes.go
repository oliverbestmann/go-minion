package rest

import (
	"github.com/gorilla/mux"
	"net/http"
)

// AddPingRoute adds an optimized ping route to the given router
func AddPingRoute(router *mux.Router) *mux.Route {
	return router.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"pong":true}`))
	}).Methods("GET")
}

// AddStaticResourcesRoute uses a FileServer to server static resources
// from the local directory path under the url prefix.
func AddStaticResourcesRoute(router *mux.Router, urlPrefix string, directory string) *mux.Route {
	return router.Handle(urlPrefix, http.FileServer(http.Dir(directory)))
}
