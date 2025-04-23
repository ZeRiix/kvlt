package routes

import (
	"kvlt/store"
	"kvlt/utils"
	"net/http"

	"github.com/gorilla/mux"
)

func getValueRoute(router *mux.Router) {
	router.HandleFunc("/value/{key}", func(w http.ResponseWriter, r *http.Request) {
		var requestRouteParams struct {
			Key string `string:"key"`
		}

		if !utils.ValidatePathParams(w, r, &requestRouteParams) {
			return
		}

		s := store.Get()
		value, exists := s.GetValue(requestRouteParams.Key)

		if !exists {
			utils.ErrorResponse(w, http.StatusNotFound, "Key not found")
			return
		}

		utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
			"key":   requestRouteParams.Key,
			"value": value,
		})
	}).Methods("GET")
}
