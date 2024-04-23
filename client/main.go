package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// Carga de archivos estaticos
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	router.PathPrefix("/salida/").Handler(http.StripPrefix("/salida/", http.FileServer(http.Dir("../salida"))))

	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/a-convenios", conveniosHandler)
	router.HandleFunc("/a-modelos", aModelosHandler)
	router.HandleFunc("/migrador", migradorHandler)

	srv := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	fmt.Println("Listening at :8000 port")
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println(err.Error())
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Puedes usar plantillas si deseas
	renderTemplate(w, "./templates/index.html", nil)
}

func conveniosHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "./templates/convenios.html", nil)
}

func aModelosHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "./templates/modelos.html", nil)
}

func migradorHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "./templates/migrador.html", nil)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
