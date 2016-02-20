package minion // import "github.com/oliverbestmann/go-minion"

import (
  "net/http"
  "github.com/rcrowley/go-metrics"
  "time"
  "os"
  "github.com/vistarmedia/go-datadog"
  "os/exec"
  "errors"
  "log"
  "strings"

  "./rest"
)

type timerHandler struct {
  Timer   metrics.Timer
  Handler rest.RestHandler
}

func (th *timerHandler) Handle(req *http.Request, vars map[string]string) (result interface{}) {
  th.Timer.Time(func() {
    result = th.Handler.Handle(req, vars)
  })

  return
}

func NewTimerHandler(timer metrics.Timer, handler rest.RestHandler) rest.RestHandler {
  return &timerHandler{timer, handler}
}

func NewNamedTimerHandler(name string, r metrics.Registry, handler rest.RestHandler) rest.RestHandler {
  return &timerHandler{metrics.NewRegisteredTimer(name, r), handler}
}

type MetricsConfig struct {
  SampleInterval time.Duration

  // You might want to specify the hostname of the system.
  // This will override hostname auto detection.
  Hostname       string

  Datadog        struct {
                   ApiKey string

                   // Right now, tags are not supported. This is a problem with
                   // the Datadog reporter library. Should be added shortly.
                   // FIXME support tags!
                   Tags   []string
                 }

  Console        bool
}

// Tries to get the hostname of the system by exploiting different
// sources.
func hostname() (string, error) {
  name, err := os.Hostname()
  if err == nil {
    return name, nil
  }

  if name = os.Getenv("HOSTNAME"); name != "" {
    return name, nil
  }

  output, err := exec.Command("hostname").Output()
  if err == nil {
    return strings.TrimSpace(string(output)), nil
  }

  output, err = exec.Command("uname", "-n").Output()
  if err == nil {
    return strings.TrimSpace(string(output)), nil
  }

  return "", errors.New("Could not determine hostname")
}

func SetupDefaultMetrics(r metrics.Registry, config MetricsConfig) {
  if r != nil {
    r = metrics.DefaultRegistry
  }

  metrics.RegisterRuntimeMemStats(r)
  metrics.CaptureRuntimeMemStats(r, config.SampleInterval)

  if config.Console {
    metrics.WriteJSON(r, config.SampleInterval, os.Stdout)
  }

  if config.Datadog.ApiKey != "" {
    host := config.Hostname
    if host == "" {
      var err error
      if host, err = hostname(); err != nil {
        log.Print("Could not get hostname")
      }
    }

    client := datadog.Client{host, config.Datadog.ApiKey}
    client.Reporter(r).Start(config.SampleInterval)
  }
}
