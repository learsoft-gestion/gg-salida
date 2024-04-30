package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func MigradorGetEmpresas(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := "SELECT DISTINCT empresa FROM extractor.ext_archivos"

		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, "Error en la consulta "+err.Error(), http.StatusInternalServerError)
			return
		}

		var empresas []string

		for rows.Next() {
			var empresa string
			err := rows.Scan(&empresa)
			if err != nil {
				http.Error(w, "Error al escanear fila "+err.Error(), http.StatusInternalServerError)
				return
			}
			empresas = append(empresas, empresa)
		}

		jsonData, err := json.Marshal(empresas)
		if err != nil {
			http.Error(w, "Error al convertir a JSON "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	}
}
