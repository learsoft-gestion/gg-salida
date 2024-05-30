package handlers

import (
	"Nueva/modelos"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var Models []modelos.Modelo

func ModelosHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "PATCH" {

			var model modelos.PatchModelo
			err := json.NewDecoder(r.Body).Decode(&model)
			if err != nil {
				fmt.Println(err.Error())
				http.Error(w, "Error decodificando JSON: "+err.Error(), http.StatusBadRequest)
				return
			}

			query := "UPDATE extractor.ext_modelos SET vigente = $1 where id_modelo = $2"

			result, err := db.Exec(query, model.Vigente, model.Id)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("Error al ejecutar query: ", err.Error())
				http.Error(w, "Error en el servidor: "+err.Error(), http.StatusBadRequest)
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
				http.Error(w, "Error en el servidor: "+err.Error(), http.StatusInternalServerError)
			}
		} else if r.Method == "POST" {

			var model modelos.Modelo
			err := json.NewDecoder(r.Body).Decode(&model)
			if err != nil {
				fmt.Println(err.Error())
				http.Error(w, "Error decodificando JSON: "+err.Error(), http.StatusBadRequest)
				return
			}

			query := "INSERT INTO extractor.ext_modelos (id_empresa_adm, id_convenio, id_concepto, id_tipo, nombre, filtro_personas, filtro_recibos, formato_salida, archivo_modelo, archivo_nomina) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)"

			// fmt.Println(query)

			result, err := db.Exec(query, model.Id_empresa, model.Id_convenio, model.Id_concepto, model.Id_tipo, model.Nombre, model.Filtro_personas, model.Filtro_recibos, model.Formato_salida, model.Archivo_modelo, model.Archivo_nomina)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("Error al ejecutar query: ", err.Error())
				http.Error(w, "Error en el servidor: "+err.Error(), http.StatusInternalServerError)
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
				http.Error(w, "Error en el servidor: "+err.Error(), http.StatusInternalServerError)

			}

		} else {

			Models = nil

			idConvenio := r.URL.Query().Get("convenio")
			vigente := r.URL.Query().Get("vigente")
			idEmpresa := r.URL.Query().Get("empresa")
			idConcepto := r.URL.Query().Get("concepto")
			idTipo := r.URL.Query().Get("tipo")
			jurisdiccion := r.URL.Query().Get("jurisdiccion")

			query := "select em.id_modelo, em.id_empresa_adm, em.id_concepto, em.id_convenio, em.id_tipo, ea.razon_social, ea.reducido, ec.nombre as nombre_concepto, c.nombre as nombre_convenio, et.nombre as nombre_tipo, em.nombre, c.filtro, em.filtro_personas, em.filtro_recibos, em.formato_salida, em.ult_ejecucion, em.id_query, em.archivo_modelo, em.vigente, em.archivo_nomina, regexp_replace( regexp_replace( regexp_replace(em.columna_estado, '<', '&lt;', 'g'), '>', '&gt;', 'g' ), E'\\n', '<BR>', 'g' ) AS columna_estado, regexp_replace( regexp_replace( regexp_replace(em.select_control, '<', '&lt;', 'g'), '>', '&gt;', 'g' ), E'\\n', '<BR>', 'g' ) AS select_control FROM extractor.ext_modelos em JOIN extractor.ext_empresas_adm ea ON em.id_empresa_adm = ea.id_empresa_adm JOIN extractor.ext_convenios c ON em.id_convenio = c.id_convenio JOIN extractor.ext_conceptos ec ON em.id_concepto = ec.id_concepto join extractor.ext_tipos et ON em.id_tipo = et.id_tipo "

			// Construir las condiciones WHERE basadas en los parámetros recibidos
			var conditions []string

			if idConvenio != "" {
				conditions = append(conditions, fmt.Sprintf("em.id_convenio = %s", idConvenio))
			}

			if vigente == "true" {
				conditions = append(conditions, "em.vigente = true")
			} else if vigente == "false" {
				conditions = append(conditions, "em.vigente = false")
			}

			if idEmpresa != "" {
				conditions = append(conditions, fmt.Sprintf("em.id_empresa_adm = %s", idEmpresa))
			}

			if idConcepto != "" {
				conditions = append(conditions, fmt.Sprintf("em.id_concepto = '%s'", idConcepto))
			}

			if idTipo != "" {
				conditions = append(conditions, fmt.Sprintf("em.id_tipo = '%s'", idTipo))
			}

			if jurisdiccion != "" {
				conditions = append(conditions, fmt.Sprintf("UPPER(em.nombre) LIKE '%%%s%%'", strings.ToUpper(jurisdiccion)))
			}

			// Combinar todas las condiciones en una sola cadena
			if len(conditions) > 0 {
				query += " WHERE " + strings.Join(conditions, " AND ")
			}

			query += " order by c.nombre, ea.razon_social, nombre_concepto, nombre_tipo, em.nombre"

			// fmt.Println(query)

			rows, err := db.Query(query)
			if err != nil {
				fmt.Println(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			for rows.Next() {
				var modelo modelos.Modelo
				var ult_ejecucion sql.NullTime
				var filtroPersonas sql.NullString
				var filtroRecibos sql.NullString
				var columna_estado sql.NullString
				var select_control sql.NullString

				if err = rows.Scan(&modelo.Id_modelo, &modelo.Id_empresa, &modelo.Id_concepto, &modelo.Id_convenio, &modelo.Id_tipo, &modelo.Empresa, &modelo.EmpReducido, &modelo.Concepto, &modelo.Convenio, &modelo.Tipo, &modelo.Nombre, &modelo.Filtro_convenio, &filtroPersonas, &filtroRecibos, &modelo.Formato_salida, &ult_ejecucion, &modelo.Query, &modelo.Archivo_modelo, &modelo.Vigente, &modelo.Archivo_nomina, &columna_estado, &select_control); err != nil {
					fmt.Println(err.Error())
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
				if columna_estado.Valid {
					modelo.Columna_estado = columna_estado.String
				}
				if select_control.Valid {
					modelo.Select_control = select_control.String
				}
				modelo.Ruta_archivo_modelo = "/templates/" + modelo.Archivo_modelo
				modelo.Ruta_archivo_nomina = "/templates/" + modelo.Archivo_nomina

				Models = append(Models, modelo)

			}

			w.Header().Set("Content-Type", "application/json")

			if err := json.NewEncoder(w).Encode(Models); err != nil {
				fmt.Println(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}
