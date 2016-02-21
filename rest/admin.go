package rest

import (
  "net/http"
  "github.com/gorilla/mux"
  "github.com/rcrowley/go-metrics"
)

func AddMetricsRoute(router *mux.Router, registry metrics.Registry) {
  router.HandleFunc("/metrics.txt", func(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    metrics.WriteOnce(registry, w)
  })

  router.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    metrics.WriteJSONOnce(registry, w)
  })
}

func AddConfigRoute(router *mux.Router, config interface{}) *mux.Route {
  return AddObjectRoute(router, "/config", config)
}

func AddObjectRoute(router *mux.Router, path string, value interface{}) *mux.Route {
  return router.Handle(path, RestHandlerFunc(func(*http.Request, map[string]string) interface{} {
    return value
  }))
}


