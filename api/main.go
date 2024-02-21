package main

import (
	"Nueva/conexiones"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Proceso struct {
	Id          int    `json:"id"`
	Nombre      string `json:"nombre"`
	Query       string `json:"query"`
	Modelo      string `json:"modelo"`
	Cant_fechas int    `json:"cant_fechas"`
}

type DTOproceso struct {
	Id          int    `json:"id"`
	Nombre      string `json:"nombre"`
	Cant_fechas int    `json:"cant_fechas"`
}

type ProcesosTemplate struct {
	Procesos []DTOproceso
	Title    string
}

type Page struct {
	Title string
	Body  []byte
}

type DTOdatos struct {
	Id         int
	FechaDesde string
	FechaHasta string
}

// Modelos para lectura de tabla
type Registro struct {
	Ids     string
	Valores map[string]interface{}
}

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

var procesos []Proceso

func getProcesos(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		rows, err := db.Query("SELECT ID_PROCESO, NOMBRE, QUERY, ARCHIVO_MODELO, CANT_FECHAS FROM EXTRACTOR.EXT_PROCESOS")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()
		defer rows.Close()

		var dtoProcesos []DTOproceso
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
			dtoProceso := DTOproceso{
				Id:          id,
				Nombre:      nombre,
				Cant_fechas: cant_fechas,
			}
			dtoProcesos = append(dtoProcesos, dtoProceso)
			proceso := Proceso{
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
		var datos DTOdatos
		err := json.NewDecoder(r.Body).Decode(&datos)
		if err != nil {
			http.Error(w, "Error decodificando JSON", http.StatusBadRequest)
			return
		}
		fmt.Println("Datos recibidos: ", datos)

		// var nombre string
		// var query string
		var proc Proceso
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

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
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

func ManejoErrores(postgresDb *sql.DB, idLogDetalle int, nombre string, err error) {
	text := "Fallo en: " + err.Error()
	postgresDb.Exec("CALL extractor.act_log_detalle($1, 'E', $2)", idLogDetalle, text)
	fmt.Printf("Error en %s: %s \n", nombre, text)
}

func leerCrearExcel(proceso Proceso, datos DTOdatos, nombreSalida string) error {

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

	// Comienza el proceso
	var id_log int
	err = db.QueryRow("SELECT extractor.etl_start()").Scan(&id_log)
	if err != nil {
		return err
	}
	fmt.Println("id_log: ", id_log)

	// Inicializamos registro en ext logueo detalle
	var idLogDetalle int
	err = db.QueryRow("SELECT extractor.start_log_detalle($1, $2)", id_log, proceso.Nombre).Scan(&idLogDetalle)
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
	}

	// Reemplazo de fecha en query
	var queryFinal string
	if proceso.Cant_fechas > 1 {
		query := strings.Replace(proceso.Query, "$1", datos.FechaDesde, 1)
		queryFinal = strings.Replace(query, "$2", datos.FechaHasta, 1)
	} else {
		queryFinal = strings.Replace(proceso.Query, "$1", datos.FechaDesde, 1)
	}

	// Ejecucion de query y lectura de resultados
	rows, err := sql.Query(queryFinal)
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
	}
	defer rows.Close()

	columnas, err := rows.Columns()
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
	}
	columnasNum := len(columnas)

	registros := make(map[string]Registro)

	valores := make([]interface{}, columnasNum)
	for i := range valores {
		valores[i] = new(interface{})
	}

	for rows.Next() {

		if err := rows.Scan(valores...); err != nil {
			ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		}

		registroMapa := make(map[string]interface{})
		for i, colNombre := range columnas {
			registroMapa[colNombre] = *(valores[i].(*interface{}))
		}

		id := *valores[0].(*interface{})
		idString := fmt.Sprintf("%v", id)

		registro := Registro{
			Ids:     idString,
			Valores: registroMapa,
		}
		registros[idString] = registro
		fmt.Println("Registro: ", registro)
	}

	// // Construir path
	// path := fmt.Sprintf("../templates/%s.xlsx", proceso.Nombre)

	// // Abrir archivo Excel
	// f, err := excelize.OpenFile(path)
	// if err != nil {
	// 	ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
	// }

	// // Nombre de la hoja
	// sheetName := f.GetSheetName(0)

	// // Obtener las celdas de la primera fila
	// filas, err := f.GetRows(sheetName)
	// if err != nil {
	// 	ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
	// }
	// fileNuevo := excelize.NewFile()

	// // var estilos []int
	// for i := range filas[0] {
	// 	colLetter, _ := excelize.ColumnNumberToName(i + 1)
	// 	cell := colLetter + "1"
	// 	styles, err := f.GetCellStyle(sheetName, cell)
	// 	if err != nil {
	// 		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
	// 	}
	// 	estilo, err := f.GetStyle(styles)
	// 	if err != nil {
	// 		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
	// 	}
	// 	fileNuevo.NewStyle(estilo)
	// 	fileNuevo.SetCellStyle(sheetName, colLetter, "1", styles)
	// 	// estilos = append(estilos, styles)
	// }

	// // Extraer valores de la primera fila
	// var primeraFila []string
	// primeraFila = append(primeraFila, filas[0]...)

	// fileNuevo.SetSheetName("Sheet1", sheetName)

	// // Escribir encabezados en el Excel
	// for colIndex, columna := range primeraFila {
	// 	columnaMin := strings.ToLower(columna)
	// 	colLetter, _ := excelize.ColumnNumberToName(colIndex + 1)
	// 	cell := colLetter + "1"
	// 	f.SetCellValue(sheetName, cell, columnaMin)
	// }

	// // err := fileNuevo.SetRowStyle(sheetName, 1, 10, estilos[1])
	// // if err != nil {
	// // 	return err
	// // }

	// // Escribir datos en el archivo Excel
	// for rowIndex, persona := range data { // row ----> 1 Gabi CABA 27
	// 	for colIndex, columna := range primeraFila {
	// 		columnaMin := strings.ToLower(columna)
	// 		colLetter, _ := excelize.ColumnNumberToName(colIndex + 1)
	// 		cell := colLetter + strconv.Itoa(rowIndex+2)
	// 		valor, ok := persona[columnaMin]
	// 		if ok {
	// 			f.SetCellValue(sheetName, cell, valor)
	// 		}
	// 		// fmt.Printf("sheetName: %s cell: %s value: %s\n", sheetName, cell, valor)
	// 	}
	// }

	// // Guardar archivo
	// if err = f.SaveAs(proceso.NombreSalida); err != nil {
	// 	ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
	// }
	_, err = db.Exec("CALL extractor.act_log_detalle($1, 'F', $2)", idLogDetalle, fmt.Sprintf("Archivo guardado en: \"%s\"", nombreSalida))
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
	}
	_, err = db.Exec("CALL extractor.etl_ending($1)", id_log)
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
	}

	return nil
}
