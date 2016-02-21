package main // import "github.com/oliverbestmann/go-minion/admin"

import (
  "log"
  "github.com/gorilla/mux"
  "github.com/rcrowley/go-metrics"
  "github.com/oliverbestmann/go-minion"
  "github.com/oliverbestmann/go-minion/rest"
  "github.com/goji/httpauth"
)

type Config struct {
  Name    string
  Metrics minion.MetricsConfig
}

func main() {
  auth := httpauth.SimpleBasicAuth("dave", "secret")

  router := mux.NewRouter()
  config := Config{}
  config.Name = "my admin test page"

  log.Println("Metrics")
  minion.SetupMetrics(metrics.DefaultRegistry, config.Metrics)

  log.Println("Admin routes")
  admin := mux.NewRouter()
  rest.AddConfigRoute(admin, config)
  rest.AddMetricsRoute(admin, metrics.DefaultRegistry)
  rest.Mount(router, "/admin", auth(admin))

  log.Println("Listening")
  rest.ListenAndServe(8080, router)
}

