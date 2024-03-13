package modelos

import "database/sql"

type Proceso struct {
	Id              int    `json:"id"`
	Id_empresa      int    `json:"id_empresa"`
	Nombre          string `json:"nombre"`
	Filtro_convenio string `json:"filtro_convenio"`
	Filtro_personas string `json:"filtro_personas"`
	Filtro_recibos  string `json:"filtro_recibos"`
	Formato_salida  string `json:"formato_salida"`
	Query           string `json:"query"`
	Archivo_modelo  string `json:"archivo_modelo"`
}

type Empresa struct {
	Id     int    `json:"id"`
	Nombre string `json:"nombre"`
}
type Convenio struct {
	Id     int    `json:"id"`
	Nombre string `json:"nombre"`
}

type DTOselect struct {
	Empresas  []Empresa
	Convenios []Convenio
	Procesos  []DTOproceso
}

type DTOproceso struct {
	Id               int
	Convenio         string
	Empresa          string
	Concepto         string
	Nombre           string
	Tipo             string
	Fecha_desde      sql.NullString
	Fecha_hasta      sql.NullString
	Nombre_salida    sql.NullString
	Ultima_ejecucion sql.NullString
}

type DTOdatos struct {
	IDs     []int
	Fecha   string
	Fecha2  string
	Forzado bool
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

type Cabecera struct {
	Formato   string
	Separador string
}

type Campo struct {
	Titulo  string
	Nombre  string
	Inicio  int
	Fin     int
	Columna string
	Tipo    string
	Formato string
}

type Variable struct {
	Nombre string
	Datos  []Option
}

type Plantilla struct {
	Cabecera  Cabecera   `json:"cabecera"`
	Campos    []Campo    `json:"campo"`
	Variables []Variable `json:"variables"`
}

type ErrorFormateado struct {
	Mensaje   string
	Procesado bool
}

type Option struct {
	Id     int    `json:"id"`
	Nombre string `json:"nombre"`
}

type Respuesta struct {
	Mensaje         string   `json:"mensaje"`
	Archivos_salida []string `json:"archivos_salida"`
	Procesado       bool     `json:"procesado"`
}

type Conceptos struct {
	Conceptos []Concepto
	Tipos     []Concepto
}

type Concepto struct {
	Id     string
	Nombre string
}
