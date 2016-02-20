package minion

import "net/http"

type RestHandlerFunc func(req *http.Request, vars map[string]string) interface{}

type funcRestHandler struct {
  handle RestHandlerFunc
}

func (h funcRestHandler) Handle(req *http.Request, vars map[string]string) interface{} {
  return h.handle(req, vars)
}

// HandlerFromFunc creates a new RestHandler from the given function.
func HandlerFromFunc(f RestHandlerFunc) RestHandler {
  return &funcRestHandler{f}
}
