package handlers

import (
	"Nueva/modelos"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func GetPiDetalle(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		idPi := vars["idPi"]
		query := "select id_pi, cuil, fecha_ingreso, remuneracion_total, categoria, descuenta_cuota_sindical from extractor.ext_pi_detalle where id_pi = $1"

		rows, err := db.Query(query, idPi)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var piDetalles []modelos.PiDetalle
		for rows.Next() {
			var piDetalle modelos.PiDetalle
			rows.Scan(&piDetalle.IdPi, &piDetalle.Cuil, &piDetalle.FechaIngreso, &piDetalle.RemTotal, &piDetalle.Categoria, &piDetalle.DescuentaCuotaSindical)
			piDetalles = append(piDetalles, piDetalle)
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(piDetalles); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
