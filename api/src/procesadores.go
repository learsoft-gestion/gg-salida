package src

import (
	"Nueva/conexiones"
	"Nueva/modelos"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

func ProcesadorSalida(db *sql.DB, proceso modelos.Proceso, fecha string, fecha2 string, version int, procesado_salida bool) (string, int, modelos.ErrorFormateado, *sql.DB) {

	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error: ", err.Error())
	}

	// Conexion al origen de datos
	sql, err := conexiones.ConectarBase("recibos", "prod", "sqlserver")
	if err != nil {
		return "", 0, modelos.ErrorFormateado{Mensaje: err.Error()}, nil
	}

	id_log, idLogDetalle, err := Logueo(db, proceso.Nombre)
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", 0, modelos.ErrorFormateado{Mensaje: err.Error()}, nil
	}

	// Reemplazar Alicuotas
	var select_salida modelos.Select_control
	err = db.QueryRow("SELECT extractor.obt_salida($1, $2, $3)", proceso.Id_modelo, fecha, fecha2).Scan(&select_salida.Query)
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", 0, modelos.ErrorFormateado{Mensaje: err.Error()}, nil
	}
	if select_salida.Query.Valid {
		proceso.Select_salida = select_salida.Query.String
	} else {
		proceso.Select_salida = ""
	}

	var query string
	var queryFinal string
	db.QueryRow("SELECT texto_query FROM extractor.ext_query where id_query = $1", proceso.Id_query).Scan(&query)
	// var queryReplace string
	// db.QueryRow("SELECT valor from extractor.ext_variables where variable = $1", proceso.Select_salida).Scan(&queryReplace)
	queryFinal = strings.Replace(query, "$SELECT$", proceso.Select_salida, 1)
	proceso.Query = queryFinal

	registros, err := Extractor(db, sql, proceso, fecha, fecha2, idLogDetalle, "salida")
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", 0, modelos.ErrorFormateado{Mensaje: err.Error()}, nil
	}

	fmt.Println("Cantidad de registros: ", len(registros))

	var name string

	if proceso.Archivo_modelo == "" || fecha != fecha2 {
		// No genera salida

		// Insertar nuevo proceso en ext_procesados
		if idProc, err := ProcesadosSalida(db, proceso.Id_modelo, fecha, fecha2, version, len(registros), ""); err != nil {
			ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return "", 0, modelos.ErrorFormateado{Mensaje: err.Error()}, nil
		} else if idProc > 0 {
			proceso.Id_procesado = idProc
		}

		// Logueo
		_, err = db.Exec("CALL extractor.act_log_detalle($1, 'F', $2)", idLogDetalle, "Modelo no genera salida")
		if err != nil {
			ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return "", 0, modelos.ErrorFormateado{Mensaje: err.Error()}, nil
		}
	} else if len(registros) == 0 {
		// No se han encontrado registros

		// Insertar nuevo proceso en ext_procesados
		if idProc, err := ProcesadosSalida(db, proceso.Id_modelo, fecha, fecha2, version, len(registros), "Sin datos"); err != nil {
			ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return "", 0, modelos.ErrorFormateado{Mensaje: err.Error()}, nil
		} else if idProc > 0 {
			proceso.Id_procesado = idProc
		}

		// Logueo
		_, err = db.Exec("CALL extractor.act_log_detalle($1, 'F', $2)", idLogDetalle, "No se han encontrado registros")
		if err != nil {
			ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return "", 0, modelos.ErrorFormateado{Mensaje: err.Error()}, nil
		}

		_, err = db.Exec("CALL extractor.etl_ending($1)", id_log)
		if err != nil {
			ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return "", 0, modelos.ErrorFormateado{Mensaje: err.Error()}, nil
		}

		return "No se han encontrado registros", proceso.Id_procesado, modelos.ErrorFormateado{Mensaje: ""}, sql
	} else {
		// Genera salida

		// Fecha para el nombre de salida
		var fechaSalida string
		if fecha == fecha2 {
			fechaSalida = fecha
		} else {
			fechaSalida = fecha + "-" + fecha2
		}

		// Directorio del archivo main.go
		directorioActual, err := os.Getwd()
		if err != nil {
			fmt.Println("Error al obtener el directorio actual:", err)
			ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return "", 0, modelos.ErrorFormateado{Mensaje: err.Error()}, nil
		}

		var nombreSalida string
		proceso_periodo := fecha + "-" + fecha2
		// Construir la ruta de la carpeta de salida
		rutaCarpeta := filepath.Join(directorioActual, ".", "salida", proceso.Nombre_empresa_reducido, proceso.Nombre_convenio, proceso_periodo, proceso.Nombre)

		// Verificar si la carpeta de salida existe, si no, crearla
		if _, err := os.Stat(rutaCarpeta); os.IsNotExist(err) {
			if err := os.MkdirAll(rutaCarpeta, 0755); err != nil {
				fmt.Println("Error al crear la carpeta de salida:", err)
				ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
				return "", 0, modelos.ErrorFormateado{Mensaje: err.Error()}, nil
			}
		}

		if version > 1 {
			nombreSalida = fmt.Sprintf("%s_%s_%s%s_%s_%s(%v)", proceso.Nombre_convenio, proceso.Nombre_empresa_reducido, proceso.Id_concepto, proceso.Id_tipo, proceso.Nombre, fechaSalida, version)
		} else {
			nombreSalida = fmt.Sprintf("%s_%s_%s%s_%s_%s", proceso.Nombre_convenio, proceso.Nombre_empresa_reducido, proceso.Id_concepto, proceso.Id_tipo, proceso.Nombre, fechaSalida)
		}

		// Formato del archivo de salida
		formato := strings.ToLower(proceso.Formato_salida)

		if formato == "xls" {
			// Ruta completa del archivo
			nombreSalida += ".xlsx"
			rutaArchivo := filepath.Join(rutaCarpeta, nombreSalida)
			plantilla := "./templates/" + proceso.Archivo_modelo

			name, err = CargarExcel(db, idLogDetalle, proceso, registros, rutaArchivo, plantilla, "salida", "", version)
			if err != nil {
				ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
				return "", 0, modelos.ErrorFormateado{Mensaje: err.Error()}, nil
			}

		} else if formato == "txt" {
			// Ruta completa del archivo
			nombreSalida += ".txt"
			rutaArchivo := filepath.Join(rutaCarpeta, nombreSalida)

			// Utilizar funcion para txt
			name, err = CargarTxt(db, idLogDetalle, proceso, registros, rutaArchivo)
			if err != nil {
				ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
				return "", 0, modelos.ErrorFormateado{Mensaje: err.Error()}, nil
			}
		} else if formato == "xml" {
			// Ruta completa del archivo
			nombreSalida += ".xml"
			rutaArchivo := filepath.Join(rutaCarpeta, nombreSalida)
			// Utilizar funcion para txt
			name, err = CargarXml(db, idLogDetalle, proceso, registros, rutaArchivo)
			if err != nil {
				ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
				return "", 0, modelos.ErrorFormateado{Mensaje: err.Error()}, nil
			}
		}

		// Insertar nuevo proceso en ext_procesados
		if idProc, err := ProcesadosSalida(db, proceso.Id_modelo, fecha, fecha2, version, len(registros), filepath.Join(rutaCarpeta, nombreSalida)); err != nil {
			ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return "", 0, modelos.ErrorFormateado{Mensaje: err.Error()}, nil
		} else if idProc > 0 {
			proceso.Id_procesado = idProc
		}

		// Logueo
		_, err = db.Exec("CALL extractor.act_log_detalle($1, 'F', $2)", idLogDetalle, fmt.Sprintf("Archivo guardado en: \"%s\"", filepath.Join(rutaCarpeta, nombreSalida)))
		if err != nil {
			ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return "", 0, modelos.ErrorFormateado{Mensaje: err.Error()}, nil
		}
	}

	_, err = db.Exec("CALL extractor.etl_ending($1)", id_log)
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", 0, modelos.ErrorFormateado{Mensaje: err.Error()}, nil
	}

	return name, proceso.Id_procesado, modelos.ErrorFormateado{Mensaje: ""}, sql
}

func ProcesadorNomina(db *sql.DB, sql *sql.DB, proceso modelos.Proceso, fecha string, fecha2 string, version int, generarNomina bool) (string, modelos.ErrorFormateado) {

	id_log, idLogDetalle, err := Logueo(db, proceso.Nombre)
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", modelos.ErrorFormateado{Mensaje: err.Error()}
	}

	var name string

	var query string
	var queryFinal string
	db.QueryRow("SELECT texto_query FROM extractor.ext_query where id_query = $1", proceso.Id_query).Scan(&query)
	var queryReplace string
	db.QueryRow("SELECT valor from extractor.ext_variables where variable = 'SELECT'").Scan(&queryReplace)
	queryFinal = strings.Replace(query, "$SELECT$", queryReplace, 1)
	proceso.Query = queryFinal

	var registros []modelos.Registro
	if generarNomina {
		registros, err = Extractor(db, sql, proceso, fecha, fecha2, idLogDetalle, "nomina")
		if err != nil {
			ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return err.Error(), modelos.ErrorFormateado{Mensaje: err.Error()}
		}
	}

	if len(registros) == 0 && generarNomina {
		if err = ProcesadosNomina(db, proceso.Id_procesado, 0, "Sin datos"); err != nil {
			fmt.Println(err.Error())
			return "", modelos.ErrorFormateado{Mensaje: err.Error()}
		}

		// Logueo
		_, err = db.Exec("CALL extractor.act_log_detalle($1, 'F', $2)", idLogDetalle, "No se han encontrado registros")
		if err != nil {
			ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return "", modelos.ErrorFormateado{Mensaje: err.Error()}
		}

		_, err = db.Exec("CALL extractor.etl_ending($1)", id_log)
		if err != nil {
			ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return "", modelos.ErrorFormateado{Mensaje: err.Error()}
		}

		return "Sin datos", modelos.ErrorFormateado{Mensaje: ""}

	} else {
		if generarNomina {
			fmt.Println("Cantidad de registros: ", len(registros))
		}

		if fecha == fecha2 {
			if generarNomina {
				// Construir la parte de las columnas y los placeholders para los valores
				columnas_slice := []string{"id_modelo", "fecha", "num_version"}
				placeholders := make([]string, len(registros[0].Columnas))
				for i, col := range registros[0].Columnas {
					columnas_slice = append(columnas_slice, strings.ToLower(col))
					placeholders[i] = fmt.Sprintf("$%d", i+4)
				}
				columnas := strings.Join(columnas_slice, ", ")
				placeholdersFinales := []string{"$1", "$2", "$3"}
				placeholdersFinales = append(placeholdersFinales, placeholders...)
				placeholdersStr := strings.Join(placeholdersFinales, ", ")

				var registrosInsert [][]interface{}
				for _, reg := range registros {
					valoresFinales := []interface{}{proceso.Id_modelo, fecha, version}
					for _, columna := range reg.Columnas {
						valoresFinales = append(valoresFinales, reg.Valores[strings.ToUpper(columna)])
					}
					registrosInsert = append(registrosInsert, valoresFinales)
				}

				// fmt.Printf("Columnas: \n%s\nPlacehoders:\n%s", columnas, placeholdersStr)

				// Construir la consulta SQL
				query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", "extractor.ext_nomina_congelada", columnas, placeholdersStr)

				// Ejecutar la consulta
				err = MultipleInsertSQL(db, registrosInsert, query)
				if err != nil {
					fmt.Println("Error al obtener el directorio actual:", err)
					ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
					return "", modelos.ErrorFormateado{Mensaje: err.Error()}
				}
			}

			// Obtener control congelado
			filas_congelados, err := db.Query("select * from extractor.ext_nomina_congelada	where (fecha, id_modelo, num_version) in ( select fecha, id_modelo, max(num_version) as version	from extractor.ext_nomina_congelada where ((substring(fecha,1,4) = substring($1,1,4) and fecha < $1) or (substring($1,5,6) = '01' and substring(fecha,1,4) = cast(cast(substring($1,1,4) as int) - 1 as varchar) and substring(fecha,5,6) = '12' )) and id_modelo = $2 group by fecha, id_modelo) order by fecha;", fecha, proceso.Id_modelo)
			if err != nil {
				ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
				return "", modelos.ErrorFormateado{Mensaje: err.Error()}
			}

			columnas_congeladas, err := filas_congelados.Columns()
			if err != nil {
				fmt.Println("Error al hacer la query del extractor")
				ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
				return "", modelos.ErrorFormateado{Mensaje: err.Error()}
			}

			val_congelados := make([]interface{}, len(columnas_congeladas))
			for i := range val_congelados {
				val_congelados[i] = new(interface{})
			}
			var congelados_nomina []modelos.Registro
			for filas_congelados.Next() {

				if err := filas_congelados.Scan(val_congelados...); err != nil {
					ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
					return "", modelos.ErrorFormateado{Mensaje: err.Error()}
				}

				registroMapa := make(map[string]interface{})
				for i, colNombre := range columnas_congeladas {
					colName := strings.ToUpper(colNombre)
					registroMapa[colName] = *(val_congelados[i].(*interface{}))
				}

				registro := modelos.Registro{
					Columnas: columnas_congeladas,
					Valores:  registroMapa,
				}
				congelados_nomina = append(congelados_nomina, registro)

			}

			nuevoSlice := make([]modelos.Registro, len(congelados_nomina)+len(registros))

			// Copiar el slice a agregar al principio en el nuevo slice
			copy(nuevoSlice, congelados_nomina)

			// Copiar el slice original al final del nuevo slice
			copy(nuevoSlice[len(congelados_nomina):], registros)

			// Asignar el nuevo slice a la variable original
			registros = nuevoSlice

		}

		// Fecha para el nombre de salida
		var fechaControl string
		if fecha == fecha2 {
			fechaControl = fecha
		} else {
			fechaControl = fecha + "-" + fecha2
		}

		// Directorio del archivo main.go
		directorioActual, err := os.Getwd()
		if err != nil {
			fmt.Println("Error al obtener el directorio actual:", err)
			ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return "", modelos.ErrorFormateado{Mensaje: err.Error()}
		}

		// Construir la ruta de la carpeta de salida
		var nombreControl string
		proceso_periodo := fecha + "-" + fecha2
		procesoNombre := proceso.Nombre + "-Nomina"
		rutaCarpeta := filepath.Join(directorioActual, ".", "salida", proceso.Nombre_empresa_reducido, proceso.Nombre_convenio, proceso_periodo, procesoNombre)

		// Verificar si la carpeta de salida existe, si no, crearla
		if _, err := os.Stat(rutaCarpeta); os.IsNotExist(err) {
			if err := os.MkdirAll(rutaCarpeta, 0755); err != nil {
				fmt.Println("Error al crear la carpeta de salida: ", err)
				ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
				return "", modelos.ErrorFormateado{Mensaje: "Error al crear la carpeta de salida: " + err.Error()}
			}
		}

		if !generarNomina {
			nombreControl = "Consulta-"
		}

		if version > 1 {
			nombreControl += fmt.Sprintf("%s_%s_%s%s_%s_%s(%v)", proceso.Nombre_convenio, proceso.Nombre_empresa_reducido, proceso.Id_concepto, proceso.Id_tipo, procesoNombre, fechaControl, version)
		} else {
			nombreControl += fmt.Sprintf("%s_%s_%s%s_%s_%s", proceso.Nombre_convenio, proceso.Nombre_empresa_reducido, proceso.Id_concepto, proceso.Id_tipo, procesoNombre, fechaControl)
		}

		// Ruta completa del archivo
		nombreControl += ".xlsx"
		rutaArchivo := filepath.Join(rutaCarpeta, nombreControl)
		plantilla := "./templates/" + proceso.Archivo_nomina

		name, err = CargarExcel(db, idLogDetalle, proceso, registros, rutaArchivo, plantilla, "nomina", "", version)
		if err != nil {
			ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return "", modelos.ErrorFormateado{Mensaje: err.Error()}
		}

		if generarNomina {
			// Insertar o actualizar proceso en ext_procesados
			if err = ProcesadosNomina(db, proceso.Id_procesado, len(registros), filepath.Join(rutaCarpeta, nombreControl)); err != nil {
				ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
				return "", modelos.ErrorFormateado{Mensaje: err.Error()}
			}
		}

		// Logueo
		_, err = db.Exec("CALL extractor.act_log_detalle($1, 'F', $2)", idLogDetalle, fmt.Sprintf("Archivo guardado en: \"%s\"", filepath.Join(rutaCarpeta, nombreControl)))
		if err != nil {
			ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return "", modelos.ErrorFormateado{Mensaje: err.Error()}
		}
	}
	_, err = db.Exec("CALL extractor.etl_ending($1)", id_log)
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", modelos.ErrorFormateado{Mensaje: err.Error()}
	}

	return name, modelos.ErrorFormateado{Mensaje: ""}
}

func ProcesadorControl(db *sql.DB, sql *sql.DB, proceso modelos.Proceso, fecha string, fecha2 string, version int, generarControl bool) (string, modelos.ErrorFormateado) {

	id_log, idLogDetalle, err := Logueo(db, proceso.Nombre)
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", modelos.ErrorFormateado{Mensaje: err.Error()}
	}

	// Formato del archivo de salida
	var name string

	// Reemplazar Alicuotas
	var select_control modelos.Select_control
	err = db.QueryRow("SELECT extractor.obt_control($1, $2, $3)", proceso.Id_modelo, fecha, fecha2).Scan(&select_control.Query)
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", modelos.ErrorFormateado{Mensaje: err.Error()}
	}
	if select_control.Query.Valid {
		proceso.Select_control = select_control.Query.String
	} else {
		proceso.Select_control = ""
	}

	// fmt.Printf("Select_control despues del replace: \n%s\n", proceso.Select_control)

	var query string
	var queryFinal string
	db.QueryRow("SELECT texto_query FROM extractor.ext_query where id_query = $1", proceso.Id_query).Scan(&query)
	if proceso.Select_control != "" {
		queryFinal = strings.Replace(query, "$SELECT$", proceso.Select_control, 1)
	} else {
		// Logueo
		if err = ProcesadosControl(db, proceso.Id_procesado, "Sin datos"); err != nil {
			fmt.Println(err.Error())
			return err.Error(), modelos.ErrorFormateado{Mensaje: "error al loguear en procesados"}
		}

		_, err = db.Exec("CALL extractor.act_log_detalle($1, 'F', $2)", idLogDetalle, "el control no esta configurado para este modelo")
		if err != nil {
			ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return "", modelos.ErrorFormateado{Mensaje: err.Error()}
		}
		return "Sin datos", modelos.ErrorFormateado{Mensaje: ""}
	}
	proceso.Query = queryFinal

	var registros []modelos.Registro
	if generarControl {
		registros, err = Extractor(db, sql, proceso, fecha, fecha2, idLogDetalle, "control")
		if err != nil {
			ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return err.Error(), modelos.ErrorFormateado{Mensaje: err.Error()}
		}
	}

	if len(registros) == 0 && generarControl {
		if err = ProcesadosControl(db, proceso.Id_procesado, "Sin datos"); err != nil {
			fmt.Println(err.Error())
			return "", modelos.ErrorFormateado{Mensaje: err.Error()}
		}

		// Logueo
		_, err = db.Exec("CALL extractor.act_log_detalle($1, 'F', $2)", idLogDetalle, "No se han encontrado registros")
		if err != nil {
			ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return "", modelos.ErrorFormateado{Mensaje: err.Error()}
		}

		_, err = db.Exec("CALL extractor.etl_ending($1)", id_log)
		if err != nil {
			ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return "", modelos.ErrorFormateado{Mensaje: err.Error()}
		}

		return "Sin datos", modelos.ErrorFormateado{Mensaje: ""}

	} else {
		if generarControl {
			fmt.Println("Cantidad de registros: ", len(registros))
		}

		// Fecha para el nombre de salida
		var fechaControl string
		if fecha == fecha2 {
			fechaControl = fecha
		} else {
			fechaControl = fecha + "-" + fecha2
		}

		// Directorio del archivo main.go
		directorioActual, err := os.Getwd()
		if err != nil {
			fmt.Println("Error al obtener el directorio actual:", err)
			ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return "", modelos.ErrorFormateado{Mensaje: err.Error()}
		}

		var nombreControl string
		proceso_periodo := fecha + "-" + fecha2
		// Construir la ruta de la carpeta de salida
		procesoNombre := proceso.Nombre + "-Control"
		rutaCarpeta := filepath.Join(directorioActual, ".", "salida", proceso.Nombre_empresa_reducido, proceso.Nombre_convenio, proceso_periodo, procesoNombre)

		// Verificar si la carpeta de salida existe, si no, crearla
		if _, err := os.Stat(rutaCarpeta); os.IsNotExist(err) {
			if err := os.MkdirAll(rutaCarpeta, 0755); err != nil {
				fmt.Println("Error al crear la carpeta de salida:", err)
				ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
				return "", modelos.ErrorFormateado{Mensaje: err.Error()}
			}
		}

		if !generarControl {
			nombreControl = "Consulta-"
		}

		if version > 1 {
			nombreControl += fmt.Sprintf("%s_%s_%s%s_%s_%s(%v)", proceso.Nombre_convenio, proceso.Nombre_empresa_reducido, proceso.Id_concepto, proceso.Id_tipo, procesoNombre, fechaControl, version)
		} else {
			nombreControl += fmt.Sprintf("%s_%s_%s%s_%s_%s", proceso.Nombre_convenio, proceso.Nombre_empresa_reducido, proceso.Id_concepto, proceso.Id_tipo, procesoNombre, fechaControl)
		}

		// Ruta completa del archivo
		nombreControl += ".xlsx"
		rutaArchivo := filepath.Join(rutaCarpeta, nombreControl)
		// plantilla := "./templates/" + proceso.Archivo_control

		// Obtener descripcion para la solapa info del control
		infoText, err := getSQLResult(db, proceso.Id_modelo, fecha, fecha2)
		if err != nil {
			ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return "", modelos.ErrorFormateado{Mensaje: err.Error()}
		}

		if fecha == fecha2 {
			if generarControl {
				// for key, value := range registros[0].Valores {
				// 	fmt.Printf("Key: %s, Type: %s, Value: %v\n", key, reflect.TypeOf(value), value)
				// }

				// Convertir []uint8 en float64
				valores, err := ConvertirBytesAFloat64(registros[0].Valores)
				if err != nil {
					ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
					return "", modelos.ErrorFormateado{Mensaje: err.Error()}
				}

				// Convertir registros en un json
				jsonData, err := json.Marshal(modelos.Registro{Columnas: registros[0].Columnas, Valores: valores})
				if err != nil {
					ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
					return "", modelos.ErrorFormateado{Mensaje: err.Error()}
				}

				queryControlCongelado := "INSERT INTO extractor.ext_control_congelado (id_modelo, fecha, num_version, json_control) VALUES ($1,$2,$3,$4)"
				_, err = db.Exec(queryControlCongelado, proceso.Id_modelo, fecha, version, jsonData)
				if err != nil {
					ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
					return "", modelos.ErrorFormateado{Mensaje: err.Error()}
				}
			}

			// Obtener control congelado
			filas_congelados, err := db.Query("select * from extractor.ext_control_congelado where (fecha, id_modelo, num_version) in ( select fecha, id_modelo, max(num_version) as version from extractor.ext_control_congelado where ((substring(fecha,1,4) = substring($1,1,4) and fecha < $1) or (substring($1,5,6) = '01' and substring(fecha,1,4) = cast(cast(substring($1,1,4) as int) - 1 as varchar) and substring(fecha,5,6) = '12' )) and id_modelo = $2 group by fecha, id_modelo) order by fecha;", fecha, proceso.Id_modelo)
			if err != nil {
				ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
				return "", modelos.ErrorFormateado{Mensaje: err.Error()}
			}

			var congelados []modelos.Control_congelado
			for filas_congelados.Next() {
				var congelado modelos.Control_congelado
				err = filas_congelados.Scan(&congelado.Id_modelo, &congelado.Fecha, &congelado.Num_version, &congelado.Json_control)
				if err != nil {
					ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
					return "", modelos.ErrorFormateado{Mensaje: err.Error()}
				}
				congelados = append(congelados, congelado)
			}

			var controles []modelos.Registro
			for _, congelado := range congelados {
				var reg modelos.Registro
				// var control_congelado map[string]interface{}
				if err = json.Unmarshal(congelado.Json_control, &reg); err != nil {
					ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
					return "", modelos.ErrorFormateado{Mensaje: err.Error()}
				}
				// for key := range control_congelado {
				// 	reg.Columnas = append(reg.Columnas, key)
				// }
				// reg.Valores = control_congelado
				reg.Valores["NUM_VERSION"] = congelado.Num_version
				controles = append(controles, reg)
			}

			nuevoSlice := make([]modelos.Registro, len(controles)+len(registros))

			// Copiar el slice a agregar al principio en el nuevo slice
			copy(nuevoSlice, controles)

			// Copiar el slice original al final del nuevo slice
			copy(nuevoSlice[len(controles):], registros)

			// Asignar el nuevo slice a la variable original
			registros = nuevoSlice

		}

		name, err = CargarExcel(db, idLogDetalle, proceso, registros, rutaArchivo, "plantilla", "control", infoText, version)
		if err != nil {
			ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return "", modelos.ErrorFormateado{Mensaje: err.Error()}
		}

		if generarControl {
			// Insertar o actualizar proceso en ext_procesados
			if err = ProcesadosControl(db, proceso.Id_procesado, filepath.Join(rutaCarpeta, nombreControl)); err != nil {
				ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
				return "", modelos.ErrorFormateado{Mensaje: err.Error()}
			}
		}

		// Logueo
		_, err = db.Exec("CALL extractor.act_log_detalle($1, 'F', $2)", idLogDetalle, fmt.Sprintf("Archivo guardado en: \"%s\"", filepath.Join(rutaCarpeta, nombreControl)))
		if err != nil {
			ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return "", modelos.ErrorFormateado{Mensaje: err.Error()}
		}
	}
	_, err = db.Exec("CALL extractor.etl_ending($1)", id_log)
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", modelos.ErrorFormateado{Mensaje: err.Error()}
	}

	return name, modelos.ErrorFormateado{Mensaje: ""}
}

// Función para ejecutar la consulta SQL y obtener el resultado como texto
func getSQLResult(db *sql.DB, idModelo int, fecha1, fecha2 string) (string, error) {

	var result string
	query := "SELECT extractor.describir_filtros($1, $2, $3)"
	err := db.QueryRow(query, idModelo, fecha1, fecha2).Scan(&result)
	if err != nil {
		return "", err
	}

	return result, nil
}
