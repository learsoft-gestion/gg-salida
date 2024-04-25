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

func Sender(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			Procesos = nil
			var datos modelos.DTOdatos
			err := json.NewDecoder(r.Body).Decode(&datos)
			if err != nil {
				http.Error(w, "Error decodificando JSON", http.StatusBadRequest)
				return
			}
			datos.Fecha = src.FormatoFecha(datos.Fecha)
			datos.Fecha2 = src.FormatoFecha(datos.Fecha2)

			queryModelos := "SELECT em.id_modelo, em.id_empresa_adm, ea.razon_social as nombre_empresa, c.id_convenio as id_convenio, c.nombre as nombre_convenio, em.nombre, c.filtro as filtro_convenio, em.filtro_personas, em.filtro_recibos, em.formato_salida, em.archivo_modelo, em.filtro_having, em.archivo_nomina, em.columna_estado, em.id_query, em.select_control FROM extractor.ext_modelos em JOIN datos.empresas_adm ea ON em.id_empresa_adm = ea.id_empresa_adm JOIN extractor.ext_convenios c ON em.id_convenio = c.id_convenio where vigente and em.id_modelo = $1"
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
				var estado sql.NullString
				var select_control sql.NullString
				err = rows.Scan(&proceso.Id_modelo, &proceso.Id_empresa, &proceso.Nombre_empresa, &proceso.Id_convenio, &proceso.Nombre_convenio, &proceso.Nombre, &proceso.Filtro_convenio, &proceso.Filtro_personas, &proceso.Filtro_recibos, &proceso.Formato_salida, &proceso.Archivo_modelo, &proceso.Filtro_having, &proceso.Archivo_nomina, &estado, &proceso.Id_query, &select_control)
				if err != nil {
					fmt.Println(err.Error())
					http.Error(w, "Error al escanear proceso", http.StatusBadRequest)
					return
				}
				if estado.Valid {
					proceso.Columna_estado = estado.String
				} else {
					proceso.Columna_estado = ""
				}
				if select_control.Valid {
					proceso.Select_control = select_control.String
				} else {
					proceso.Select_control = ""
				}
				proceso.Id_procesado = datos.Id_procesado
				Procesos = append(Procesos, proceso)
			}

			version := 1
			var archivo_salida bool
			// Verificar si el proceso ya se corrió
			var archivoSalida sql.NullString
			var num_version sql.NullInt32
			err = db.QueryRow("select num_version, archivo_salida from extractor.ext_procesados where id_modelo = $1 and fecha_desde = $2 and fecha_hasta = $3 order by num_version desc limit 1;", datos.Id_modelo, datos.Fecha, datos.Fecha2).Scan(&num_version, &archivoSalida)
			if err != nil && err != sql.ErrNoRows {
				fmt.Println(err.Error())
				http.Error(w, "Error al escanear proceso", http.StatusBadRequest)
				return
			}
			if archivoSalida.Valid {
				archivo_salida = true
			}
			if num_version.Valid {
				version += int(num_version.Int32)
			}

			datos.Version = version

			var resultado []string
			result, id_procesado, errFormateado := src.ProcesadorSalida(Procesos[0], datos.Fecha, datos.Fecha2, version, archivo_salida)
			if result != "" {
				resultado = append(resultado, result)
			}
			datos.Id_procesado = id_procesado
			Procesos[0].Id_procesado = id_procesado
			if errFormateado.Mensaje != "" {
				errString := "Error en " + Procesos[0].Nombre + ": " + errFormateado.Mensaje
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

			// Ejecutar nomina
			respuesta_nomina := Nomina(datos)

			// Ejecutar control
			respuesta_control := Control(datos)

			if respuesta_nomina.Archivos_nomina != nil {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				respuesta := modelos.Respuesta{
					Mensaje:          "Informe generado exitosamente",
					Archivos_salida:  resultado,
					Archivos_nomina:  respuesta_nomina.Archivos_nomina,
					Archivos_control: respuesta_control.Archivos_control,
				}
				jsonResp, _ := json.Marshal(respuesta)
				w.Write(jsonResp)
			} else {
				w.WriteHeader(http.StatusBadRequest)
				jsonResp, _ := json.Marshal(respuesta_nomina)
				w.Write(jsonResp)
			}

		}

	}
}

func MultipleSend(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			Procesos = nil

			var placeholders []string
			for i := range Restantes.Id {
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
			for _, arg := range Restantes.Id {
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
				Procesos = append(Procesos, proceso)
			}

			datos := modelos.DTOdatos{
				Fecha:  Restantes.Fecha1,
				Fecha2: Restantes.Fecha2,
			}

			var resultado_salida []string
			var resultado_nomina []string
			for _, proc := range Procesos {
				var archivoSalida bool
				var archivo_salida sql.NullString
				var version int
				var cuenta sql.NullInt32

				// Verificar si el proceso ya se corrió
				err = db.QueryRow("select num_version, archivo_salida from extractor.ext_procesados where id_modelo = $1 and fecha_desde = $2 and fecha_hasta = $3 order by num_version desc limit 1", proc.Id_modelo, Restantes.Fecha1, Restantes.Fecha2).Scan(&cuenta, &archivo_salida)
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
				result_salida, id_procesado, err := src.ProcesadorSalida(proc, Restantes.Fecha1, Restantes.Fecha2, version, archivoSalida)
				if err.Mensaje != "" {
					fmt.Println(err.Mensaje)
					http.Error(w, err.Mensaje, http.StatusBadRequest)
					return
				}
				if result_salida != "" {
					resultado_salida = append(resultado_salida, result_salida)
				}
				datos.Id_modelo = proc.Id_modelo
				datos.Id_procesado = id_procesado
				datos.Version = version
				// El proceso termino, reinicio procesos
				// procesos = nil

				// Ejecutar nomina
				result_nomina := Nomina(datos)

				if result_nomina.Archivos_nomina != nil {
					resultado_nomina = append(resultado_nomina, result_nomina.Archivos_nomina[0])
				}
			}

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			respuesta := modelos.Respuesta{
				Mensaje:         "Datos recibidos y procesados",
				Archivos_salida: resultado_salida,
				Archivos_nomina: resultado_nomina,
			}
			jsonResp, _ := json.Marshal(respuesta)
			w.Write(jsonResp)
		}
	}
}

func Nomina(datos modelos.DTOdatos) modelos.Respuesta {

	var resultado []string
	result, errFormateado := src.ProcesadorNomina(Procesos[0], datos.Fecha, datos.Fecha2, datos.Version)
	if result != "" {
		resultado = append(resultado, result)
	}
	if errFormateado.Mensaje != "" {
		errString := "Error en " + Procesos[0].Nombre + ": " + errFormateado.Mensaje

		respuesta := modelos.Respuesta{
			Mensaje:         errString,
			Archivos_nomina: nil,
		}

		return respuesta
	}

	respuesta := modelos.Respuesta{
		Mensaje:         "Informe generado exitosamente",
		Archivos_nomina: resultado,
	}
	return respuesta
}

func Control(datos modelos.DTOdatos) modelos.Respuesta {

	var resultado []string
	result, errFormateado := src.ProcesadorControl(Procesos[0], datos.Fecha, datos.Fecha2, datos.Version)
	if result != "" {
		resultado = append(resultado, result)
	}
	if errFormateado.Mensaje != "" {
		errString := "Error en " + Procesos[0].Nombre + ": " + errFormateado.Mensaje

		respuesta := modelos.Respuesta{
			Mensaje:          errString,
			Archivos_control: nil,
		}

		return respuesta
	}

	respuesta := modelos.Respuesta{
		Mensaje:          "Informe generado exitosamente",
		Archivos_control: resultado,
	}
	return respuesta
}
