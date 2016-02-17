package go_minion

import (
  "net/http"
  "encoding/json"
  "github.com/gorilla/mux"
  "strconv"
)

type httpBasicHandler struct {
  RestHandler
}

func (h httpBasicHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  vars := mux.Vars(req)
  result := h.Handle(req, vars)
  writeResult(result, w)
}

// Serializes and writes some value to the ResponseWriter. Different value types
// are supported. You can serialize errors, raw byte slices or any other object
// as json.
func writeResult(value interface{}, w http.ResponseWriter) {
  switch result := value.(type) {
  case error:
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte(result.Error()))

  case Response:
    writeJsonResponse(w, result.Status, result.Content)

  case []byte:
    w.WriteHeader(http.StatusOK)
    w.Write(result)

  default:
    writeJsonResponse(w, http.StatusOK, result)
  }
}

// Tries to serialize the given value as json and writes it to the
// ResponseWriter if possible. Also sets the provided http status if
// serialization was successful.
func writeJsonResponse(w http.ResponseWriter, status int, value interface{}) {
  bytes, err := json.Marshal(value)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte(err.Error()))
    return
  }

  // everything is fine.
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("Content-Length", strconv.Itoa(len(bytes)))
  w.WriteHeader(status)
  w.Write(bytes)
}
