package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func MigradorGetPeriodos(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := "SELECT DISTINCT periodo FROM extractor.ext_archivos"

		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, "Error en la consulta "+err.Error(), http.StatusInternalServerError)
			return
		}

		var periodos []string

		for rows.Next() {
			var periodo string
			err := rows.Scan(&periodo)
			if err != nil {
				http.Error(w, "Error al escanear fila "+err.Error(), http.StatusInternalServerError)
				return
			}
			periodos = append(periodos, periodo)
		}

		jsonData, err := json.Marshal(periodos)
		if err != nil {
			http.Error(w, "Error al convertir a JSON "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	}
}
