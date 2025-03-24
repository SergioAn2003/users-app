package handler

import (
	"encoding/json"
	"net/http"
)

type ResponseError struct {
	Message string `json:"message"`
}

func (h *Handler) sendErr(w http.ResponseWriter, code int, err error, msg string) {
	h.log.ErrorF("api error: %s, code = %d", err.Error(), code)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(ResponseError{Message: msg}); err != nil {
		h.log.ErrorF("failed to send error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type emptyJSON struct{}

func (h *Handler) sendJSON(w http.ResponseWriter, code int, data any) {
	if data == nil {
		data = emptyJSON{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.log.ErrorF("failed to send json: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
