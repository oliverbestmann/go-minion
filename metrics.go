package go_minion

import (
  "net/http"
  "github.com/rcrowley/go-metrics"
)

type timerHandler struct {
  Timer   metrics.Timer
  Handler RestHandler
}

func (th *timerHandler) Handle(req *http.Request, vars map[string]string) (result interface{}) {
  th.Timer.Time(func() {
    result = th.Handler.Handle(req, vars)
  })

  return
}

func NewTimerHandler(timer metrics.Timer, handler RestHandler) RestHandler {
  return &timerHandler{timer, handler}
}

func NewNamedTimerHandler(name string, r metrics.Registry, handler RestHandler) RestHandler {
  return &timerHandler{metrics.NewRegisteredTimer(name, r), handler}
}


