package main

import (
	"net/http"
)

// type jsonResponse struct {
// 	Error   bool   `json:"error"`
// 	Message string `json:"message"`
// 	Data    any    `json:"data,omitempty"`
// }

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)

	// out, _ := json.MarshalIndent(payload, "", "/t")
	// w.Header().Set("Content-type", "application/json")
	// w.WriteHeader(http.StatusAccepted)
	// w.Write(out)
}
