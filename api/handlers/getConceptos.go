package handlers

import (
	"Nueva/modelos"
	"Nueva/src"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func GetConceptos(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extraigo IDs de la request
		vars := mux.Vars(r)

		id_convenio := vars["id_convenio"]
		id_empresa := vars["id_empresa"]
		query := "select ec.id_concepto, ec.nombre as nombre_concepto, et.id_tipo, et.nombre as nombre_tipo from extractor.ext_modelos em join extractor.ext_conceptos ec on em.id_concepto = ec.id_concepto join extractor.ext_tipos et on em.id_tipo = et.id_tipo"
		if len(id_convenio) > 0 && len(id_empresa) > 0 {
			query += fmt.Sprintf(" where em.id_convenio = %s and em.id_empresa_adm = %s", id_convenio, id_empresa)
		} else if len(id_convenio) > 0 {
			query += fmt.Sprintf(" where em.id_convenio = %s ", id_convenio)
		} else if len(id_empresa) > 0 {
			query += fmt.Sprintf(" where em.id_empresa_adm = %s", id_empresa)
		}
		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var dtoConceptos modelos.Conceptos
		var conceptos []modelos.Concepto
		var tipos []modelos.Concepto
		for rows.Next() {
			var id_concepto string
			var concepto string
			var id_tipo string
			var tipo string
			if err = rows.Scan(&id_concepto, &concepto, &id_tipo, &tipo); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			conceptos = src.AddToSetConceptos(conceptos, modelos.Concepto{Id: id_concepto, Nombre: concepto})
			tipos = src.AddToSetConceptos(tipos, modelos.Concepto{Id: id_tipo, Nombre: tipo})

		}
		dtoConceptos = modelos.Conceptos{
			Conceptos: conceptos,
			Tipos:     tipos,
		}
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(dtoConceptos); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
