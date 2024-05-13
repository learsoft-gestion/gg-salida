package handlers

import (
	"Nueva/modelos"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func GetValoresAlicuotas(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extraigo IDs de la request
		vars := mux.Vars(r)

		idAlicuota := vars["idAlicuota"]
		query := "select id_valores_alicuota, id_alicuota, vigencia_desde, valor from extractor.ext_valores_alicuotas where id_alicuota = $1"

		rows, err := db.Query(query, idAlicuota)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var valoresAlicuotas []modelos.ValorAlicuota
		for rows.Next() {
			var valorAlicuota modelos.ValorAlicuota
			rows.Scan(&valorAlicuota.IdValoresAlicuota, &valorAlicuota.IdAlicuota, &valorAlicuota.VigenciaDesde, &valorAlicuota.Valor)

			valoresAlicuotas = append(valoresAlicuotas, valorAlicuota)
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(valoresAlicuotas); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
