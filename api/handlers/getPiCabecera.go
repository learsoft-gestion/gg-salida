package handlers

import (
	"Nueva/modelos"
	"database/sql"
	"encoding/json"
	"net/http"
)

func GetPiCabecera(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := "select id_pi, id_empresa_adm, periodo from extractor.ext_pi_cabecera"

		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var piCabeceras []modelos.PiCabecera
		for rows.Next() {
			var piCabecera modelos.PiCabecera
			rows.Scan(&piCabecera.IdPi, &piCabecera.IdEmpresaAdm, &piCabecera.Periodo)
			piCabeceras = append(piCabeceras, piCabecera)
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(piCabeceras); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
