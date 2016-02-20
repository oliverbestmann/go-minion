package rest

import (
  "net/http"
  "github.com/gorilla/mux"
)

func handlePingRequest(req *http.Request, vars map[string]string) interface{} {
  return RawResponse{
    http.StatusOK,
    "application/json",
    []byte(`{"pong":true}`),
  }
}

// AddPingRoute adds an optimized ping route to the given router
func AddPingRoute(router *Router) *mux.Route {
  return router.AddRestHandler("/ping", HandlerFromFunc(func(req *http.Request, vars map[string]string) interface{} {
    return RawResponse{
      http.StatusOK,
      "application/json",
      []byte(`{"pong":true}`),
    }
  })).Methods("GET")
}

// AddStaticResourcesRoute uses a FileServer to server static resources
// from the local directory path under the url prefix.
func AddStaticResourcesRoute(router *Router, urlPrefix string, directory string) *mux.Route {
  return router.Handle(urlPrefix, http.FileServer(http.Dir(directory)))
}
