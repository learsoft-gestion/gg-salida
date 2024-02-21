package src

import (
	"database/sql"
	"fmt"
)

func ManejoErrores(postgresDb *sql.DB, idLogDetalle int, nombre string, err error) {
	text := "Fallo en: " + err.Error()
	postgresDb.Exec("CALL extractor.act_log_detalle($1, 'E', $2)", idLogDetalle, text)
	fmt.Printf("Error en %s: %s \n", nombre, text)
}
