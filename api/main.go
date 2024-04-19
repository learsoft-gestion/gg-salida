package main

import (
	"Nueva/conexiones"
	"Nueva/handlers"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error: ", err.Error())
	}

	db, err := conexiones.ConectarBase("postgres", "test", "postgres")
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	router := mux.NewRouter()

	// Carga de archivos estaticos
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("../client/static"))))
	router.PathPrefix("/salida/").Handler(http.StripPrefix("/salida/", http.FileServer(http.Dir("../salida"))))

	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/a-convenios", conveniosHandler)
	router.HandleFunc("/a-modelos", aModelosHandler)
	router.HandleFunc("/migrador", migradorHandler)

	router.HandleFunc("/modelos", handlers.ModelosHandler(db))
	router.HandleFunc("/convenios", handlers.GetConvenios(db))
	router.HandleFunc("/empresas", handlers.GetEmpresas(db))
	router.HandleFunc("/empresas/{id_convenio}", handlers.GetEmpresas(db))
	router.HandleFunc("/conceptos", handlers.GetConceptos(db))
	router.HandleFunc("/conceptos/{id_convenio}", handlers.GetConceptos(db))
	router.HandleFunc("/conceptos/{id_empresa}", handlers.GetConceptos(db))
	router.HandleFunc("/conceptos/{id_convenio}/{id_empresa}", handlers.GetConceptos(db))
	router.HandleFunc("/procesos", handlers.GetProcesos(db))
	router.HandleFunc("/send", handlers.Sender(db)).Methods("POST")
	router.HandleFunc("/multiple", handlers.MultipleSend(db)).Methods("POST")
	router.HandleFunc("/restantes", handlers.ProcesosRestantes(db))
	router.HandleFunc("/clientes", handlers.GetClientes(db))
	router.HandleFunc("/migrador/procesos", handlers.Migrador(db))

	srv := &http.Server{
		Addr:    os.Getenv("SV_ADDR"),
		Handler: router,
	}

	fmt.Println("Listening...")
	srv.ListenAndServe()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Puedes usar plantillas si deseas
	renderTemplate(w, "../client/templates/index.html", nil)
}

func conveniosHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "../client/templates/convenios.html", nil)
}

func aModelosHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "../client/templates/modelos.html", nil)
}

func migradorHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "../client/templates/migrador.html", nil)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
