package routes

import (
	"kvlt/store"
	"kvlt/utils"
	"net/http"

	"github.com/gorilla/mux"
)

func setValueRoute(router *mux.Router) {
	router.HandleFunc("/value", func(w http.ResponseWriter, r *http.Request) {
		var requestBody struct {
			Key   string      `json:"key"`
			Value interface{} `json:"value"`
		}

		if !utils.ValidateRequestBody(w, r, &requestBody) {
			return
		}

		if requestBody.Key == "" {
			utils.ErrorResponse(w, http.StatusBadRequest, "Missing 'key' in request body")
			return
		}

		s := store.Get()
		s.SetValue(requestBody.Key, requestBody.Value)

		utils.JSONResponse(w, http.StatusCreated, map[string]interface{}{
			"key":   requestBody.Key,
			"value": requestBody.Value,
		})
	}).Methods("POST")
}
