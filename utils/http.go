package utils

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Response struct {
	Status int         `json:"-"`               // HTTP status code
	Info   string      `json:"info,omitempty"`  // Information add in header
	Data   interface{} `json:"data,omitempty"`  // Data to be returned in the response
	Error  string      `json:"error,omitempty"` // Error message to be returned in the response
}

// HTTP utility functions for handling JSON responses, request body parsing, and query parameter parsing.
func MakeResponse(w http.ResponseWriter, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)

	responseMap := make(map[string]interface{})

	if response.Info != "" {
		w.Header().Set("information", response.Info)
	}

	if response.Data != nil {
		responseMap["data"] = response.Data
	}

	if response.Error != "" {
		responseMap["error"] = response.Error
	}

	json.NewEncoder(w).Encode(responseMap)
}

// ParseJSONBody parses the JSON body of an HTTP request into the provided struct.
func ParseJSONBody(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// ValidateRequestBody checks if the request body is valid JSON and unmarshals it into the provided struct.
func ValidateRequestBody(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := ParseJSONBody(r, v); err != nil {
		MakeResponse(w, Response{
			Status: http.StatusBadRequest,
			Error:  "Invalid request body",
			Info:   "error.bodyParams",
		})
		return false
	}
	return true
}

// ParseQueryParams parses the query parameters from the URL into the provided struct.
func ParseQueryParams(r *http.Request, v interface{}) error {
	query := r.URL.Query()
	if err := json.Unmarshal([]byte(query.Encode()), v); err != nil {
		return err
	}
	return nil
}

// ValidateQueryParams checks if the query parameters are valid and unmarshals them into the provided struct.
func ValidateQueryParams(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := ParseQueryParams(r, v); err != nil {
		MakeResponse(w, Response{
			Status: http.StatusBadRequest,
			Error:  "Invalid query parameters",
			Info:   "error.queryParams",
		})
		return false
	}
	return true
}

// ParsePathParams parses the path parameters from the URL into the provided struct.
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

// ValidatePathParams checks if the path parameters are valid and unmarshals them into the provided struct.
func ValidatePathParams(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := ParsePathParams(r, v); err != nil {
		MakeResponse(w, Response{
			Status: http.StatusBadRequest,
			Error:  "Invalid path parameters",
			Info:   "error.routeParams",
		})
		return false
	}
	return true
}
