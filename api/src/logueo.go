package src

import (
	"database/sql"
	"fmt"
)

func Logueo(db *sql.DB, nombre string) (int, int, error) {
	// Comienza el proceso
	var id_log int
	err := db.QueryRow("SELECT extractor.etl_start()").Scan(&id_log)
	if err != nil {
		return 0, 0, err
	}
	fmt.Println("id_log: ", id_log)

	// Inicializamos registro en ext logueo detalle
	var idLogDetalle int
	err = db.QueryRow("SELECT extractor.start_log_detalle($1, $2)", id_log, nombre).Scan(&idLogDetalle)
	if err != nil {
		ManejoErrores(db, idLogDetalle, nombre, err)
	}
	return id_log, idLogDetalle, nil
}

func Procesados(db *sql.DB, id int, fecha1 string, fecha2 string, version int, cant_registros int, nombre_salida string) error {

	_, err := db.Exec("insert into extractor.ext_procesados (id_modelo, fecha_desde, fecha_hasta, version, cant_registros, nombre_salida) values ($1,$2,$3,$4,$5,$6)", id, fecha1, fecha2, version, cant_registros, nombre_salida)
	if err != nil {
		return err
	}

	return nil
}
