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
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// func (p *Page) save() error {
// 	filename := p.Title + ".txt"
// 	return os.WriteFile(filename, p.Body, 0600)
// }

// func loadPage(title string) (*Page, error) {
// 	filename := title + ".txt"
// 	body, err := os.ReadFile(filename)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &Page{Title: title, Body: body}, nil
// }

func homeHandler(w http.ResponseWriter, r *http.Request) {
	templates = template.Must(template.New("../client/index.html").ParseFiles("../client/index.html"))
	w.WriteHeader(http.StatusAccepted)

	// p, err := loadPage(title)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }
	renderTemplate(w, "index", nil)

}

// var procesos []modelos.Proceso
// var dtoProcesos []modelos.DTOproceso

// func getProcesos(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		rows, err := db.Query("SELECT em.id_modelo, ea.razon_social as nombre_empresa_adm, c.nombre as nombre_convenio, em.nombre, em.filtro_personas, em.filtro_recibos, em.formato_salida, em.archivo_modelo FROM extractor.ext_modelos em JOIN datos.empresas_adm ea ON em.id_empresa_adm = ea.id_empresa_adm JOIN datos.convenios c ON em.id_convenio = c.id_convenio where vigente")
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		for rows.Next() {
// 			var id int
// 			var empresa string
// 			var sindicato string
// 			var nombre string
// 			var filtro_personas sql.NullString
// 			var filtroPersonas string
// 			var filtro_recibos sql.NullString
// 			var filtroRecibos string
// 			var formato_salida string
// 			var archivo_modelo string

// 			if err := rows.Scan(&id, &empresa, &sindicato, &nombre, &filtro_personas, &filtro_recibos, &formato_salida, &archivo_modelo); err != nil {
// 				http.Error(w, err.Error(), http.StatusInternalServerError)
// 				return
// 			}
// 			dtoProceso := modelos.DTOproceso{
// 				Id:        id,
// 				Empresa:   empresa,
// 				Sindicato: sindicato,
// 				Nombre:    nombre,
// 			}
// 			dtoProcesos = append(dtoProcesos, dtoProceso)

// 			if filtro_personas.Valid {
// 				filtroPersonas = filtro_personas.String
// 			}
// 			if filtro_recibos.Valid {
// 				filtroRecibos = filtro_recibos.String
// 			}
// 			proceso := modelos.Proceso{
// 				Id:              id,
// 				Nombre:          nombre,
// 				Filtro_personas: filtroPersonas,
// 				Filtro_recibos:  filtroRecibos,
// 				Formato_salida:  formato_salida,
// 				Archivo_modelo:  archivo_modelo,
// 			}
// 			procesos = append(procesos, proceso)
// 		}
// 		rows.Close()

// 		dtoSelect := modelos.DTOselect{
// 			Empresas:  empresas,
// 			Convenios: convenios,
// 			Procesos:  dtoProcesos,
// 		}

// 		w.Header().Set("Content-Type", "application/json")

// 		if err := json.NewEncoder(w).Encode(dtoSelect); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 	}
// }

var convenios []modelos.Option

func getConvenios(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rows, err := db.Query("SELECT em.id_convenio, c.nombre as nombre_convenio FROM extractor.ext_modelos em JOIN extractor.ext_convenios c ON em.id_convenio = c.id_convenio where vigente")
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

var empresas []modelos.Option

func getEmpresas(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id_convenio, err := strconv.Atoi(vars["id_convenio"])
		if err != nil {
			http.Error(w, "ID invalido", http.StatusBadRequest)
			return
		}
		fmt.Println("CONVENIO: ", id_convenio)
		rows, err := db.Query("SELECT em.id_empresa_adm, ea.razon_social as nombre_empresa_adm FROM extractor.ext_modelos em JOIN datos.empresas_adm ea ON em.id_empresa_adm = ea.id_empresa_adm where id_convenio = $1 and vigente", id_convenio)
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

var DTOprocesos []modelos.Option
var procesos []modelos.Proceso

func getProcesos(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Reinicio variables
		DTOprocesos = nil
		procesos = nil

		// Extraigo IDs de la request
		vars := mux.Vars(r)
		id_convenio, err := strconv.Atoi(vars["id_convenio"])
		if err != nil {
			http.Error(w, "ID invalido", http.StatusBadRequest)
			return
		}
		id_empresa, err := strconv.Atoi(vars["id_empresa"])
		if err != nil {
			http.Error(w, "ID invalido", http.StatusBadRequest)
			return
		}

		rows, err := db.Query("SELECT em.id_modelo, ea.razon_social as nombre_empresa_adm, c.nombre as nombre_convenio, em.nombre, em.filtro_personas, em.filtro_recibos, em.formato_salida, em.archivo_modelo FROM extractor.ext_modelos em JOIN datos.empresas_adm ea ON em.id_empresa_adm = ea.id_empresa_adm JOIN extractor.ext_convenios c ON em.id_convenio = c.id_convenio where em.id_convenio = $1 and em.id_empresa_adm = $2 and vigente", id_convenio, id_empresa)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for rows.Next() {
			var id int
			var empresa string
			var sindicato string
			var nombre string
			var filtro_personas sql.NullString
			var filtroPersonas string
			var filtro_recibos sql.NullString
			var filtroRecibos string
			var formato_salida string
			var archivo_modelo string

			if err := rows.Scan(&id, &empresa, &sindicato, &nombre, &filtro_personas, &filtro_recibos, &formato_salida, &archivo_modelo); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			DTOprocesos = src.AddToSet(DTOprocesos, modelos.Option{Id: id, Nombre: nombre})

			if filtro_personas.Valid {
				filtroPersonas = filtro_personas.String
			}
			if filtro_recibos.Valid {
				filtroRecibos = filtro_recibos.String
			}
			proceso := modelos.Proceso{
				Id:              id,
				Nombre:          nombre,
				Filtro_personas: filtroPersonas,
				Filtro_recibos:  filtroRecibos,
				Formato_salida:  formato_salida,
				Archivo_modelo:  archivo_modelo,
			}
			procesos = append(procesos, proceso)
		}
		rows.Close()

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(DTOprocesos); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

// func getConceptos(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// Extraigo IDs de la request
// 		vars := mux.Vars(r)
// 		id_convenio, err := strconv.Atoi(vars["id_convenio"])
// 		if err != nil {
// 			http.Error(w, "ID invalido", http.StatusBadRequest)
// 			return
// 		}
// 		id_empresa, err := strconv.Atoi(vars["id_empresa"])
// 		if err != nil {
// 			http.Error(w, "ID invalido", http.StatusBadRequest)
// 			return
// 		}

// 		rows, err := db.Query("SELECT em.id_empresa_adm, ea.razon_social as nombre_empresa_adm FROM extractor.ext_modelos em JOIN datos.empresas_adm ea ON em.id_empresa_adm = ea.id_empresa_adm where id_convenio = $1 and vigente", id_convenio)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		var dtoConceptos modelos.Conceptos
// 		var conceptos []string
// 		var tipos []string
// 		for rows.Next() {
// 			var concepto string
// 			var tipo string
// 			if err = rows.Scan(&concepto, &tipo); err != nil {
// 				http.Error(w, err.Error(), http.StatusInternalServerError)
// 				return
// 			}

// 			conceptos = src.AddToSet(conceptos, {Id: id, Nombre: empresa})

// 		}
// 		w.Header().Set("Content-Type", "application/json")

// 		if err := json.NewEncoder(w).Encode(empresas); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 	}
// }

func sender(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		var datos modelos.DTOdatos
		err := json.NewDecoder(r.Body).Decode(&datos)
		if err != nil {
			http.Error(w, "Error decodificando JSON", http.StatusBadRequest)
			return
		}
		// fmt.Println("Datos recibidos: ", datos)
		var procs []modelos.Proceso
		// var proc modelos.Proceso
		for _, id := range datos.IDs {
			for _, element := range procesos {
				if element.Id == id {
					// proc = element
					procs = append(procs, element)
				}
			}
		}

		var resultado []string
		for _, proc := range procs {
			result, errFormateado := procesador(proc, datos.Fecha, datos.Fecha2, datos.Forzado)
			if result != "" {
				resultado = append(resultado, result)
			}
			if (errFormateado.Mensaje != "") && (!errFormateado.Procesado) {
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
			} else if errFormateado.Procesado {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				respuesta := modelos.Respuesta{
					Mensaje:         errFormateado.Mensaje,
					Archivos_salida: nil,
					Procesado:       true,
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

func force(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var datos modelos.DTOdatos
		err := json.NewDecoder(r.Body).Decode(&datos)
		if err != nil {
			http.Error(w, "Error decodificando JSON", http.StatusBadRequest)
			return
		}
		// fmt.Println("Datos recibidos: ", datos)
		var procs []modelos.Proceso
		for _, id := range datos.IDs {
			for _, element := range procesos {
				if element.Id == id {
					procs = append(procs, element)
				}
			}
		}

		var resultado string
		for _, proc := range procs {
			result, errFormateado := procesador(proc, datos.Fecha, datos.Fecha2, datos.Forzado)
			if errFormateado.Mensaje != "" && !errFormateado.Procesado {
				errString := "Error en " + proc.Nombre + ": " + errFormateado.Mensaje
				// http.Error(w, errString, http.StatusBadRequest)
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				respuesta := map[string]string{
					"mensaje":        errString,
					"archivo_salida": "",
				}
				jsonResp, _ := json.Marshal(respuesta)
				w.Write(jsonResp)
				return
				// } else if errFormateado.Procesado {
				// 	w.WriteHeader(http.StatusOK)
				// 	w.Header().Set("Content-Type", "application/json")
				// 	respuesta := map[string]string{
				// 		"mensaje":        "",
				// 		"archivo_salida": "",
				// 		"procesado":      "si",
				// 	}
				// 	jsonResp, _ := json.Marshal(respuesta)
				// 	w.Write(jsonResp)
				// 	fmt.Println(err)
				// 	return
			}
			resultado = result
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		respuesta := map[string]string{
			"mensaje":        "Datos recibidos y procesados",
			"archivo_salida": resultado,
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
	router.HandleFunc("/empresas/{id_convenio}", getEmpresas(db))
	// router.HandleFunc("/conceptos/{id_convenio}/{id_empresa}", getConceptos(db))
	router.HandleFunc("/procesos/{id_convenio}/{id_empresa}", getProcesos(db))
	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/send", sender).Methods("POST")
	router.HandleFunc("/force", force).Methods("POST")

	srv := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	fmt.Println("Listening...")
	srv.ListenAndServe()
}

func procesador(proceso modelos.Proceso, fecha string, fecha2 string, forzado bool) (string, modelos.ErrorFormateado) {

	db, err := conexiones.ConectarBase("postgres", "test", "postgres")
	if err != nil {
		return "", modelos.ErrorFormateado{Mensaje: err.Error()}
	}
	defer db.Close()

	var cuenta int
	var version int
	// Verificar si el proceso ya se corrió
	if !forzado { // Si la ejecucion no viene forzada continuo evaluando si este modelo ya se procesó
		fecha_desde, _ := time.Parse("200601", fecha)
		fmt.Println(fecha_desde)
		if fecha2 == "" {
			err = db.QueryRow("select count(*) from extractor.ext_procesados where id_modelo = $1 and fecha_desde = $2", proceso.Id, fecha_desde).Scan(&cuenta)
			if err != nil {
				return "", modelos.ErrorFormateado{Mensaje: err.Error()}
			}
		} else {
			fecha_hasta, _ := time.Parse("200601", fecha2)
			err = db.QueryRow("select count(*) from extractor.ext_procesados where id_modelo = $1 and fecha_desde = $2 and fecha_hasta = $3", proceso.Id, fecha_desde, fecha_hasta).Scan(&cuenta)
			if err != nil {
				return "", modelos.ErrorFormateado{Mensaje: err.Error()}
			}
		}
	}
	if cuenta > 0 {
		fmt.Println("Este modelo ya ha sido procesado.")
		version = cuenta + 1
		// return "", modelos.ErrorFormateado{Mensaje: fmt.Errorf("el modelo ya ha sido procesado").Error(), Procesado: true}
	} else {
		version = cuenta + 1
	}

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
		return "", modelos.ErrorFormateado{Mensaje: err.Error()}
	}

	// Fecha para el nombre de salida
	fechaSalida := time.Now()
	fechaFormateada := fechaSalida.Format("20060102")
	var nombreSalida string

	// Formato del archivo de salida
	formato := strings.ToLower(proceso.Formato_salida)
	var name string
	if formato == "xls" {
		nombreSalida = fmt.Sprintf("../salida/%s_%s.xlsx", proceso.Nombre, fechaFormateada)
		name, err = src.CargarExcel(db, idLogDetalle, proceso, registros, nombreSalida)
		if err != nil {
			src.ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return "", modelos.ErrorFormateado{Mensaje: err.Error()}
		}
	} else if formato == "txt" {
		nombreSalida = fmt.Sprintf("../salida/%s_%s.txt", proceso.Nombre, fechaFormateada)
		// Utilizar funcion para txt
		name, err = src.CargarTxt(db, idLogDetalle, proceso, registros, nombreSalida)
		if err != nil {
			src.ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return "", modelos.ErrorFormateado{Mensaje: err.Error()}
		}
	}

	// Insertar nuevo proceso en ext_procesados
	if err = src.Procesados(db, proceso.Id, fecha, fecha2, version, len(registros), nombreSalida); err != nil {
		src.ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", modelos.ErrorFormateado{Mensaje: err.Error()}
	}

	// Logueo
	_, err = db.Exec("CALL extractor.act_log_detalle($1, 'F', $2)", idLogDetalle, fmt.Sprintf("Archivo guardado en: \"%s\"", nombreSalida))
	if err != nil {
		src.ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", modelos.ErrorFormateado{Mensaje: err.Error()}
	}
	_, err = db.Exec("CALL extractor.etl_ending($1)", id_log)
	if err != nil {
		src.ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", modelos.ErrorFormateado{Mensaje: err.Error()}
	}

	return name, modelos.ErrorFormateado{Mensaje: ""}
}
