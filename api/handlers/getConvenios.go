package handlers

import (
	"Nueva/modelos"
	"Nueva/src"
	"database/sql"
	"encoding/json"
	"net/http"
)

var Convenios []modelos.Option

func GetConvenios(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rows, err := db.Query("SELECT em.id_convenio, c.nombre, c.filtro FROM extractor.ext_modelos em JOIN extractor.ext_convenios c ON em.id_convenio = c.id_convenio where vigente order by c.nombre")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for rows.Next() {
			var id int
			var convenio string
			var filtro sql.NullString
			var filtroJson string

			if err = rows.Scan(&id, &convenio, &filtro); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if filtro.Valid {
				filtroJson = filtro.String
			}

			Convenios = src.AddToSet(Convenios, modelos.Option{Id: id, Nombre: convenio, Filtro: filtroJson})

		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(Convenios); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
