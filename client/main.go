package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"html/template"
	"io"
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

	router := mux.NewRouter()

	// Middleware :: To prevent unauthorized access
	router.Use(authMiddleware)

	// Carga de archivos estaticos
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	// router.PathPrefix("/api/salida/").Handler(http.StripPrefix("/api/salida/", http.FileServer(http.Dir("../api/salida"))))
	// router.PathPrefix("/api/templates/").Handler(http.StripPrefix("/api/templates/", http.FileServer(http.Dir("../api/templates"))))

	// :: GG :: Protected routes
	protectedRoutes := router.NewRoute().Subrouter()
	protectedRoutes.HandleFunc("/a-convenios", conveniosHandler)
	protectedRoutes.HandleFunc("/a-modelos", aModelosHandler)
	protectedRoutes.HandleFunc("/a-alicuotas", aAlicuotasHandler)
	protectedRoutes.HandleFunc("/a-personalinterno", aPersonalinternoHandler)
	protectedRoutes.HandleFunc("/migrador", migradorHandler)
	protectedRoutes.HandleFunc("/logout", logoutHandler)
	protectedRoutes.HandleFunc("/consulta", consultaHandler)
	protectedRoutes.HandleFunc("/proyeccion", proyeccionHandler)

	// :: GG :: Public routes
	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/login", loginHandler).Methods("GET", "POST")
	router.HandleFunc("/backend-url", backendUrl())

	srv := &http.Server{
		Addr:    os.Getenv("SV_ADDR"),
		Handler: router,
	}

	fmt.Println("Listening at ", os.Getenv("SV_ADDR"))
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println(err.Error())
	}
}

func backendUrl() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Cargar variables de entorno
		if err := godotenv.Load(); err != nil {
			fmt.Println("Error: ", err.Error())
		}

		prefijoURL := os.Getenv("URL_BACK")

		type Response struct {
			PrefijoURL string `json:"prefijoURL"`
		}
		response := Response{PrefijoURL: prefijoURL}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Puedes usar plantillas si deseas
	renderTemplate(w, "./templates/index.html", nil)
}

// :: GG :: Login handler
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		cookie, err := r.Cookie("sindicatos-loggedin")
		if err == nil && cookie.Value == "true" {
			// User already logged in, redirect to home
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		// Show login page
		renderTemplate(w, "./templates/login.html", nil)
		return
	} else if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Check if username or password is empty
		if username == "" || password == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Execute the login process
		response, err := loginUser(w, username, password)
		fmt.Println("Response login :: ", response)
		if err != nil {
			// On error, redirect back to the login with an error
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Assuming the response contains some form of success indicator
		if response == "success" { // Adjust the condition based on actual API response
			// On success, redirect to the home page
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			// If login is not successful, redirect back to the login page
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	}
}

// :: GG :: Login user
func loginUser(w http.ResponseWriter, username string, password string) (string, error) {

	token, err := connectToAPI()
	if err != nil {
		return "No fue posible establecer la conexión... ", err
	}

	// Create an instance of the struct to hold the parsed data
	var tokenStructure map[string]any

	errToken := json.Unmarshal([]byte(token), &tokenStructure)

	if errToken != nil {
		fmt.Println("Errores:", err)
	}

	tokenParsed := tokenStructure["accessToken"]

	// Hardcoded credentials for testing
	hardcodedUsername := os.Getenv("TEST_USER")
	hardcodedPassword := os.Getenv("TEST_PASS")

	// Check if provided credentials match
	if username == hardcodedUsername && password == hardcodedPassword {
		// Set a cookie to signify successful login
		http.SetCookie(w, &http.Cookie{
			Name:     "sindicatos-loggedin",
			Value:    "true",
			Path:     "/",
			MaxAge:   3600, // Expires after one hour
			HttpOnly: true, // Prevents JavaScript access to the cookie
		})
		return "success", nil // Simulated successful login response
	}

	// JSON request body
	requestBody, err := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})

	if err != nil {
		return "", err
	}

	// POST request
	req, err := http.NewRequest("POST", os.Getenv("API_URL")+"/logistica/login_user", bytes.NewBuffer(requestBody))

	if err != nil {
		return "", err
	}

	// Encode credentials and set the Authorization header
	auth := tokenParsed.(string)
	// encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", "Bearer "+auth)

	// Set content type
	req.Header.Add("Content-Type", "application/json")

	// Send the request using a new HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}

// :: GG :: Logout handler
func logoutHandler(w http.ResponseWriter, r *http.Request) {

	// Clear the login cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "sindicatos-loggedin",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Immediately expire the cookie
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func conveniosHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "./templates/convenios.html", nil)
}

func aModelosHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "./templates/modelos.html", nil)
}

func aAlicuotasHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "./templates/alicuotas.html", nil)
}

func aPersonalinternoHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "./templates/personalinterno.html", nil)
}

func migradorHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "./templates/migrador.html", nil)
}

func consultaHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "./templates/consulta.html", nil)
}

func proyeccionHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "./templates/proyeccion.html", nil)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		http.Error(w, "Failed to load template: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := t.Execute(w, data); err != nil {
		http.Error(w, "Failed to execute template: "+err.Error(), http.StatusInternalServerError)
	}
}

// :: GG :: Connect to API
func connectToAPI() (string, error) {

	// Create the request body as a JSON
	requestBody, err := json.Marshal(map[string]string{
		"provider": os.Getenv("API_PROVIDER"),
	})

	if err != nil {
		return "", err
	}

	// Simulated API request
	resp, err := http.Post(os.Getenv("API_URL")+"/auth", "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	token, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(token), nil

}

// :: GG :: Auth middleware
func authMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Exclude login page from authentication check
		if r.URL.Path == "/login" {
			next.ServeHTTP(w, r)
			return
		}

		// Check for the cookie
		cookie, err := r.Cookie("sindicatos-loggedin")

		if err != nil || cookie.Value != "true" {
			// If not logged in, go to login
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// If logged in, proceed with the original handler
		next.ServeHTTP(w, r)
	})

}
