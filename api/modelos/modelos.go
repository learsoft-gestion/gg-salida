package modelos

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

type DTOdatos struct {
	Id         int
	FechaDesde string
	FechaHasta string
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
