package handlers

import (
	"Nueva/modelos"
	"Nueva/src"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var DTOprocesos []modelos.DTOproceso
var Procesos []modelos.Proceso

func GetProcesos(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Reinicio variables
		DTOprocesos = nil
		Procesos = nil

		// Extraigo params de la request
		id_convenio := r.URL.Query().Get("convenio")
		id_empresa := r.URL.Query().Get("empresa")
		id_concepto := r.URL.Query().Get("concepto")
		id_tipo := r.URL.Query().Get("tipo")
		fecha1 := r.URL.Query().Get("fecha1")
		fechaFormateada := src.FormatoFecha(fecha1)
		fecha2 := r.URL.Query().Get("fecha2")
		fechaFormateada2 := src.FormatoFecha(fecha2)
		jurisdiccion := r.URL.Query().Get("jurisdiccion")
		procesadoStr := r.URL.Query().Get("procesado")
		var procesado bool
		if procesadoStr == "true" {
			procesado = true
		} else if procesadoStr == "false" {
			procesado = false
		}

		if len(fecha1) == 0 || len(fecha2) == 0 {
			http.Error(w, "Fecha desde y fecha hasta son obligatorios", http.StatusInternalServerError)
			return
		}

		query := fmt.Sprintf("select em.id_modelo, c.nombre as nombre_convenio, ea.reducido as nombre_empresa_adm, ec.nombre as nombre_concepto, em.nombre, et.nombre as nombre_tipo, ep.fecha_desde, ep.fecha_hasta, ep.num_version, ep.archivo_salida, ep.fecha_ejecucion,	case when fecha_ejecucion is null then 'lanzar' when fecha_ejecucion = max(ep.fecha_ejecucion) over(partition by em.id_modelo) then 'relanzar' end boton, ep.archivo_nomina, ep.archivo_control, ep.id_proceso, ep.bloqueado from extractor.ext_modelos em left join extractor.ext_procesados ep on em.id_modelo = ep.id_modelo and ep.fecha_desde = '%s' and ep.fecha_hasta = '%s' join extractor.ext_empresas_adm ea ON em.id_empresa_adm = ea.id_empresa_adm join extractor.ext_convenios c ON em.id_convenio = c.id_convenio join extractor.ext_conceptos ec on em.id_concepto = ec.id_concepto join extractor.ext_tipos et on em.id_tipo = et.id_tipo where em.vigente", fechaFormateada, fechaFormateada2)

		if len(id_convenio) > 0 {
			query += fmt.Sprintf(" and em.id_convenio = %v", id_convenio)
		}
		if len(id_empresa) > 0 {
			query += fmt.Sprintf(" and em.id_empresa_adm = %s", id_empresa)
		}
		if len(id_concepto) > 0 {
			query += fmt.Sprintf(" and em.id_concepto = '%s'", id_concepto)
		}
		if len(id_tipo) > 0 {
			query += fmt.Sprintf(" and em.id_tipo = '%s'", id_tipo)
		}
		if len(procesadoStr) > 0 {
			query += fmt.Sprintf(" and coalesce(archivo_salida is not null, false) = %v", procesado)
		}
		if len(jurisdiccion) > 0 {
			query += " and UPPER(em.nombre) like '%" + strings.ToUpper(jurisdiccion) + "%'"
		}
		query += " ORDER BY nombre_empresa_adm, nombre_concepto, nombre, nombre_tipo, ep.num_version desc;"
		rows, err := db.Query(query)
		if err != nil {
			fmt.Println("Fallo el query: ", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			// var proceso modelos.Proceso
			var DTOproceso modelos.DTOproceso
			var fecha_desde sql.NullString
			var fecha_hasta sql.NullString
			var version sql.NullInt16
			var nombre_salida sql.NullString
			var ult_ejecucion sql.NullString
			var boton sql.NullString
			var nombre_nomina sql.NullString
			var nombre_control sql.NullString
			var id_proceso sql.NullInt32

			if err := rows.Scan(&DTOproceso.Id_modelo, &DTOproceso.Convenio, &DTOproceso.Empresa, &DTOproceso.Concepto, &DTOproceso.Nombre, &DTOproceso.Tipo, &fecha_desde, &fecha_hasta, &version, &nombre_salida, &ult_ejecucion, &boton, &nombre_nomina, &nombre_control, &id_proceso, &DTOproceso.Bloqueado); err != nil {
				println("Fallo el scan", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if fecha_desde.Valid {
				DTOproceso.Fecha_desde = fecha_desde.String
			}
			if fecha_hasta.Valid {
				DTOproceso.Fecha_hasta = fecha_hasta.String
			}
			if version.Valid {
				DTOproceso.Version = fmt.Sprintf("%v", version.Int16)
			}
			if ult_ejecucion.Valid {
				fecha_ult_ejecucion, err := time.Parse("2006-01-02T15:04:05.999999Z", ult_ejecucion.String)
				if err != nil {
					fmt.Println("Error al parsear fecha de ultima ejecucion")
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				DTOproceso.Ultima_ejecucion = fecha_ult_ejecucion.Format("2006-01-02 15:04")

				DTOproceso.Nombre_salida = nombre_salida.String

			}
			if nombre_nomina.Valid {
				DTOproceso.Nombre_nomina = nombre_nomina.String
			}
			if nombre_control.Valid {
				DTOproceso.Nombre_control = nombre_control.String
			}
			if boton.Valid {
				DTOproceso.Boton = boton.String
			}
			if id_proceso.Valid {
				DTOproceso.Id_procesado = int(id_proceso.Int32)
			}
			DTOprocesos = append(DTOprocesos, DTOproceso)

		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(DTOprocesos); err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}
