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
				http.Error(w, "Error decodificando JSON: "+err.Error(), http.StatusBadRequest)
				return
			}
			datos.Fecha = src.FormatoFecha(datos.Fecha)
			datos.Fecha2 = src.FormatoFecha(datos.Fecha2)
			fmt.Println("MODELO: ", datos.Id_modelo)
			queryModelos := "SELECT em.id_modelo, em.id_empresa_adm, ea.razon_social as nombre_empresa, ea.reducido as nombre_empresa_reducido, c.id_convenio as id_convenio, c.nombre as nombre_convenio, em.id_concepto, em.id_tipo, em.nombre, c.filtro as filtro_convenio, em.filtro_personas, em.filtro_recibos, em.formato_salida, em.archivo_modelo, em.archivo_nomina, em.columna_estado, em.id_query, em.select_control, em.select_salida FROM extractor.ext_modelos em JOIN extractor.ext_empresas_adm ea ON em.id_empresa_adm = ea.id_empresa_adm JOIN extractor.ext_convenios c ON em.id_convenio = c.id_convenio where vigente and em.id_modelo = $1"
			// fmt.Println("Query modelos: ", queryModelos)
			stmt, err := db.Prepare(queryModelos)
			if err != nil {
				fmt.Println(err.Error())
				http.Error(w, "Error al preparar query: "+err.Error(), http.StatusBadRequest)
				return
			}
			defer stmt.Close()
			var args []interface{}
			args = append(args, datos.Id_modelo)
			rows, err := stmt.Query(args...)
			if err != nil {
				http.Error(w, "Error al ejecutar el query: "+err.Error(), http.StatusBadRequest)
				return
			}
			defer rows.Close()

			for rows.Next() {
				var proceso modelos.Proceso
				var estado sql.NullString
				var select_control sql.NullString
				var filtro_recibos sql.NullString
				var filtro_personas sql.NullString

				err = rows.Scan(&proceso.Id_modelo, &proceso.Id_empresa, &proceso.Nombre_empresa, &proceso.Nombre_empresa_reducido, &proceso.Id_convenio, &proceso.Nombre_convenio, &proceso.Id_concepto, &proceso.Id_tipo, &proceso.Nombre, &proceso.Filtro_convenio, &filtro_personas, &filtro_recibos, &proceso.Formato_salida, &proceso.Archivo_modelo, &proceso.Archivo_nomina, &estado, &proceso.Id_query, &select_control, &proceso.Select_salida)
				if err != nil {
					fmt.Println(err.Error())
					http.Error(w, "Error al escanear proceso: "+err.Error(), http.StatusBadRequest)
					return
				}
				if filtro_personas.Valid {
					proceso.Filtro_personas = filtro_personas.String
				} else {
					proceso.Filtro_personas = ""
				}
				if filtro_recibos.Valid {
					proceso.Filtro_recibos = filtro_recibos.String
				} else {
					proceso.Filtro_recibos = ""
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
				http.Error(w, "Error al escanear proceso: "+err.Error(), http.StatusBadRequest)
				return
			}
			if archivoSalida.Valid {
				archivo_salida = true
			}
			if num_version.Valid {
				version += int(num_version.Int32)
			}

			datos.Version = version

			fmt.Printf("## Salida de %s ##\n", Procesos[0].Nombre)
			var resultado []string
			result, id_procesado, errFormateado, sql := src.ProcesadorSalida(db, Procesos[0], datos.Fecha, datos.Fecha2, version, archivo_salida)
			if errFormateado.Mensaje != "" {
				_, err = src.ProcesadosSalida(db, datos.Id_modelo, datos.Fecha, datos.Fecha2, version, 0, "Error")
				if err != nil {
					fmt.Println(err.Error())
					http.Error(w, "Error al loguear en procesados: "+err.Error(), http.StatusBadRequest)
					return
				}
				errString := "Error en " + Procesos[0].Nombre + ": " + errFormateado.Mensaje
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
			if result != "" {
				resultado = append(resultado, result)
			}
			datos.Id_procesado = id_procesado
			Procesos[0].Id_procesado = id_procesado

			// Ejecutar nomina
			respuesta_nomina := Nomina(db, sql, datos, Procesos[0])

			// Ejecutar control
			respuesta_control := Control(db, sql, datos, Procesos[0])

			if respuesta_nomina.Archivos_nomina[0] == "No se han encontrado registros" {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				respuesta := modelos.Respuesta{
					Mensaje:          "No se han encontrado registros",
					Archivos_salida:  resultado,
					Archivos_nomina:  respuesta_nomina.Archivos_nomina,
					Archivos_control: respuesta_control.Archivos_control,
				}
				jsonResp, _ := json.Marshal(respuesta)
				w.Write(jsonResp)
			} else if (respuesta_nomina.Archivos_nomina != nil) && (respuesta_control.Archivos_control != nil) {
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
			} else if respuesta_nomina.Archivos_nomina == nil {
				w.WriteHeader(http.StatusBadRequest)
				jsonResp, _ := json.Marshal(respuesta_nomina)
				w.Write(jsonResp)
			} else {
				w.WriteHeader(http.StatusBadRequest)
				jsonResp, _ := json.Marshal(respuesta_control)
				w.Write(jsonResp)
			}

		} else {
			http.Error(w, "Esta ruta solo admite una solicitud POST", http.StatusBadRequest)
			return
		}
	}
}

func MultipleSend(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// println("Restantes pre func", Restantes)
		if r.Method == "POST" {
			Procesos = nil

			var placeholders []string
			for i := range Restantes.Id {
				placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
			}
			// fmt.Println(Restantes)
			// fmt.Println(placeholders)
			if len(Restantes.Id) < 1 {
				fmt.Println("Restantes vacio")
				http.Error(w, "Restantes vacio", http.StatusBadRequest)
				return
			}

			queryModelos := fmt.Sprintf("SELECT em.id_modelo, em.id_empresa_adm, ea.razon_social as nombre_empresa, ea.reducido as nombre_empresa_reducido, c.id_convenio as id_convenio, c.nombre as nombre_convenio, em.id_concepto, em.id_tipo, em.nombre, c.filtro as filtro_convenio, em.filtro_personas, em.filtro_recibos, em.formato_salida, em.archivo_modelo, em.archivo_nomina, em.columna_estado, em.id_query, em.select_control, em.select_salida FROM extractor.ext_modelos em JOIN extractor.ext_empresas_adm ea ON em.id_empresa_adm = ea.id_empresa_adm JOIN extractor.ext_convenios c ON em.id_convenio = c.id_convenio where vigente and em.id_modelo in (%s)", strings.Join(placeholders, ","))

			stmt, err := db.Prepare(queryModelos)
			if err != nil {
				fmt.Println(len(Restantes.Id))
				fmt.Println("Error al preparar query: ", err.Error())
				// fmt.Printf("Error al preparar query: %s \nQuery: %s\n", err.Error(), queryModelos)
				strRes := fmt.Sprintf("Error al preparar query. Restantes: %v", len(Restantes.Id))
				http.Error(w, strRes, http.StatusInternalServerError)
				return
			}
			defer stmt.Close()
			var args []interface{}
			for _, arg := range Restantes.Id {
				args = append(args, arg)
			}

			rows, err := stmt.Query(args...)
			if err != nil {
				http.Error(w, "Error al ejecutar el query: "+err.Error(), http.StatusBadRequest)
				return
			}
			defer rows.Close()
			for rows.Next() {
				var proceso modelos.Proceso
				var estado sql.NullString
				var select_control sql.NullString
				var filtro_recibos sql.NullString
				var filtro_personas sql.NullString

				err = rows.Scan(&proceso.Id_modelo, &proceso.Id_empresa, &proceso.Nombre_empresa, &proceso.Nombre_empresa_reducido, &proceso.Id_convenio, &proceso.Nombre_convenio, &proceso.Id_concepto, &proceso.Id_tipo, &proceso.Nombre, &proceso.Filtro_convenio, &filtro_personas, &filtro_recibos, &proceso.Formato_salida, &proceso.Archivo_modelo, &proceso.Archivo_nomina, &estado, &proceso.Id_query, &select_control, &proceso.Select_salida)
				if err != nil {
					fmt.Println(err.Error())
					http.Error(w, "Error al escanear proceso: "+err.Error(), http.StatusBadRequest)
					return
				}
				if filtro_personas.Valid {
					proceso.Filtro_personas = filtro_personas.String
				} else {
					proceso.Filtro_personas = ""
				}
				if filtro_recibos.Valid {
					proceso.Filtro_recibos = filtro_recibos.String
				} else {
					proceso.Filtro_recibos = ""
				}
				if estado.Valid {
					proceso.Columna_estado = estado.String
				}
				if select_control.Valid {
					proceso.Select_control = select_control.String
				}
				Procesos = append(Procesos, proceso)
			}

			datos := modelos.DTOdatos{
				Fecha:  Restantes.Fecha1,
				Fecha2: Restantes.Fecha2,
			}

			var resultado_salida []string
			var resultado_nomina []string
			var resultado_control []string

			for _, proc := range Procesos {
				var archivoSalida bool
				var archivo_salida sql.NullString
				version := 1
				var cuenta sql.NullInt32

				// Verificar si el proceso ya se corrió
				err = db.QueryRow("select num_version, archivo_salida from extractor.ext_procesados where id_modelo = $1 and fecha_desde = $2 and fecha_hasta = $3 order by num_version desc limit 1", proc.Id_modelo, Restantes.Fecha1, Restantes.Fecha2).Scan(&cuenta, &archivo_salida)
				if err != nil && err != sql.ErrNoRows {
					fmt.Println(err.Error())
					http.Error(w, "Error al escanear proceso: "+err.Error(), http.StatusBadRequest)
					return
				}
				if archivo_salida.Valid {
					archivoSalida = true
				}
				if cuenta.Valid {
					version += int(cuenta.Int32) + 1
				}

				fmt.Printf("## Salida de %s ##\n", proc.Nombre)
				result_salida, id_procesado, errFormateado, sql := src.ProcesadorSalida(db, proc, Restantes.Fecha1, Restantes.Fecha2, version, archivoSalida)
				if errFormateado.Mensaje != "" {
					_, err := src.ProcesadosSalida(db, proc.Id_modelo, Restantes.Fecha1, Restantes.Fecha2, version, 0, "Error")
					if err != nil {
						fmt.Println(err.Error())
						http.Error(w, "Error al loguear en procesados: "+err.Error(), http.StatusBadRequest)
						return
					}

					fmt.Println(errFormateado.Mensaje)
					http.Error(w, errFormateado.Mensaje, http.StatusBadRequest)
					return
				}
				if result_salida != "" {
					resultado_salida = append(resultado_salida, result_salida)
				}
				datos.Id_modelo = proc.Id_modelo
				datos.Id_procesado = id_procesado
				proc.Id_procesado = id_procesado
				datos.Version = version

				// Ejecutar nomina
				result_nomina := Nomina(db, sql, datos, proc)

				// Ejecutar control
				result_control := Control(db, sql, datos, proc)

				if result_nomina.Archivos_nomina != nil {
					resultado_nomina = append(resultado_nomina, result_nomina.Archivos_nomina[0])
				} else {
					fmt.Println("Error en la nomina: ", result_nomina.Mensaje)
					http.Error(w, "Error interno del servidor", http.StatusBadRequest)
					return
				}
				if result_control.Archivos_control != nil {
					resultado_control = append(resultado_nomina, result_control.Archivos_control[0])
				} else {
					fmt.Println("Error en el control: ", result_control.Mensaje)
					http.Error(w, "Error interno del servidor", http.StatusBadRequest)
					return
				}
			}

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			respuesta := modelos.Respuesta{
				Mensaje:          "Datos recibidos y procesados",
				Archivos_salida:  resultado_salida,
				Archivos_nomina:  resultado_nomina,
				Archivos_control: resultado_control,
			}
			jsonResp, _ := json.Marshal(respuesta)
			w.Write(jsonResp)
		}
	}
}

func Nomina(db *sql.DB, sql *sql.DB, datos modelos.DTOdatos, proceso modelos.Proceso) modelos.Respuesta {
	fmt.Printf("## Nomina de %s ##\n", proceso.Nombre)
	var resultado []string
	result, errFormateado := src.ProcesadorNomina(db, sql, proceso, datos.Fecha, datos.Fecha2, datos.Version)
	if result != "" {
		resultado = append(resultado, result)
	}
	if errFormateado.Mensaje != "" {
		if err := src.ProcesadosNomina(db, proceso.Id_procesado, 0, "Error"); err != nil {
			errString := "Error al loguear en procesados: " + errFormateado.Mensaje

			respuesta := modelos.Respuesta{
				Mensaje:         errString,
				Archivos_nomina: nil,
			}

			return respuesta
		}

		errString := "Error en " + proceso.Nombre + ": " + errFormateado.Mensaje

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

func Control(db *sql.DB, sql *sql.DB, datos modelos.DTOdatos, proceso modelos.Proceso) modelos.Respuesta {
	fmt.Printf("## Control de %s ##\n", proceso.Nombre)

	var resultado []string
	result, errFormateado := src.ProcesadorControl(db, sql, proceso, datos.Fecha, datos.Fecha2, datos.Version)
	if result != "" {
		resultado = append(resultado, result)
	}
	if errFormateado.Mensaje != "" {

		if err := src.ProcesadosControl(db, proceso.Id_procesado, "Error"); err != nil {
			errString := "Error al loguear en procesados: " + errFormateado.Mensaje

			respuesta := modelos.Respuesta{
				Mensaje:          errString,
				Archivos_control: nil,
			}

			return respuesta
		}

		errString := "Error en " + proceso.Nombre + ": " + errFormateado.Mensaje

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

	fmt.Printf("------------------------------------------------\n")
	return respuesta
}
