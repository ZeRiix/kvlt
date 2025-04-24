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
			Key      string      `json:"key"`
			Value    interface{} `json:"value"`
			Duration int64       `json:"duration"`
		}

		if !utils.ValidateRequestBody(w, r, &requestBody) {
			return
		}

		if requestBody.Key == "" {
			utils.MakeResponse(w, utils.Response{
				Status: http.StatusBadRequest,
				Error:  "Key is required",
				Info:   "error.keyRequired",
			})
			return
		}

		s := store.Get()
		s.SetValue(requestBody.Key, requestBody.Value, requestBody.Duration)

		utils.MakeResponse(w, utils.Response{
			Status: http.StatusOK,
			Info:   "success.keySet",
		})
	}).Methods("PUT")
}
