package modelos

type Proceso struct {
	Id              int    `json:"id"`
	Nombre          string `json:"nombre"`
	Filtro_personas string `json:"filtro_personas"`
	Filtro_recibos  string `json:"filtro_recibos"`
	Formato_salida  string `json:"formato_salida"`
	Query           string `json:"query"`
}

type DTOproceso struct {
	Id     int    `json:"id"`
	Nombre string `json:"nombre"`
	// Cant_fechas int    `json:"cant_fechas"`
}

type DTOdatos struct {
	Id    int
	Fecha string
}

type ProcesosTemplate struct {
	Procesos []DTOproceso
	Title    string
}

type Page struct {
	Title string
	Body  []byte
}

// Modelos para lectura de tabla
type Registro struct {
	Ids     string
	Valores map[string]interface{}
}

type Campo struct {
	Nombre  string
	Inicio  int
	Fin     int
	Tipo    string
	Formato string
}

type Plantilla struct {
	Campos []Campo `json:"campo"`
}
