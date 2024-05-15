package handlers

import (
	"Nueva/modelos"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func GetAlicuotas(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extraigo IDs de la request
		vars := mux.Vars(r)

		idConvenio := vars["idConvenio"]
		query := "select id_alicuota, id_convenio, nombre, descripcion from extractor.ext_alicuotas where id_convenio = $1 order by nombre asc"

		rows, err := db.Query(query, idConvenio)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var alicuotas []modelos.Alicuota
		for rows.Next() {
			var alicuota modelos.Alicuota
			rows.Scan(&alicuota.IdAlicuota, &alicuota.IdConvenio, &alicuota.Nombre, &alicuota.Descripcion)

			alicuotas = append(alicuotas, alicuota)
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(alicuotas); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
