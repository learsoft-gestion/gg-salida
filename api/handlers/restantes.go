package handlers

import (
	"Nueva/modelos"
	"Nueva/src"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var Restantes modelos.Restantes

func ProcesosRestantes(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		query := fmt.Sprintf("select modelo.id_modelo from extractor.ext_modelos modelo where modelo.id_convenio = %v and not exists  (select 1 from extractor.ext_procesados ep where ep.id_modelo = modelo.id_modelo and ep.fecha_desde = '%s' and ep.fecha_hasta = '%s')", id_convenio, fechaFormateada, fechaFormateada2)

		if len(id_empresa) > 0 {
			query += fmt.Sprintf(" and modelo.id_empresa_adm = %s", id_empresa)
		}
		if len(id_concepto) > 0 {
			query += fmt.Sprintf(" and modelo.id_concepto = '%s'", id_concepto)
		}
		if len(id_tipo) > 0 {
			query += fmt.Sprintf(" and modelo.id_tipo = '%s'", id_tipo)
		}
		if len(procesadoStr) > 0 {
			query += fmt.Sprintf(" and coalesce(archivo_salida is not null, false) = %v", procesado)
		}
		if len(jurisdiccion) > 0 {
			query += " and UPPER(modelo.nombre) like '%" + strings.ToUpper(jurisdiccion) + "%'"
		}

		rows, err := db.Query(query)
		if err != nil {
			fmt.Println(fmt.Println("Error al ejecutar query de procesosRestantes"))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var id_modelos []int
		var nombre_conv string
		var nombre_emp string
		var nombre_concepto string
		var nombre_tipo string

		for rows.Next() {
			var id_modelo int
			rows.Scan(&id_modelo)
			id_modelos = append(id_modelos, id_modelo)
		}

		err = db.QueryRow(fmt.Sprintf("SELECT nombre from extractor.ext_convenios where id_convenio = %v", id_convenio)).Scan(&nombre_conv)
		if err != nil {
			fmt.Println(err.Error())
		}

		if id_empresa != "" {
			err = db.QueryRow(fmt.Sprintf("select razon_social from datos.empresas_adm where id_empresa_adm = %s", id_empresa)).Scan(&nombre_emp)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		if id_concepto != "" {
			err = db.QueryRow(fmt.Sprintf("select ec.nombre from extractor.ext_conceptos ec where id_concepto = '%s'", strings.ToUpper(id_concepto))).Scan(&nombre_concepto)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		if id_tipo != "" {
			err = db.QueryRow(fmt.Sprintf("select et.nombre from extractor.ext_tipos et where id_tipo = '%s'", strings.ToUpper(id_tipo))).Scan(&nombre_tipo)
			if err != nil {
				fmt.Println(err.Error())
			}
		}

		Restantes = modelos.Restantes{
			Id:       id_modelos,
			Convenio: nombre_conv,
			Fecha1:   fechaFormateada,
			Fecha2:   fechaFormateada2,
		}
		cantidad := len(id_modelos)
		var resString string
		var btn string
		if cantidad == 0 {
			resString = fmt.Sprintf("Ya han sido generados todos los informes para el convenio %s", nombre_conv)
			if nombre_emp != "" {
				resString += ", empresa " + nombre_emp
			}
			if nombre_concepto != "" {
				resString += ", concepto " + nombre_concepto
			}
			if nombre_tipo != "" {
				resString += ", tipo " + nombre_tipo
			}
			if jurisdiccion != "" {
				resString += ", jurisdiccion " + jurisdiccion
			}
		} else {
			resString = fmt.Sprintf("Faltan generar %v informes para el convenio %s", cantidad, nombre_conv)
			if nombre_emp != "" {
				resString += ", empresa " + nombre_emp
			}
			if nombre_concepto != "" {
				resString += ", concepto " + nombre_concepto
			}
			if nombre_tipo != "" {
				resString += ", tipo " + nombre_tipo
			}
			if jurisdiccion != "" {
				resString += ", jurisdiccion " + jurisdiccion
			}
			btn = "Generar"
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		respuesta := modelos.RespuestaRestantes{
			Mensaje: resString,
			Boton:   btn,
		}

		jsonResp, _ := json.Marshal(respuesta)
		w.Write(jsonResp)

	}
}
