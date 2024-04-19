package handlers

import (
	"Nueva/modelos"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

var Models []modelos.Modelo

func ModelosHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PATCH" {
			var vigente bool
			idURL := r.URL.Query().Get("id")
			vigenteURL := r.URL.Query().Get("vigente")

			id, err := strconv.Atoi(idURL)
			if err != nil {
				http.Error(w, "Debe enviar un ID numerico", http.StatusBadRequest)
				return
			}

			if vigenteURL == "true" {
				vigente = true
			} else if vigenteURL == "false" {
				vigente = false
			} else {
				http.Error(w, "Debe enviar vigente en formato bool y un ID numerico", http.StatusBadRequest)
				return
			}

			query := "UPDATE extractor.ext_modelos SET vigente = $1 where id_modelo = $2"

			result, err := db.Exec(query, vigente, id)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("Error al ejecutar query: ", err.Error())
				http.Error(w, "Error en el servidor", http.StatusBadRequest)
				return
			}

			if cuenta, err := result.RowsAffected(); cuenta > 1 {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				respuesta := "Modelo actualizado exitosamente"
				jsonResp, _ := json.Marshal(respuesta)
				w.Write(jsonResp)
				return
			} else if err != nil {
				fmt.Println("Error al actualizar modelo: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				http.Error(w, "Error en el servidor", http.StatusInternalServerError)
			}
		} else if r.Method == "POST" {

			var model modelos.Modelo
			err := json.NewDecoder(r.Body).Decode(&model)
			if err != nil {
				fmt.Println(err.Error())
				http.Error(w, "Error decodificando JSON", http.StatusBadRequest)
				return
			}

			query := "INSERT INTO extractor.ext_modelos (id_empresa_adm, id_convenio, id_concepto, id_tipo, nombre, filtro_personas, filtro_recibos, filtro_having, formato_salida, archivo_control, archivo_modelo, archivo_nomina) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)"

			// fmt.Println(query)

			result, err := db.Exec(query, model.Id_empresa, model.Id_convenio, model.Id_concepto, model.Id_tipo, model.Nombre, model.Filtro_personas, model.Filtro_recibos, model.Filtro_having, model.Formato_salida, model.Archivo_control, model.Archivo_modelo, model.Archivo_nomina)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("Error al ejecutar query: ", err.Error())
				http.Error(w, "Error en el servidor", http.StatusInternalServerError)
				return
			}

			if cuenta, err := result.RowsAffected(); cuenta > 1 {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				respuesta := fmt.Sprintf("Modelo %s creado exitosamente", model.Nombre)
				jsonResp, _ := json.Marshal(respuesta)
				w.Write(jsonResp)
				return
			} else if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("No se han insertado registros: ", err.Error())
				http.Error(w, "Error en el servidor", http.StatusInternalServerError)

			}

		} else {

			Models = nil

			id_convenio := r.URL.Query().Get("convenio")
			vigente := r.URL.Query().Get("vigente")

			query := "select em.id_modelo, em.id_empresa_adm, em.id_concepto, em.id_convenio, em.id_tipo, ea.razon_social as nombre_empresa_adm, ec.nombre as nombre_concepto, c.nombre as nombre_convenio, et.nombre as nombre_tipo, em.nombre, em.filtro_personas, em.filtro_recibos, em.formato_salida, em.ult_ejecucion, em.id_query, em.archivo_modelo, em.vigente, em.filtro_having, em.archivo_control, em.archivo_nomina from extractor.ext_modelos em join datos.empresas_adm ea ON em.id_empresa_adm = ea.id_empresa_adm join extractor.ext_convenios c ON em.id_convenio = c.id_convenio join extractor.ext_conceptos ec on em.id_concepto = ec.id_concepto join extractor.ext_tipos et on em.id_tipo = et.id_tipo "

			if id_convenio != "" {
				query += "where em.id_convenio = " + id_convenio
			}
			if vigente == "true" {
				if id_convenio != "" {
					query += " and em.vigente"
				} else {
					query += "where vigente"
				}
			} else if vigente == "false" {
				if id_convenio != "" {
					query += " and em.vigente = false"
				} else {
					query += "where em.vigente = false"
				}
			}

			query += " order by em.id_modelo"

			// fmt.Println(query)

			rows, err := db.Query(query)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			for rows.Next() {
				var modelo modelos.Modelo
				var ult_ejecucion sql.NullTime
				var filtroPersonas sql.NullString
				var filtroRecibos sql.NullString
				var filtroHaving sql.NullString

				if err = rows.Scan(&modelo.Id_modelo, &modelo.Id_empresa, &modelo.Id_concepto, &modelo.Id_convenio, &modelo.Id_tipo, &modelo.Empresa, &modelo.Concepto, &modelo.Convenio, &modelo.Tipo, &modelo.Nombre, &filtroPersonas, &filtroRecibos, &modelo.Formato_salida, &ult_ejecucion, &modelo.Query, &modelo.Archivo_modelo, &modelo.Vigente, &filtroHaving, &modelo.Archivo_control, &modelo.Archivo_nomina); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				if ult_ejecucion.Valid {
					modelo.Ultima_ejecucion = ult_ejecucion.Time.String()
				}
				if filtroPersonas.Valid {
					modelo.Filtro_personas = filtroPersonas.String
				}
				if filtroRecibos.Valid {
					modelo.Filtro_recibos = filtroRecibos.String
				}
				if ult_ejecucion.Valid {
					modelo.Filtro_having = filtroHaving.String
				}

				Models = append(Models, modelo)

			}

			w.Header().Set("Content-Type", "application/json")

			if err := json.NewEncoder(w).Encode(Models); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}
