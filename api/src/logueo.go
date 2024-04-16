package src

import (
	"database/sql"
	"fmt"
	"time"
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

func ProcesadosSalida(db *sql.DB, id_modelo int, fecha1 string, fecha2 string, version int, cant_registros int, nombre_salida string) (int, error) {

	var id_proceso int
	err := db.QueryRow("insert into extractor.ext_procesados (id_modelo, fecha_desde, fecha_hasta, num_version, cant_registros_salida, archivo_salida) values ($1,$2,$3,$4,$5,$6) returning id_proceso", id_modelo, fecha1, fecha2, version, cant_registros, nombre_salida).Scan(&id_proceso)
	if err != nil {
		return 0, err
	}
	return id_proceso, nil

}

func ProcesadosNomina(db *sql.DB, id_proceso int, id_modelo int, fecha1 string, fecha2 string, version int, cant_registros int, nombre_nomina string) error {

	fecha_completa := time.Now()
	fecha_actual := fecha_completa.Format("2006-01-02 15:04:05")

	_, err := db.Exec("update extractor.ext_procesados ep set cant_registros_nomina = $1, archivo_nomina = $2, fecha_ejecucion = $3 where ep.id_proceso = $4", cant_registros, nombre_nomina, fecha_actual, id_proceso)
	if err != nil {
		return err
	}

	return nil
}
