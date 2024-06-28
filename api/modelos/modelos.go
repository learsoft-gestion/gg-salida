package modelos

import (
	"database/sql"
)

type Modelo struct {
	Id_modelo           int
	Id_empresa          int
	Id_concepto         string
	Id_convenio         int
	Id_tipo             string
	Empresa             string
	EmpReducido         string
	Concepto            string
	Convenio            string
	Tipo                string
	Nombre              string
	Filtro_convenio     string
	Filtro_personas     string
	Filtro_recibos      string
	Formato_salida      string
	Ultima_ejecucion    string
	Query               string
	Ruta_archivo_modelo string
	Archivo_modelo      string
	Vigente             string
	Ruta_archivo_nomina string
	Archivo_nomina      string
	Columna_estado      string
	Select_control      string
}

type Proceso struct {
	Id_modelo               int
	Id_empresa              int
	Nombre_empresa          string
	Nombre_empresa_reducido string
	Id_convenio             int
	Nombre_convenio         string
	Id_concepto             string
	Id_tipo                 string
	Nombre                  string
	Filtro_convenio         string
	Filtro_personas         string
	Filtro_recibos          string
	Formato_salida          string
	Query                   string
	Archivo_modelo          string
	Archivo_nomina          string
	Id_procesado            int
	Id_consultado           int
	Columna_estado          string
	Id_query                int
	Select_control          string
	Select_salida           string
}

type ModeloProyeccion struct {
	Id_modelo               int
	Id_empresa              int
	Nombre_empresa          string
	Nombre_empresa_reducido string
	Id_convenio             int
	Nombre_convenio         string
	Id_concepto             string
	Id_tipo                 string
	Nombre                  string
	Filtro_convenio         sql.NullString
	Filtro_personas         sql.NullString
	Filtro_recibos          sql.NullString
	Formato_salida          string
	Query                   string
	Archivo_modelo          string
	Archivo_nomina          string
	Id_procesado            int
	Columna_estado          string
	Id_query                int
	Select_control          string
	Select_salida           string
}

type ModeloDTO struct {
	Id_modelo int
	Convenio  string
	Empresa   string
	Concepto  string
	Nombre    string
	Tipo      string
	Fecha     string
	Total     float64
}

type Migrador struct {
	ID             int
	Procesado      bool
	FechaProcesado sql.NullTime
	ArchivoEntrada string
	ArchivoFinal   string
	Empresa        string
	Periodo        string
	Convenio       string
	Estado         sql.NullString
	Descripcion    sql.NullString
	RutaEntrada    sql.NullString
	RutaFinal      sql.NullString
}

type PatchModelo struct {
	Id      int
	Vigente bool
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
	Bloqueado        bool
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
	Ids      string
	Columnas []string
	Valores  map[string]interface{}
}

type Cabecera struct {
	Formato            string
	Separador          string
	XmlTag             string
	Tag                string
	Children           string
	Sentido_encabezado string
	Estilo             string
	Encabezados        string
	FinalDeLinea       string
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

type Alicuota struct {
	IdAlicuota  string
	IdConvenio  string
	Nombre      string
	Descripcion string
}

type ValorAlicuota struct {
	IdValoresAlicuota string
	IdAlicuota        string
	VigenciaDesde     string
	Valor             string
}

type AlicuotaBack struct {
	ValorAli   string
	NombreAli  string
	ReplaceAli string
}

type Select_control struct {
	Query sql.NullString
}

type PersonalInterno struct {
	Cuil string
}

type PiCabecera struct {
	IdPi         string
	IdEmpresaAdm string
	Periodo      string
}

type PiDetalle struct {
	IdPi                   string
	Cuil                   string
	FechaIngreso           string
	RemTotal               int
	Categoria              int
	DescuentaCuotaSindical bool
}

type Estilos struct {
	StyleMoneda            int
	StyleNumero            int
	StyleNumeroDecimal     int
	StyleEncabezadoNomina  int
	StyleEncabezadoControl int
	StyleColumnaControl    int
	StyleTotalesControl    int
	StyleControlCeleste    int
	StyleAligned           int
	StyleDefaultCabecera   int
	StyleDefault           int
	StyleVertical          int
	StyleColumnaInfo       int
	StyleValorInfo         int
}

type Control_congelado struct {
	Id_modelo    int
	Fecha        string
	Num_version  int
	Json_control []byte
}

type Bloqueado struct {
	Id_procesado int
	Bloquear     bool
}
