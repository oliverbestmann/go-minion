package rest // import "github.com/oliverbestmann/go-minion/rest"

import (
  "os"
  "log"
  "fmt"
  "net/http"
  "github.com/gorilla/mux"
  "github.com/gorilla/handlers"
  "path"
)

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
  (&restHandlerImpl{h}).ServeHTTP(w, req)
}

// ListenAndServe starts a http server on the given port and uses the
// provided router for request dispatching and handling.
func ListenAndServe(port int, router *mux.Router, transformers ...func (http.Handler) http.Handler) {
  var handler http.Handler = router
  for _, mw := range transformers {
    handler = mw(handler)
  }

  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port),
    handlers.LoggingHandler(os.Stdout,
      handlers.RecoveryHandler()(handler))))
}

func Mount(router *mux.Router, prefix string, handler http.Handler) *mux.Route {
  prefix = path.Clean(prefix)
  return router.PathPrefix(prefix + "/").Handler(http.StripPrefix(prefix, handler))
}
