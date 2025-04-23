package utils

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func JSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func ErrorResponse(w http.ResponseWriter, status int, message string) {
	JSONResponse(w, status, map[string]string{"error": message})
}

func ParseJSONBody(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func ValidateRequestBody(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := ParseJSONBody(r, v); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return false
	}
	return true
}

func ParseQueryParams(r *http.Request, v interface{}) error {
	query := r.URL.Query()
	if err := json.Unmarshal([]byte(query.Encode()), v); err != nil {
		return err
	}
	return nil
}

func ValidateQueryParams(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := ParseQueryParams(r, v); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid query parameters")
		return false
	}
	return true
}

func ParsePathParams(r *http.Request, v interface{}) error {
	vars := mux.Vars(r)
	jsonBytes, err := json.Marshal(vars)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(jsonBytes, v); err != nil {
		return err
	}
	return nil
}

func ValidatePathParams(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := ParsePathParams(r, v); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "Invalid path parameters")
		return false
	}
	return true
}
