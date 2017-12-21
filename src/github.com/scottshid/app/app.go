package app

import (
    "net/http"
    "encoding/json"
)

// Abstraction for rest endpoints which will take care of taking a map and creating json from it
type RestResponse struct {
    Error   error
    Payload map[string]interface{}
    Code    int
}

type AppHandler func(http.ResponseWriter, *http.Request) *RestResponse

func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    response := fn(w, r); if response != nil {
        if response.Error != nil {
            http.Error(w, response.Error.Error(), response.Code)
        } else {
            if response.Payload != nil {
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(response.Code)
                err := json.NewEncoder(w).Encode(response.Payload); if err != nil {
                    http.Error(w, "Internal error", http.StatusInternalServerError)
                    return
                }
            } else {
                w.WriteHeader(response.Code)
            }
        }
    }
}
