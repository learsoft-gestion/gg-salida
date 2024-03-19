package main

import (
	"Nueva/conexiones"
	"Nueva/modelos"
	"Nueva/src"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	templates = template.Must(template.New("../client/index.html").ParseFiles("../client/index.html"))
	w.WriteHeader(http.StatusAccepted)

	renderTemplate(w, "index", nil)

}

var convenios []modelos.Option
var empresas []modelos.Option
var DTOprocesos []modelos.DTOproceso
var procesos []modelos.Proceso

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

		query := fmt.Sprintf("select em.id_modelo, c.nombre as nombre_convenio, ea.razon_social as nombre_empresa_adm, ec.nombre as nombre_concepto, em.nombre, et.nombre as nombre_tipo, ep.fecha_desde, ep.fecha_hasta, ep.nombre_salida, ep.version, ep.fecha_ejecucion, coalesce(nombre_salida is not null, false) procesado, case when version is null then 'lanzar' when version = max(ep.version) over(partition by em.id_modelo, em.id_empresa_adm, em.id_concepto, em.id_convenio) then 'relanzar' end boton from extractor.ext_modelos em left join extractor.ext_procesados ep on em.id_modelo = ep.id_modelo join datos.empresas_adm ea ON em.id_empresa_adm = ea.id_empresa_adm join extractor.ext_convenios c ON em.id_convenio = c.id_convenio join extractor.ext_conceptos ec on em.id_concepto = ec.id_concepto join extractor.ext_tipos et on em.id_tipo = et.id_tipo where em.id_convenio = %v and ((ep.fecha_desde = '%s' and ep.fecha_hasta = '%s') or ep.fecha_desde is null)", id_convenio, fechaFormateada, fechaFormateada2)

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
			query += fmt.Sprintf(" and coalesce(nombre_salida is not null, false) = %v", procesado)
		}
		if len(jurisdiccion) > 0 {
			query += " and UPPER(em.nombre) like '%" + strings.ToUpper(jurisdiccion) + "%'"
		}
		query += " order by nombre_empresa_adm, nombre_concepto, nombre, nombre_tipo, ep.version"
		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			// var proceso modelos.Proceso
			var DTOproceso modelos.DTOproceso
			var version sql.NullInt16
			var ult_ejecucion sql.NullString
			var boton sql.NullString

			if err := rows.Scan(&DTOproceso.Id, &DTOproceso.Convenio, &DTOproceso.Empresa, &DTOproceso.Concepto, &DTOproceso.Nombre, &DTOproceso.Tipo, &DTOproceso.Fecha_desde, &DTOproceso.Fecha_hasta, &DTOproceso.Nombre_salida, &version, &ult_ejecucion, &DTOproceso.Procesado, &boton); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
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
			}
			if boton.Valid {
				DTOproceso.Boton = boton.String
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
			var datos modelos.DTOdatos
			err := json.NewDecoder(r.Body).Decode(&datos)
			if err != nil {
				http.Error(w, "Error decodificando JSON", http.StatusBadRequest)
				return
			}
			datos.Fecha = src.FormatoFecha(datos.Fecha)
			datos.Fecha2 = src.FormatoFecha(datos.Fecha2)
			var placeholders []string
			// for i := range datos.Id {
			// 	placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
			// }
			placeholders = append(placeholders, "$1")
			queryModelos := fmt.Sprintf("SELECT em.id_modelo, em.id_empresa_adm, ea.razon_social as nombre_empresa, c.nombre as nombre_convenio, em.nombre, c.filtro as filtro_convenio, em.filtro_personas, em.filtro_recibos, em.formato_salida, em.archivo_modelo FROM extractor.ext_modelos em JOIN datos.empresas_adm ea ON em.id_empresa_adm = ea.id_empresa_adm JOIN extractor.ext_convenios c ON em.id_convenio = c.id_convenio where vigente and em.id_modelo in (%s)", strings.Join(placeholders, ","))
			// fmt.Println("Query modelos: ", queryModelos)
			stmt, err := db.Prepare(queryModelos)
			if err != nil {
				http.Error(w, "Error al preparar query", http.StatusBadRequest)
				return
			}
			defer stmt.Close()
			var args []interface{}
			// for _, arg := range datos.IDs {
			// 	args = append(args, arg)
			// }
			args = append(args, datos.Id)
			rows, err := stmt.Query(args...)
			if err != nil {
				http.Error(w, "Error al ejecutar el query", http.StatusBadRequest)
				return
			}
			defer rows.Close()
			for rows.Next() {
				var proceso modelos.Proceso
				err = rows.Scan(&proceso.Id, &proceso.Id_empresa, &proceso.Nombre_empresa, &proceso.Nombre_convenio, &proceso.Nombre, &proceso.Filtro_convenio, &proceso.Filtro_personas, &proceso.Filtro_recibos, &proceso.Formato_salida, &proceso.Archivo_modelo)
				if err != nil {
					fmt.Println(err.Error())
					http.Error(w, "Error al escanear proceso", http.StatusBadRequest)
					return
				}
				procesos = append(procesos, proceso)
			}

			var resultado []string
			for _, proc := range procesos {
				var cuenta int
				var version int
				// Verificar si el proceso ya se corrió

				err = db.QueryRow("select count(*) from extractor.ext_procesados where id_modelo = $1 and fecha_desde = $2 and fecha_hasta = $3", proc.Id, datos.Fecha, datos.Fecha2).Scan(&cuenta)
				if err != nil {
					fmt.Println(err.Error())
					http.Error(w, "Error al escanear proceso", http.StatusBadRequest)
					return
				}

				version = cuenta + 1

				result, errFormateado := procesador(proc, datos.Fecha, datos.Fecha2, version)
				if result != "" {
					resultado = append(resultado, result)
				}
				if errFormateado.Mensaje != "" {
					errString := "Error en " + proc.Nombre + ": " + errFormateado.Mensaje
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
				err = rows.Scan(&proceso.Id, &proceso.Id_empresa, &proceso.Nombre_empresa, &proceso.Nombre_convenio, &proceso.Nombre, &proceso.Filtro_convenio, &proceso.Filtro_personas, &proceso.Filtro_recibos, &proceso.Formato_salida, &proceso.Archivo_modelo)
				if err != nil {
					fmt.Println(err.Error())
					http.Error(w, "Error al escanear proceso", http.StatusBadRequest)
					return
				}
				procesos = append(procesos, proceso)
			}

			var resultado []string
			for _, proc := range procesos {
				var cuenta int
				var version int

				// Verificar si el proceso ya se corrió
				err = db.QueryRow("select count(*) from extractor.ext_procesados where id_modelo = $1 and fecha_desde = $2 and fecha_hasta = $3", proc.Id, restantes.Fecha1, restantes.Fecha2).Scan(&cuenta)
				if err != nil {
					fmt.Println(err.Error())
					http.Error(w, "Error al escanear proceso", http.StatusBadRequest)
					return
				}

				version = cuenta + 1

				result, _ := procesador(proc, restantes.Fecha1, restantes.Fecha2, version)
				if result != "" {
					resultado = append(resultado, result)
				}

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

func procesosRestantes(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id_convenio := r.URL.Query().Get("convenio")
		fecha1 := r.URL.Query().Get("fecha1")
		fechaFormateada := src.FormatoFecha(fecha1)
		fecha2 := r.URL.Query().Get("fecha2")
		fechaFormateada2 := src.FormatoFecha(fecha2)

		query := fmt.Sprintf("SELECT modelo.id_modelo, c.nombre from extractor.ext_modelos modelo LEFT JOIN extractor.ext_procesados ep ON modelo.id_modelo = ep.id_modelo JOIN datos.empresas_adm ea ON modelo.id_empresa_adm = ea.id_empresa_adm JOIN extractor.ext_convenios c ON modelo.id_convenio = c.id_convenio JOIN extractor.ext_conceptos ec ON modelo.id_concepto = ec.id_concepto JOIN extractor.ext_tipos et ON modelo.id_tipo = et.id_tipo WHERE modelo.id_convenio = %v AND ((ep.fecha_desde = '%s' AND ep.fecha_hasta = '%s') OR ep.fecha_desde IS NULL) AND ep.version IS NULL GROUP BY modelo.id_modelo, c.nombre, modelo.nombre, ep.version;", id_convenio, fechaFormateada, fechaFormateada2)

		rows, err := db.Query(query)
		if err != nil {
			fmt.Println(fmt.Println("Error al ejecutar query de procesosRestantes"))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var id_modelos []int
		var nombre_conv string
		for rows.Next() {
			var id_modelo int
			var convenio string
			rows.Scan(&id_modelo, &convenio)
			id_modelos = append(id_modelos, id_modelo)
			nombre_conv = convenio
		}

		restantes = modelos.Restantes{
			Id:       id_modelos,
			Convenio: nombre_conv,
			Fecha1:   fechaFormateada,
			Fecha2:   fechaFormateada2,
		}
		cantidad := len(id_modelos)
		resString := fmt.Sprintf("Faltan generar %v informes para el convenio %s", cantidad, nombre_conv)

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		respuesta := modelos.RespuestaRestantes{
			Mensaje: resString,
		}
		jsonResp, _ := json.Marshal(respuesta)
		w.Write(jsonResp)

	}
}

// var folder embed.FS
var templates *template.Template

func renderTemplate(w http.ResponseWriter, tmpl string, p *modelos.Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {

	// Carga de archivos estaticos
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	db, err := conexiones.ConectarBase("postgres", "test", "postgres")
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	router := mux.NewRouter()
	router.HandleFunc("/convenios", getConvenios(db))
	router.HandleFunc("/empresas", getEmpresas(db))
	router.HandleFunc("/empresas/{id_convenio}", getEmpresas(db))
	router.HandleFunc("/conceptos", getConceptos(db))
	router.HandleFunc("/conceptos/{id_convenio}", getConceptos(db))
	router.HandleFunc("/conceptos/{id_empresa}", getConceptos(db))
	router.HandleFunc("/conceptos/{id_convenio}/{id_empresa}", getConceptos(db))
	router.HandleFunc("/procesos", getProcesos(db))
	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/send", sender(db)).Methods("POST")
	router.HandleFunc("/multiple", multipleSend(db)).Methods("POST")
	router.HandleFunc("/restantes", procesosRestantes(db))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Listening...")
	srv.ListenAndServe()
}

func procesador(proceso modelos.Proceso, fecha string, fecha2 string, version int) (string, modelos.ErrorFormateado) {

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

	id_log, idLogDetalle, err := src.Logueo(db, proceso.Nombre)
	if err != nil {
		src.ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", modelos.ErrorFormateado{Mensaje: err.Error()}
	}

	query := ""
	db.QueryRow("SELECT texto_query FROM extractor.ext_query;").Scan(&query)
	proceso.Query = query

	registros, err := src.Extractor(db, sql, proceso, fecha, fecha2, idLogDetalle)
	if err != nil {
		src.ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return err.Error(), modelos.ErrorFormateado{Mensaje: err.Error()}
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
		src.ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", modelos.ErrorFormateado{Mensaje: err.Error()}
	}

	var nombreSalida string
	proceso_periodo := fecha + "-" + fecha2
	// Construir la ruta de la carpeta de salida
	rutaCarpeta := filepath.Join(directorioActual, "..", "salida", proceso.Nombre_empresa, proceso.Nombre_convenio, proceso_periodo)

	// Verificar si la carpeta de salida existe, si no, crearla
	if _, err := os.Stat(rutaCarpeta); os.IsNotExist(err) {
		if err := os.MkdirAll(rutaCarpeta, 0755); err != nil {
			fmt.Println("Error al crear la carpeta de salida:", err)
			src.ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
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

		name, err = src.CargarExcel(db, idLogDetalle, proceso, registros, rutaArchivo)
		if err != nil {
			src.ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return "", modelos.ErrorFormateado{Mensaje: err.Error()}
		}

	} else if formato == "txt" {
		// Ruta completa del archivo
		nombreSalida += ".txt"
		rutaArchivo := filepath.Join(rutaCarpeta, nombreSalida)

		// Utilizar funcion para txt
		name, err = src.CargarTxt(db, idLogDetalle, proceso, registros, rutaArchivo)
		if err != nil {
			src.ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return "", modelos.ErrorFormateado{Mensaje: err.Error()}
		}
	} else if formato == "xml" {
		// Ruta completa del archivo
		nombreSalida += ".xml"
		rutaArchivo := filepath.Join(rutaCarpeta, nombreSalida)
		// Utilizar funcion para txt
		name, err = src.CargarXml(db, idLogDetalle, proceso, registros, rutaArchivo)
		if err != nil {
			src.ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return "", modelos.ErrorFormateado{Mensaje: err.Error()}
		}
	}

	// Insertar nuevo proceso en ext_procesados
	if err = src.Procesados(db, proceso.Id, fecha, fecha2, version, len(registros), filepath.Join(rutaCarpeta, nombreSalida)); err != nil {
		src.ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", modelos.ErrorFormateado{Mensaje: err.Error()}
	}

	// Logueo
	_, err = db.Exec("CALL extractor.act_log_detalle($1, 'F', $2)", idLogDetalle, fmt.Sprintf("Archivo guardado en: \"%s\"", filepath.Join(rutaCarpeta, nombreSalida)))
	if err != nil {
		src.ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", modelos.ErrorFormateado{Mensaje: err.Error()}
	}
	_, err = db.Exec("CALL extractor.etl_ending($1)", id_log)
	if err != nil {
		src.ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", modelos.ErrorFormateado{Mensaje: err.Error()}
	}
	// El proceso termino, reinicio procesos
	procesos = nil
	return name, modelos.ErrorFormateado{Mensaje: ""}
}
