package src

import (
	"Nueva/conexiones"
	"Nueva/modelos"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ProcesadorSalida(proceso modelos.Proceso, fecha string, fecha2 string, version int) (string, modelos.ErrorFormateado) {

	db, err := conexiones.ConectarBase("postgres", "test", "postgres")
	if err != nil {
		return "", modelos.ErrorFormateado{Mensaje: err.Error()}
	}
	defer db.Close()

	// Conexion al origen de datos
	sql, err := conexiones.ConectarBase("recibos", "prod", "sqlserver")
	if err != nil {
		return "", modelos.ErrorFormateado{Mensaje: err.Error()}
	}
	defer sql.Close()

	id_log, idLogDetalle, err := Logueo(db, proceso.Nombre)
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", modelos.ErrorFormateado{Mensaje: err.Error()}
	}

	query := ""
	db.QueryRow("SELECT texto_query FROM extractor.ext_query;").Scan(&query)
	proceso.Query = query

	registros, err := Extractor(db, sql, proceso, fecha, fecha2, idLogDetalle)
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return err.Error(), modelos.ErrorFormateado{Mensaje: err.Error()}
	}

	if len(registros) == 0 {
		if err = ProcesadosSalida(db, proceso.Id_modelo, fecha, fecha2, version, 0, ""); err != nil {
			fmt.Println(err.Error())
			return err.Error(), modelos.ErrorFormateado{Mensaje: "error al loguear en procesados"}
		}
		return "", modelos.ErrorFormateado{Mensaje: "no se han encontrado registros"}
	} else {
		fmt.Println("Cantidad de registros: ", len(registros))
	}

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
		return "", modelos.ErrorFormateado{Mensaje: err.Error()}
	}

	var nombreSalida string
	proceso_periodo := fecha + "-" + fecha2
	// Construir la ruta de la carpeta de salida
	rutaCarpeta := filepath.Join(directorioActual, "..", "salida", proceso.Nombre_empresa, proceso.Nombre_convenio, proceso_periodo, proceso.Nombre)

	// Verificar si la carpeta de salida existe, si no, crearla
	if _, err := os.Stat(rutaCarpeta); os.IsNotExist(err) {
		if err := os.MkdirAll(rutaCarpeta, 0755); err != nil {
			fmt.Println("Error al crear la carpeta de salida:", err)
			ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return "", modelos.ErrorFormateado{Mensaje: err.Error()}
		}
	}

	if version > 1 {
		nombreSalida = fmt.Sprintf("%s_%s(%v)", proceso.Nombre, fechaSalida, version)
	} else {
		nombreSalida = fmt.Sprintf("%s_%s", proceso.Nombre, fechaSalida)
	}

	// Formato del archivo de salida
	formato := strings.ToLower(proceso.Formato_salida)
	var name string
	if formato == "xls" {
		// Ruta completa del archivo
		nombreSalida += ".xlsx"
		rutaArchivo := filepath.Join(rutaCarpeta, nombreSalida)

		name, err = CargarExcel(db, idLogDetalle, proceso, registros, rutaArchivo)
		if err != nil {
			ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return "", modelos.ErrorFormateado{Mensaje: err.Error()}
		}

	} else if formato == "txt" {
		// Ruta completa del archivo
		nombreSalida += ".txt"
		rutaArchivo := filepath.Join(rutaCarpeta, nombreSalida)

		// Utilizar funcion para txt
		name, err = CargarTxt(db, idLogDetalle, proceso, registros, rutaArchivo)
		if err != nil {
			ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return "", modelos.ErrorFormateado{Mensaje: err.Error()}
		}
	} else if formato == "xml" {
		// Ruta completa del archivo
		nombreSalida += ".xml"
		rutaArchivo := filepath.Join(rutaCarpeta, nombreSalida)
		// Utilizar funcion para txt
		name, err = CargarXml(db, idLogDetalle, proceso, registros, rutaArchivo)
		if err != nil {
			ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return "", modelos.ErrorFormateado{Mensaje: err.Error()}
		}
	}

	// Insertar nuevo proceso en ext_procesados
	if err = ProcesadosSalida(db, proceso.Id_modelo, fecha, fecha2, version, len(registros), filepath.Join(rutaCarpeta, nombreSalida)); err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", modelos.ErrorFormateado{Mensaje: err.Error()}
	}

	// Logueo
	_, err = db.Exec("CALL extractor.act_log_detalle($1, 'F', $2)", idLogDetalle, fmt.Sprintf("Archivo guardado en: \"%s\"", filepath.Join(rutaCarpeta, nombreSalida)))
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", modelos.ErrorFormateado{Mensaje: err.Error()}
	}
	_, err = db.Exec("CALL extractor.etl_ending($1)", id_log)
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", modelos.ErrorFormateado{Mensaje: err.Error()}
	}

	return name, modelos.ErrorFormateado{Mensaje: ""}
}

func ProcesadorControl(proceso modelos.Proceso, fecha string, fecha2 string, version int, procesado bool) (string, modelos.ErrorFormateado) {

	db, err := conexiones.ConectarBase("postgres", "test", "postgres")
	if err != nil {
		return "", modelos.ErrorFormateado{Mensaje: err.Error()}
	}
	defer db.Close()

	// Conexion al origen de datos
	sql, err := conexiones.ConectarBase("recibos", "prod", "sqlserver")
	if err != nil {
		return "", modelos.ErrorFormateado{Mensaje: err.Error()}
	}
	defer sql.Close()

	id_log, idLogDetalle, err := Logueo(db, proceso.Nombre)
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", modelos.ErrorFormateado{Mensaje: err.Error()}
	}

	query := ""
	db.QueryRow("SELECT texto_query FROM extractor.ext_query;").Scan(&query)
	proceso.Query = query

	registros, err := Extractor(db, sql, proceso, fecha, fecha2, idLogDetalle)
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return err.Error(), modelos.ErrorFormateado{Mensaje: err.Error()}
	}

	if len(registros) == 0 {
		if err = ProcesadosControl(db, proceso.Id_procesado, proceso.Id_modelo, fecha, fecha2, version, 0, "", procesado); err != nil {
			fmt.Println(err.Error())
			return err.Error(), modelos.ErrorFormateado{Mensaje: "error al loguear en procesados"}
		}
		return "", modelos.ErrorFormateado{Mensaje: "no se han encontrado registros"}
	} else {
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
	procesoNombre := proceso.Nombre + "-Nomina"
	rutaCarpeta := filepath.Join(directorioActual, "..", "salida", proceso.Nombre_empresa, proceso.Nombre_convenio, proceso_periodo, procesoNombre)

	// Verificar si la carpeta de salida existe, si no, crearla
	if _, err := os.Stat(rutaCarpeta); os.IsNotExist(err) {
		if err := os.MkdirAll(rutaCarpeta, 0755); err != nil {
			fmt.Println("Error al crear la carpeta de salida:", err)
			ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return "", modelos.ErrorFormateado{Mensaje: err.Error()}
		}
	}

	if version > 1 {
		nombreControl = fmt.Sprintf("%s_%s(%v)", procesoNombre, fechaControl, version)
	} else {
		nombreControl = fmt.Sprintf("%s_%s", procesoNombre, fechaControl)
	}

	// Formato del archivo de salida
	var name string
	// Ruta completa del archivo
	nombreControl += ".xlsx"
	rutaArchivo := filepath.Join(rutaCarpeta, nombreControl)
	name, err = CargarExcel(db, idLogDetalle, proceso, registros, rutaArchivo)
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", modelos.ErrorFormateado{Mensaje: err.Error()}
	}

	// } else if formato == "txt" {
	// 	// Ruta completa del archivo
	// 	nombreControl += ".txt"
	// 	rutaArchivo := filepath.Join(rutaCarpeta, nombreControl)
	// 	// Utilizar funcion para txt
	// 	name, err = CargarTxt(db, idLogDetalle, proceso, registros, rutaArchivo)
	// 	if err != nil {
	// 		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
	// 		return "", modelos.ErrorFormateado{Mensaje: err.Error()}
	// 	}
	// } else if formato == "xml" {
	// 	// Ruta completa del archivo
	// 	nombreControl += ".xml"
	// 	rutaArchivo := filepath.Join(rutaCarpeta, nombreControl)
	// 	// Utilizar funcion para txt
	// 	name, err = CargarXml(db, idLogDetalle, proceso, registros, rutaArchivo)
	// 	if err != nil {
	// 		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
	// 		return "", modelos.ErrorFormateado{Mensaje: err.Error()}
	// 	}

	// Insertar o actualizar proceso en ext_procesados
	if err = ProcesadosControl(db, proceso.Id_procesado, proceso.Id_modelo, fecha, fecha2, version, len(registros), filepath.Join(rutaCarpeta, nombreControl), procesado); err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", modelos.ErrorFormateado{Mensaje: err.Error()}
	}

	// Logueo
	_, err = db.Exec("CALL extractor.act_log_detalle($1, 'F', $2)", idLogDetalle, fmt.Sprintf("Archivo guardado en: \"%s\"", filepath.Join(rutaCarpeta, nombreControl)))
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", modelos.ErrorFormateado{Mensaje: err.Error()}
	}
	_, err = db.Exec("CALL extractor.etl_ending($1)", id_log)
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", modelos.ErrorFormateado{Mensaje: err.Error()}
	}

	return name, modelos.ErrorFormateado{Mensaje: ""}
}