package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func MigradorGetConvenios(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := "SELECT DISTINCT convenio FROM extractor.ext_archivos"

		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, "Error en la consulta "+err.Error(), http.StatusInternalServerError)
			return
		}

		var convenios []string

		for rows.Next() {
			var convenio string
			err := rows.Scan(&convenio)
			if err != nil {
				http.Error(w, "Error al escanear fila "+err.Error(), http.StatusInternalServerError)
				return
			}
			convenios = append(convenios, convenio)
		}

		jsonData, err := json.Marshal(convenios)
		if err != nil {
			http.Error(w, "Error al convertir a JSON "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	}
}
