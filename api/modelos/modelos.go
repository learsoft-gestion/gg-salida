package modelos

type Modelo struct {
	Id_modelo        int
	Id_empresa       int
	Id_concepto      string
	Id_convenio      int
	Id_tipo          string
	Empresa          string
	Concepto         string
	Convenio         string
	Tipo             string
	Nombre           string
	Filtro_personas  string
	Filtro_recibos   string
	Formato_salida   string
	Ultima_ejecucion string
	Query            string
	Archivo_modelo   string
	Vigente          string
	Filtro_having    string
	Archivo_control  string
	Archivo_nomina   string
}

type Proceso struct {
	Id_modelo       int    `json:"id"`
	Id_empresa      int    `json:"id_empresa"`
	Nombre_empresa  string `json:"nombre_empresa"`
	Id_convenio     int    `json:"id_convenio"`
	Nombre_convenio string `json:"nombre_convenio"`
	Nombre          string `json:"nombre"`
	Filtro_convenio string `json:"filtro_convenio"`
	Filtro_personas string `json:"filtro_personas"`
	Filtro_recibos  string `json:"filtro_recibos"`
	Formato_salida  string `json:"formato_salida"`
	Query           string `json:"query"`
	Archivo_modelo  string `json:"archivo_modelo"`
	Filtro_having   string `json:"filtro_having"`
	Archivo_control string `json:"archivo_control"`
	Archivo_nomina  string `json:"archivo_nomina"`
	Id_procesado    int
}

type Empresa struct {
	Id     int    `json:"id"`
	Nombre string `json:"nombre"`
}
type Convenio struct {
	Id     int    `json:"id"`
	Nombre string `json:"nombre"`
}

type Cliente struct {
	Nombre string `json:"nombre"`
	Cuit   string `json:"cuit"`
}

type DTOselect struct {
	Empresas  []Empresa
	Convenios []Convenio
	Procesos  []DTOproceso
}

type DTOproceso struct {
	Id_modelo        int
	Convenio         string
	Empresa          string
	Concepto         string
	Nombre           string
	Tipo             string
	Fecha_desde      string
	Fecha_hasta      string
	Version          string
	Nombre_salida    string
	Ultima_ejecucion string
	Boton            string
	Nombre_control   string
	Nombre_nomina    string
	Id_procesado     int
}

type DTOdatos struct {
	Id_modelo    int
	Id_procesado int
	Fecha        string
	Fecha2       string
	Version      int
}
type DTOdatosMultiple struct {
	Id     []int
	Fecha  string
	Fecha2 string
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
	Formato            string
	Separador          string
	XmlTag             string
	Tag                string
	Children           string
	Sentido_encabezado string
	Estilo             string
}

type Campo struct {
	Titulo  string
	Nombre  string
	Inicio  int
	Fin     int
	Columna string
	Tipo    string
	Formato string
	Option1 string
	Option2 string
	Ancho   int
}

type Variable struct {
	Nombre string
	Datos  []LookupJson
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
	Filtro string `json:"filtro"`
}

type LookupJson struct {
	Id     string
	Nombre string
}

type Respuesta struct {
	Mensaje          string   `json:"mensaje"`
	Archivos_salida  []string `json:"archivos_salida"`
	Archivos_control []string `json:"archivos_control"`
	Archivos_nomina  []string `json:"archivos_nomina"`
	// Procesado       bool     `json:"procesado"`
}

type RespuestaRestantes struct {
	Mensaje string `json:"mensaje"`
	Boton   string `json:"boton"`
}

type Restantes struct {
	Id       []int
	Convenio string
	Fecha1   string
	Fecha2   string
}

type Conceptos struct {
	Conceptos []Concepto
	Tipos     []Concepto
}

type Concepto struct {
	Id     string
	Nombre string
}
