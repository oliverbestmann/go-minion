package rest

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rcrowley/go-metrics"
	"net/http"
	"time"
)

type timerHandler struct {
	Timer   metrics.Timer
	Handler http.Handler
}

func (th *timerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	th.Timer.Time(func() {
		th.Handler.ServeHTTP(w, req)
	})
}

func NewTimerHandler(timer metrics.Timer, handler http.Handler) http.Handler {
	return &timerHandler{timer, handler}
}

func NewNamedTimerHandler(name string, r metrics.Registry, handler http.Handler) http.Handler {
	return &timerHandler{metrics.NewRegisteredTimer(name, r), handler}
}

func MeterRequests(reg metrics.Registry) func(http.Handler) http.Handler {
	return func(back http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			start := time.Now()

			writer := statusResponseWriter{w, 200}
			back.ServeHTTP(writer, req)

			route := mux.CurrentRoute(req)
			if name := route.GetName(); name != "" {
				timer := metrics.GetOrRegisterTimer(name, reg)
				timer.UpdateSince(start)

				meter := metrics.GetOrRegisterMeter(fmt.Sprintf("%s.%d", name, writer.status), reg)
				meter.Mark(1)
			}
		})
	}
}

// This is a handler that stores the status of the response.
// This is needed for the MeterRequests middleware
type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (rw statusResponseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw statusResponseWriter) Headers() http.Header {
	return rw.ResponseWriter.Header()
}

func (rw statusResponseWriter) Write(bytes []byte) (int, error) {
	return rw.ResponseWriter.Write(bytes)
}
