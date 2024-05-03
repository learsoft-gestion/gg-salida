package main

import (
	"Nueva/conexiones"
	"Nueva/handlers"
	"fmt"
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

	db, err := conexiones.ConectarBase("postgres", os.Getenv("CONN_POSTGRES"), "postgres")
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	router := mux.NewRouter()

	// Agrega un manejador OPTIONS global
	router.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.Use(corsHandler)

	// Carga de archivos estaticos
	router.PathPrefix("/salida/").Handler(http.StripPrefix("/salida/", http.FileServer(http.Dir("./salida"))))
	router.PathPrefix("/templates/").Handler(http.StripPrefix("/templates/", http.FileServer(http.Dir("./templates"))))

	// router.HandleFunc("/", indexHandler)
	// router.HandleFunc("/a-convenios", conveniosHandler)
	// router.HandleFunc("/a-modelos", aModelosHandler)
	// router.HandleFunc("/migrador", migradorHandler)

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
	router.HandleFunc("/migrador/empresas", handlers.MigradorGetEmpresas(db))
	router.HandleFunc("/migrador/convenios", handlers.MigradorGetConvenios(db))
	router.HandleFunc("/migrador/periodos", handlers.MigradorGetPeriodos(db))
	router.HandleFunc("/migrador/archivos", handlers.ProcesarArchivo(db))

	srv := &http.Server{
		Addr:    os.Getenv("SV_ADDR"),
		Handler: router,
	}

	fmt.Println("Listening at ", os.Getenv("SV_ADDR"))
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println(err.Error())
	}
}

// Middleware para manejar CORS
func corsHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Permitir solicitudes desde cualquier origen
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Permitir ciertos métodos HTTP
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PATCH")

		// Permitir ciertos encabezados
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Si la solicitud es de tipo OPTIONS, responder con éxito y terminar la cadena de middleware
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Continuar con el siguiente manejador
		next.ServeHTTP(w, r)
	})
}

// func indexHandler(w http.ResponseWriter, r *http.Request) {
// 	// Puedes usar plantillas si deseas
// 	renderTemplate(w, "../client/templates/index.html", nil)
// }

// func conveniosHandler(w http.ResponseWriter, r *http.Request) {
// 	renderTemplate(w, "../client/templates/convenios.html", nil)
// }

// func aModelosHandler(w http.ResponseWriter, r *http.Request) {
// 	renderTemplate(w, "../client/templates/modelos.html", nil)
// }

// func migradorHandler(w http.ResponseWriter, r *http.Request) {
// 	renderTemplate(w, "../client/templates/migrador.html", nil)
// }

// func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
// 	t, err := template.ParseFiles(tmpl)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	err = t.Execute(w, data)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }
