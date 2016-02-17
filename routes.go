package go_minion

import (
  "net/http"
  "github.com/gorilla/mux"
)

// AddPingRoute adds an optimized ping route to the given router
func AddPingRoute(router *Router) *mux.Route {
  return router.AddRestHandlerFunc("/ping", func(req *http.Request, vars map[string]string) interface{} {
    return RawResponse{
      http.StatusOK,
      "application/json",
      []byte(`{"pong":true}`),
    }
  }).Methods("GET")
}
