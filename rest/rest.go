package rest // import "github.com/oliverbestmann/go-minion"

import (
  "os"
  "log"
  "fmt"
  "net/http"
  "github.com/gorilla/mux"
  "github.com/gorilla/handlers"
)

type Router struct {
  *mux.Router
}

type Response struct {
  Status  int
  Content interface{}
}

type RawResponse struct {
  Status      int
  ContentType string
  Content     []byte
}

type RestHandler interface {
  Handle(req *http.Request, vars map[string]string) interface{}
}

func (r *Router) AddRestHandler(path string, handler RestHandler) *mux.Route {
  return r.Handle(path, httpBasicHandler{handler})
}

// ListenAndServe is a convenience method to start a server with
// this router instance
func (r *Router) ListenAndServe(port int) {
  ListenAndServe(port, r.Router)
}

// ListenAndServe starts a http server on the given port and uses the
// provided router for request dispatching and handling.
func ListenAndServe(port int, router *mux.Router) {
  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port),
    handlers.LoggingHandler(os.Stdout,
      handlers.RecoveryHandler()(router))))
}
