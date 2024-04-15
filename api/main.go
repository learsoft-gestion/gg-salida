package main

import (
	"Nueva/conexiones"
	"Nueva/modelos"
	"Nueva/src"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var convenios []modelos.Option
var empresas []modelos.Option
var DTOprocesos []modelos.DTOproceso
var procesos []modelos.Proceso
var clientes []modelos.Cliente

// Almacena los registros restantes a ejecutar
var restantes modelos.Restantes

func getConvenios(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rows, err := db.Query("SELECT em.id_convenio, c.nombre as nombre_convenio FROM extractor.ext_modelos em JOIN extractor.ext_convenios c ON em.id_convenio = c.id_convenio where vigente order by c.nombre")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for rows.Next() {
			var id int
			var convenio string
			if err = rows.Scan(&id, &convenio); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			convenios = src.AddToSet(convenios, modelos.Option{Id: id, Nombre: convenio})

		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(convenios); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func getEmpresas(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var id_convenio string
		var err error

		query := "SELECT em.id_empresa_adm, ea.razon_social as nombre_empresa_adm FROM extractor.ext_modelos em JOIN datos.empresas_adm ea ON em.id_empresa_adm = ea.id_empresa_adm where vigente"
		vars := mux.Vars(r)

		if len(vars["id_convenio"]) > 0 {
			id_convenio = vars["id_convenio"]
			query += fmt.Sprintf(" and em.id_convenio = %s", id_convenio)
		}

		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		empresas = nil
		for rows.Next() {
			var id int
			var empresa string
			if err = rows.Scan(&id, &empresa); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			empresas = src.AddToSet(empresas, modelos.Option{Id: id, Nombre: empresa})

		}
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(empresas); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func getProcesos(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Reinicio variables
		DTOprocesos = nil
		procesos = nil

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

		if len(id_convenio) == 0 || len(fecha1) == 0 || len(fecha2) == 0 {
			http.Error(w, "Convenio, fecha1 y fecha2 son obligatorios", http.StatusInternalServerError)
			return
		}

		// query := fmt.Sprintf("select em.id_modelo, c.nombre as nombre_convenio, ea.razon_social as nombre_empresa_adm, ec.nombre as nombre_concepto, em.nombre, et.nombre as nombre_tipo, ep.fecha_desde, ep.fecha_hasta, ep.archivo_salida, ep.version, ep.fecha_ejecucion_salida, coalesce(archivo_salida is not null, false) procesado, case when version is null then 'lanzar' when version = max(ep.version) over(partition by em.id_modelo) then 'relanzar' end boton, case when ep.version >= 2 and ep.version = max(ep.version) over(partition by em.id_modelo) then true else false end as ult_version from extractor.ext_modelos em left join extractor.ext_procesados ep on em.id_modelo = ep.id_modelo and ep.fecha_desde = '%s' and ep.fecha_hasta = '%s' join datos.empresas_adm ea ON em.id_empresa_adm = ea.id_empresa_adm join extractor.ext_convenios c ON em.id_convenio = c.id_convenio join extractor.ext_conceptos ec on em.id_concepto = ec.id_concepto join extractor.ext_tipos et on em.id_tipo = et.id_tipo where em.id_convenio = %v", fechaFormateada, fechaFormateada2, id_convenio)
		query := fmt.Sprintf("select em.id_modelo, c.nombre as nombre_convenio, ea.razon_social as nombre_empresa_adm, ec.nombre as nombre_concepto, em.nombre, et.nombre as nombre_tipo, ep.fecha_desde, ep.fecha_hasta, ep.version, ep.archivo_salida, ep.fecha_ejecucion_salida,	case when fecha_ejecucion_salida is null then 'lanzar' when fecha_ejecucion_salida = max(ep.fecha_ejecucion_salida) over(partition by em.id_modelo) then 'relanzar' end boton_salida, ep.archivo_control, ep.fecha_ejecucion_control,	case when fecha_ejecucion_control is null then 'lanzar' when fecha_ejecucion_control = max(ep.fecha_ejecucion_control) over(partition by em.id_modelo) then 'relanzar' end boton_control, ep.id_proceso from extractor.ext_modelos em left join extractor.ext_procesados ep on em.id_modelo = ep.id_modelo and ep.fecha_desde = '%s' and ep.fecha_hasta = '%s' join datos.empresas_adm ea ON em.id_empresa_adm = ea.id_empresa_adm join extractor.ext_convenios c ON em.id_convenio = c.id_convenio join extractor.ext_conceptos ec on em.id_concepto = ec.id_concepto join extractor.ext_tipos et on em.id_tipo = et.id_tipo where em.id_convenio = %v", fechaFormateada, fechaFormateada2, id_convenio)

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
		query += " ORDER BY nombre_empresa_adm, nombre_concepto, nombre, nombre_tipo, ep.version desc;"
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
			var ult_ejecucion_salida sql.NullString
			var boton_salida sql.NullString
			var nombre_control sql.NullString
			var ult_ejecucion_control sql.NullString
			var boton_control sql.NullString
			var id_proceso sql.NullInt32

			if err := rows.Scan(&DTOproceso.Id_modelo, &DTOproceso.Convenio, &DTOproceso.Empresa, &DTOproceso.Concepto, &DTOproceso.Nombre, &DTOproceso.Tipo, &fecha_desde, &fecha_hasta, &version, &nombre_salida, &ult_ejecucion_salida, &boton_salida, &nombre_control, &ult_ejecucion_control, &boton_control, &id_proceso); err != nil {
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
			if ult_ejecucion_salida.Valid {
				fecha_ult_ejecucion, err := time.Parse("2006-01-02T15:04:05.999999Z", ult_ejecucion_salida.String)
				if err != nil {
					fmt.Println("Error al parsear fecha de ultima ejecucion")
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				DTOproceso.Ultima_ejecucion_salida = fecha_ult_ejecucion.Format("2006-01-02 15:04")

				DTOproceso.Nombre_salida = nombre_salida.String

			}
			if ult_ejecucion_control.Valid {
				fecha_ult_ejecucion, err := time.Parse("2006-01-02T15:04:05.999999Z", ult_ejecucion_control.String)
				if err != nil {
					fmt.Println("Error al parsear fecha de ultima ejecucion")
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				DTOproceso.Ultima_ejecucion_control = fecha_ult_ejecucion.Format("2006-01-02 15:04")

				DTOproceso.Nombre_control = nombre_control.String
			}
			if boton_salida.Valid {
				DTOproceso.Boton_salida = boton_salida.String
			}
			if boton_control.Valid {
				DTOproceso.Boton_control = boton_control.String
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

func getConceptos(db *sql.DB) http.HandlerFunc {
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

func sender(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			procesos = nil
			var datos modelos.DTOdatos
			err := json.NewDecoder(r.Body).Decode(&datos)
			if err != nil {
				http.Error(w, "Error decodificando JSON", http.StatusBadRequest)
				return
			}
			datos.Fecha = src.FormatoFecha(datos.Fecha)
			datos.Fecha2 = src.FormatoFecha(datos.Fecha2)

			queryModelos := "SELECT em.id_modelo, em.id_empresa_adm, ea.razon_social as nombre_empresa, c.id_convenio as id_convenio, c.nombre as nombre_convenio, em.nombre, c.filtro as filtro_convenio, em.filtro_personas, em.filtro_recibos, em.formato_salida, em.archivo_modelo, em.filtro_having FROM extractor.ext_modelos em JOIN datos.empresas_adm ea ON em.id_empresa_adm = ea.id_empresa_adm JOIN extractor.ext_convenios c ON em.id_convenio = c.id_convenio where vigente and em.id_modelo = $1"
			// fmt.Println("Query modelos: ", queryModelos)
			stmt, err := db.Prepare(queryModelos)
			if err != nil {
				http.Error(w, "Error al preparar query", http.StatusBadRequest)
				return
			}
			defer stmt.Close()
			var args []interface{}
			args = append(args, datos.Id_modelo)
			rows, err := stmt.Query(args...)
			if err != nil {
				http.Error(w, "Error al ejecutar el query", http.StatusBadRequest)
				return
			}
			defer rows.Close()

			for rows.Next() {
				var proceso modelos.Proceso
				err = rows.Scan(&proceso.Id_modelo, &proceso.Id_empresa, &proceso.Nombre_empresa, &proceso.Id_convenio, &proceso.Nombre_convenio, &proceso.Nombre, &proceso.Filtro_convenio, &proceso.Filtro_personas, &proceso.Filtro_recibos, &proceso.Formato_salida, &proceso.Archivo_modelo, &proceso.Filtro_having)
				if err != nil {
					fmt.Println(err.Error())
					http.Error(w, "Error al escanear proceso", http.StatusBadRequest)
					return
				}
				proceso.Id_procesado = datos.Id_procesado
				procesos = append(procesos, proceso)
			}

			var version int
			var archivo_salida bool
			if datos.Id_procesado > 0 {
				// Verificar si el proceso ya se corrió
				queryCuenta := fmt.Sprintf("select version, archivo_salida from extractor.ext_procesados where id_proceso = %v", datos.Id_procesado)
				// fmt.Println("Query: ", queryCuenta)
				var archivoSalida sql.NullString
				err = db.QueryRow(queryCuenta).Scan(&version, &archivoSalida)
				if err != nil {
					fmt.Println(err.Error())
					http.Error(w, "Error al escanear proceso", http.StatusBadRequest)
					return
				}
				if archivoSalida.Valid {
					archivo_salida = true
				}

				version += 1
			} else {
				version = 1
			}

			var resultado []string
			result, errFormateado := src.ProcesadorSalida(procesos[0], datos.Fecha, datos.Fecha2, version, archivo_salida)
			if result != "" {
				resultado = append(resultado, result)
			}
			if errFormateado.Mensaje != "" {
				errString := "Error en " + procesos[0].Nombre + ": " + errFormateado.Mensaje
				// http.Error(w, errString, http.StatusBadRequest)
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				respuesta := modelos.Respuesta{
					Mensaje:         errString,
					Archivos_salida: nil,
				}
				jsonResp, _ := json.Marshal(respuesta)
				w.Write(jsonResp)
				return
			}
			// El proceso termino, reinicio procesos
			procesos = nil

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			respuesta := modelos.Respuesta{
				Mensaje:         "Informe generado exitosamente",
				Archivos_salida: resultado,
			}
			jsonResp, _ := json.Marshal(respuesta)
			w.Write(jsonResp)
		}
	}
}

func multipleSend(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			procesos = nil

			var placeholders []string
			for i := range restantes.Id {
				placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
			}

			queryModelos := fmt.Sprintf("SELECT em.id_modelo, em.id_empresa_adm, ea.razon_social as nombre_empresa, c.nombre as nombre_convenio, em.nombre, c.filtro as filtro_convenio, em.filtro_personas, em.filtro_recibos, em.formato_salida, em.archivo_modelo FROM extractor.ext_modelos em JOIN datos.empresas_adm ea ON em.id_empresa_adm = ea.id_empresa_adm JOIN extractor.ext_convenios c ON em.id_convenio = c.id_convenio where vigente and em.id_modelo in (%s)", strings.Join(placeholders, ","))

			stmt, err := db.Prepare(queryModelos)
			if err != nil {
				http.Error(w, "Error al preparar query", http.StatusBadRequest)
				return
			}
			defer stmt.Close()
			var args []interface{}
			for _, arg := range restantes.Id {
				args = append(args, arg)
			}

			rows, err := stmt.Query(args...)
			if err != nil {
				http.Error(w, "Error al ejecutar el query", http.StatusBadRequest)
				return
			}
			defer rows.Close()
			for rows.Next() {
				var proceso modelos.Proceso
				err = rows.Scan(&proceso.Id_modelo, &proceso.Id_empresa, &proceso.Nombre_empresa, &proceso.Nombre_convenio, &proceso.Nombre, &proceso.Filtro_convenio, &proceso.Filtro_personas, &proceso.Filtro_recibos, &proceso.Formato_salida, &proceso.Archivo_modelo)
				if err != nil {
					fmt.Println(err.Error())
					http.Error(w, "Error al escanear proceso", http.StatusBadRequest)
					return
				}
				procesos = append(procesos, proceso)
			}

			var resultado []string
			for _, proc := range procesos {
				var archivoSalida bool
				var archivo_salida sql.NullString
				var version int
				var cuenta sql.NullInt32

				// Verificar si el proceso ya se corrió
				err = db.QueryRow("select version, archivo_salida from extractor.ext_procesados where id_modelo = $1 and fecha_desde = $2 and fecha_hasta = $3 order by version desc limit 1", proc.Id_modelo, restantes.Fecha1, restantes.Fecha2).Scan(&cuenta, &archivo_salida)
				if err != nil && err != sql.ErrNoRows {
					fmt.Println(err.Error())
					http.Error(w, "Error al escanear proceso", http.StatusBadRequest)
					return
				}
				if archivo_salida.Valid {
					archivoSalida = true
				}
				if cuenta.Valid {
					version = int(cuenta.Int32) + 1
				} else {
					version = 1
				}
				fmt.Println("Version: ", version)
				result, _ := src.ProcesadorSalida(proc, restantes.Fecha1, restantes.Fecha2, version, archivoSalida)
				if result != "" {
					resultado = append(resultado, result)
				}
				// El proceso termino, reinicio procesos
				procesos = nil
			}

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			respuesta := modelos.Respuesta{
				Mensaje:         "Datos recibidos y procesados",
				Archivos_salida: resultado,
			}
			jsonResp, _ := json.Marshal(respuesta)
			w.Write(jsonResp)
		}
	}
}

func control(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			procesos = nil
			var datos modelos.DTOdatos
			err := json.NewDecoder(r.Body).Decode(&datos)
			if err != nil {
				http.Error(w, "Error decodificando JSON de control", http.StatusBadRequest)
				return
			}
			datos.Fecha = src.FormatoFecha(datos.Fecha)
			datos.Fecha2 = src.FormatoFecha(datos.Fecha2)

			queryModelos := "SELECT em.id_modelo, em.id_empresa_adm, ea.razon_social as nombre_empresa, c.id_convenio as id_convenio, c.nombre as nombre_convenio, em.nombre, c.filtro as filtro_convenio, em.filtro_personas, em.filtro_recibos, em.formato_salida, em.archivo_modelo, em.filtro_having, em.archivo_control FROM extractor.ext_modelos em JOIN datos.empresas_adm ea ON em.id_empresa_adm = ea.id_empresa_adm JOIN extractor.ext_convenios c ON em.id_convenio = c.id_convenio where vigente and em.id_modelo = $1"
			// fmt.Println("Query modelos: ", queryModelos)
			stmt, err := db.Prepare(queryModelos)
			if err != nil {
				http.Error(w, "Error al preparar query", http.StatusBadRequest)
				return
			}
			defer stmt.Close()
			var args []interface{}
			args = append(args, datos.Id_modelo)
			rows, err := stmt.Query(args...)
			if err != nil {
				http.Error(w, "Error al ejecutar el query", http.StatusBadRequest)
				return
			}
			defer rows.Close()

			for rows.Next() {
				var proceso modelos.Proceso
				err = rows.Scan(&proceso.Id_modelo, &proceso.Id_empresa, &proceso.Nombre_empresa, &proceso.Id_convenio, &proceso.Nombre_convenio, &proceso.Nombre, &proceso.Filtro_convenio, &proceso.Filtro_personas, &proceso.Filtro_recibos, &proceso.Formato_salida, &proceso.Archivo_modelo, &proceso.Filtro_having, &proceso.Archivo_control)
				if err != nil {
					fmt.Println(err.Error())
					http.Error(w, "Error al escanear proceso", http.StatusBadRequest)
					return
				}
				proceso.Id_procesado = datos.Id_procesado
				procesos = append(procesos, proceso)
			}

			var version int
			var archivoControl bool
			if datos.Id_procesado > 0 {
				// Verificar si el proceso ya se corrió
				queryCuenta := fmt.Sprintf("select version, archivo_control from extractor.ext_procesados where id_proceso = %v", datos.Id_procesado)
				// fmt.Println("Query: ", queryCuenta)
				var archivo_control sql.NullString
				err = db.QueryRow(queryCuenta).Scan(&version, &archivo_control)
				if err != nil {
					fmt.Println(err.Error())
					http.Error(w, "Error al escanear proceso", http.StatusBadRequest)
					return
				}
				if archivo_control.Valid {
					archivoControl = true
				} else {
					archivoControl = false
				}

				version += 1
			} else {
				version = 1
			}

			var resultado []string
			result, errFormateado := src.ProcesadorControl(procesos[0], datos.Fecha, datos.Fecha2, version, archivoControl)
			if result != "" {
				resultado = append(resultado, result)
			}
			if errFormateado.Mensaje != "" {
				errString := "Error en " + procesos[0].Nombre + ": " + errFormateado.Mensaje
				// http.Error(w, errString, http.StatusBadRequest)
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				respuesta := modelos.Respuesta{
					Mensaje:          errString,
					Archivos_control: nil,
				}
				jsonResp, _ := json.Marshal(respuesta)
				w.Write(jsonResp)
				return

			}
			// El proceso termino, reinicio procesos
			procesos = nil

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			respuesta := modelos.Respuesta{
				Mensaje:          "Informe generado exitosamente",
				Archivos_control: resultado,
			}
			jsonResp, _ := json.Marshal(respuesta)
			w.Write(jsonResp)
		}
	}
}

func procesosRestantes(db *sql.DB) http.HandlerFunc {
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

		restantes = modelos.Restantes{
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

func getClientes(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientes = nil
		query := "select razon_social, cuit from datos.clientes"
		nombre_cliente := r.URL.Query().Get("cliente")
		cuit_cliente := r.URL.Query().Get("cuit")
		if len(nombre_cliente) > 0 {
			query += " where razon_social like '%" + nombre_cliente + "%'"
		}
		if len(cuit_cliente) > 0 {
			if len(nombre_cliente) == 0 {
				query += " where cuit like '%" + cuit_cliente + "%'"
			} else {
				query += " and cuit like '%" + cuit_cliente + "%'"
			}
		}

		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for rows.Next() {
			var cliente modelos.Cliente
			if err = rows.Scan(&cliente.Nombre, &cliente.Cuit); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			clientes = append(clientes, cliente)

		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(clientes); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func main() {

	db, err := conexiones.ConectarBase("postgres", "test", "postgres")
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	router := mux.NewRouter()
	// Carga de archivos estaticos
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("../client/static"))))
	router.PathPrefix("/salida/").Handler(http.StripPrefix("/salida/", http.FileServer(http.Dir("../salida"))))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../client/index.html")
	})

	router.HandleFunc("/convenios", getConvenios(db))
	router.HandleFunc("/empresas", getEmpresas(db))
	router.HandleFunc("/empresas/{id_convenio}", getEmpresas(db))
	router.HandleFunc("/conceptos", getConceptos(db))
	router.HandleFunc("/conceptos/{id_convenio}", getConceptos(db))
	router.HandleFunc("/conceptos/{id_empresa}", getConceptos(db))
	router.HandleFunc("/conceptos/{id_convenio}/{id_empresa}", getConceptos(db))
	router.HandleFunc("/procesos", getProcesos(db))
	router.HandleFunc("/send", sender(db)).Methods("POST")
	router.HandleFunc("/multiple", multipleSend(db)).Methods("POST")
	router.HandleFunc("/restantes", procesosRestantes(db))
	router.HandleFunc("/clientes", getClientes(db))
	router.HandleFunc("/control", control(db))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Listening...")
	srv.ListenAndServe()
}
