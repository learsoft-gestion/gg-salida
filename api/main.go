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

var procesos []modelos.Proceso

func getProcesos(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rows, err := db.Query("SELECT id_modelo, nombre, filtro_personas, filtro_recibos, formato_salida FROM extractor.ext_modelos")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// defer db.Close()
		// defer rows.Close()

		var dtoProcesos []modelos.DTOproceso
		for rows.Next() {
			var id int
			var nombre string
			var filtro_personas sql.NullString
			var filtroPersonas string
			var filtro_recibos sql.NullString
			var filtroRecibos string
			var formato_salida string
			if err := rows.Scan(&id, &nombre, &filtro_personas, &filtro_recibos, &formato_salida); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			dtoProceso := modelos.DTOproceso{
				Id:     id,
				Nombre: nombre,
			}
			dtoProcesos = append(dtoProcesos, dtoProceso)
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
			}
			procesos = append(procesos, proceso)
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(dtoProcesos); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

func sender(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		var datos modelos.DTOdatos
		err := json.NewDecoder(r.Body).Decode(&datos)
		if err != nil {
			http.Error(w, "Error decodificando JSON", http.StatusBadRequest)
			return
		}
		// fmt.Println("Datos recibidos: ", datos)

		// var nombre string
		// var query string
		var proc modelos.Proceso
		for _, element := range procesos {
			if element.Id == datos.Id {
				// nombre = element.Nombre
				// query = element.Query
				proc = element
			}
		}

		err = procesador(proc, datos)
		if err != nil {
			http.Error(w, "Error decodificando JSON", http.StatusBadRequest)
			fmt.Println(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		respuesta := map[string]string{
			"mensaje": "Datos recibidos",
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

	db, err := conexiones.ConectarBase("postgres", "test", "postgres")
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	router := mux.NewRouter()
	router.HandleFunc("/procesos", getProcesos(db))
	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/send", sender).Methods("POST")

	srv := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	fmt.Println("Listening...")
	srv.ListenAndServe()
}

func procesador(proceso modelos.Proceso, datos modelos.DTOdatos) error {

	db, err := conexiones.ConectarBase("postgres", "test", "postgres")
	if err != nil {
		return err
	}
	defer db.Close()

	sql, err := conexiones.ConectarBase("recibos", "prod", "sqlserver")
	if err != nil {
		return err
	}
	defer sql.Close()

	id_log, idLogDetalle, err := src.Logueo(db, proceso.Nombre)
	if err != nil {
		src.ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
	}

	query := ""
	db.QueryRow("SELECT texto_query FROM extractor.ext_query;").Scan(&query)
	proceso.Query = query

	registros, err := src.Extractor(db, sql, proceso, datos, idLogDetalle)
	if err != nil {
		src.ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return err
	}

	// Fecha para el nombre de salida
	fecha := time.Now()
	fechaFormateada := fecha.Format("20060102")
	var nombreSalida string

	// Formato del archivo de salida
	formato := strings.ToLower(proceso.Formato_salida)

	if formato == "xls" {
		nombreSalida = fmt.Sprintf("../salida/%s_%s.xlsx", proceso.Nombre, fechaFormateada)
		err = src.CargarExcel(db, idLogDetalle, proceso, registros, nombreSalida)
		if err != nil {
			src.ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return err
		}
	} else if formato == "txt" {
		nombreSalida = fmt.Sprintf("../salida/%s_%s.txt", proceso.Nombre, fechaFormateada)
		// Utilizar funcion para txt
		err = src.CargarTxt(db, idLogDetalle, proceso, registros, nombreSalida)
		if err != nil {
			src.ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return err
		}
	}

	_, err = db.Exec("CALL extractor.act_log_detalle($1, 'F', $2)", idLogDetalle, fmt.Sprintf("Archivo guardado en: \"%s\"", nombreSalida))
	if err != nil {
		src.ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return err
	}
	_, err = db.Exec("CALL extractor.etl_ending($1)", id_log)
	if err != nil {
		src.ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return err
	}

	return nil
}
