package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func NotFound(rw http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") == "application/json" {
		JSON(rw, 404, notFoundResponse)
	} else {
		String(rw, 404, "404 page not found\n")
	}
}

func MethodNotAllowed(rw http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") == "application/json" {
		JSON(rw, 405, methodNotAllowed)
	} else {
		String(rw, 405, "405 method not allowed\n")
	}
}

func String(rw http.ResponseWriter, code int, text string) {
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(code)
	if _, err := rw.Write([]byte(text)); err != nil {
		fmt.Printf("Error on writing response: %s\n", err.Error())
	}
}

func JSON(rw http.ResponseWriter, code int, obj any) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	jsonBytes, errMarshal := json.Marshal(obj)
	if errMarshal != nil {
		fmt.Printf("JSON marshling error: %s\n", errMarshal.Error())
	}
	if _, err := rw.Write(jsonBytes); err != nil {
		fmt.Printf("Error on writing response: %s\n", err.Error())
	}
}

func ParseJSON(r *http.Request, obj any) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return nil
}