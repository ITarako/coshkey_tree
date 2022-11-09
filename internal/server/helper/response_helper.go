package serverhelper

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

func WriteError(w http.ResponseWriter, data any, responseStatusCode int, err error, errorMsg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(responseStatusCode)

	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}

	log.Error().Err(err).Msg(errorMsg)
}

func WriteSuccess(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(data)
}
