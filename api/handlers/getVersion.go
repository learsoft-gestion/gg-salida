package handlers

import (
	"encoding/json"
	"net/http"
)

func GetVersion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		version := "Version 1.3"

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(version); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
