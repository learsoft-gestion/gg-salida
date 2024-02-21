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

		rows, err := db.Query("SELECT ID_PROCESO, NOMBRE, QUERY, ARCHIVO_MODELO, CANT_FECHAS FROM EXTRACTOR.EXT_PROCESOS")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()
		defer rows.Close()

		var dtoProcesos []modelos.DTOproceso
		for rows.Next() {
			var id int
			var nombre string
			var query string
			var modelo string
			var cant_fechas int
			if err := rows.Scan(&id, &nombre, &query, &modelo, &cant_fechas); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			dtoProceso := modelos.DTOproceso{
				Id:          id,
				Nombre:      nombre,
				Cant_fechas: cant_fechas,
			}
			dtoProcesos = append(dtoProcesos, dtoProceso)
			proceso := modelos.Proceso{
				Id:          id,
				Nombre:      nombre,
				Query:       query,
				Modelo:      modelo,
				Cant_fechas: cant_fechas,
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
		fmt.Println("Datos recibidos: ", datos)

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

		err = leerCrearExcel(proc, datos, "./salida/salida2.xlsx")
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
	// Leer archivo json
	// var personas []map[string]interface{}
	// contenido, err := os.ReadFile("./datos.json")
	// if err != nil {
	// 	panic(err.Error())
	// }
	// err = json.Unmarshal(contenido, &personas)
	// if err != nil {
	// 	panic(err.Error())
	// }

	// err = leerCrearExcel("./templates/camioneros.xlsx", personas, "./salida/salida2.xlsx")
	// if err != nil {
	// 	panic(err.Error())
	// }

	// fmt.Println("Datos insertados en salida2.xlsx")
}

func leerCrearExcel(proceso modelos.Proceso, datos modelos.DTOdatos, nombreSalida string) error {

	db, err := conexiones.ConectarBase("postgres", "test", "postgres")
	if err != nil {
		return err
	}
	defer db.Close()

	sql, err := conexiones.ConectarBase("recibos", "test", "sqlserver")
	if err != nil {
		return err
	}
	defer sql.Close()

	id_log, idLogDetalle, err := src.Logueo(db, proceso.Nombre)
	if err != nil {
		src.ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
	}

	_, err = src.Extractor(db, sql, proceso, datos, idLogDetalle)
	if err != nil {
		src.ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
	}

	_, err = db.Exec("CALL extractor.act_log_detalle($1, 'F', $2)", idLogDetalle, fmt.Sprintf("Archivo guardado en: \"%s\"", nombreSalida))
	if err != nil {
		src.ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
	}
	_, err = db.Exec("CALL extractor.etl_ending($1)", id_log)
	if err != nil {
		src.ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
	}

	return nil
}
