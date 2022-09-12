package server

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

func NotFound(rw http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") == "application/json" {
		JSON(rw, http.StatusNotFound, notFoundResponse)
	} else {
		String(rw, http.StatusNotFound, "404 page not found\n")
	}
}

func MethodNotAllowed(rw http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") == "application/json" {
		JSON(rw, http.StatusMethodNotAllowed, methodNotAllowed)
	} else {
		String(rw, http.StatusMethodNotAllowed, "405 method not allowed\n")
	}
}

func String(rw http.ResponseWriter, code int, text string) {
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(code)
	if _, err := rw.Write([]byte(text)); err != nil {
		log.Error().Err(err).Msg("Writing response")
	}
}

func JSON(rw http.ResponseWriter, code int, obj any) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	jsonBytes, errMarshal := json.Marshal(obj)
	if errMarshal != nil {
		log.Error().Err(errMarshal).Msg("JSON marshling")
	}
	if _, err := rw.Write(jsonBytes); err != nil {
		log.Error().Err(err).Msg("Writing response")
	}
}

func ParseJSON(r *http.Request, obj any) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return nil
}
