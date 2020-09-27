package controller

import (
	"encoding/json"
	"net/http"
)

// WriteError sends a response with { error, trace }
func WriteError(w http.ResponseWriter, statusCode int, short string, long string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	body := map[string]string{
		"error": short,
		"trace": long,
	}
	encodedBody, _ := json.Marshal(body)
	w.Write(encodedBody)
}

// WriteResponse sends a response with the provided body
func WriteResponse(w http.ResponseWriter, body interface{}) {
	if body != nil {
		w.Header().Set("Content-Type", "application/json")
		encodedBody, err := json.Marshal(body)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "unable to marshal body", err.Error())
			return
		}
		w.Write(encodedBody)
	}
}
