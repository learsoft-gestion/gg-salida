package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func GetPersonalInterno(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := "select cuil from extractor.ext_personal_interno order by 1"

		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var cuils []string
		for rows.Next() {
			var cuil string
			rows.Scan(&cuil)

			cuils = append(cuils, cuil)
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(cuils); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
