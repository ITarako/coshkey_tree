package server

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

func writeError(w http.ResponseWriter, data any, requestStatusCode int, err error, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(requestStatusCode)

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}

	log.Error().Err(err).Msg(msg)
}

func writeSuccess(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
