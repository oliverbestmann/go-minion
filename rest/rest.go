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

type RestHandlerFunc func(*http.Request, map[string]string) interface{}

func (h RestHandlerFunc) Handle(req *http.Request, vars map[string]string) interface{} {
  return h(req, vars)
}

// ServeHTTP forwards the call to the function itself.
func (h RestHandlerFunc) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  DefaultRestHandler{h}.ServeHTTP(w, req)
}

// Adds a rest handler directly to this
func (r *Router) RestHandler(path string, handler RestHandler) *mux.Route {
  return r.Handle(path, DefaultRestHandler{handler})
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
