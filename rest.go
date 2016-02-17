package go_minion

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

type RestHandler interface {
  Handle(req *http.Request, vars map[string]string) interface{}
}

type RestHandlerFunc func(req *http.Request, vars map[string]string) interface{}

type funcRestHandler struct {
  handle RestHandlerFunc
}

func (h funcRestHandler) Handle(req *http.Request, vars map[string]string) interface{} {
  return h.handle(req, vars)
}

func (r *Router) AddRestHandler(path string, handler RestHandler) *mux.Route {
  return r.Handle(path, httpBasicHandler{handler})
}

func (r *Router) AddRestHandlerFunc(path string, handler RestHandlerFunc) *mux.Route {
  return r.AddRestHandler(path, funcRestHandler{handler})
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
