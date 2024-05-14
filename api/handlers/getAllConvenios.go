package handlers

import (
	"Nueva/modelos"
	"database/sql"
	"encoding/json"
	"net/http"
)

func GetAllConvenios(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id_convenio, nombre FROM extractor.ext_convenios order by nombre")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var convenios []modelos.Convenio
		for rows.Next() {
			var convenio modelos.Convenio
			rows.Scan(&convenio.Id, &convenio.Nombre)

			convenios = append(convenios, convenio)
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(convenios); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
