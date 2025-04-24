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
			utils.MakeResponse(w, utils.Response{
				Status: http.StatusNotFound,
				Error:  "Key not found",
				Info:   "error.keyNotFound",
			})
			return
		}

		utils.MakeResponse(w, utils.Response{
			Status: http.StatusOK,
			Data: map[string]interface{}{
				"key":   requestRouteParams.Key,
				"value": value,
			},
			Info: "success.keyFound",
		})
	}).Methods("GET")
}
