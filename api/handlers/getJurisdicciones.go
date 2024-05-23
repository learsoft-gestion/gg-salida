package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func GetJurisdicciones(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		idConvenio := vars["idConvenio"]
		idEmpresa := vars["idEmpresa"]

		query := "select distinct upper(nombre) from extractor.ext_modelos where 1 = 1"

		if idConvenio != "0" {
			query += " and id_convenio = " + idConvenio
		}
		if idEmpresa != "0" {
			query += " and id_empresa_adm = " + idEmpresa
		}
		query += " order by 1"

		rows, err := db.Query(query)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("Error al ejecutar query: ", err.Error())
			http.Error(w, "Error en el servidor: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var jurisdicciones []string
		for rows.Next() {
			var jurisdiccion string
			rows.Scan(&jurisdiccion)

			jurisdicciones = append(jurisdicciones, jurisdiccion)
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(jurisdicciones); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
