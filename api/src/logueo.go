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

	fecha_completa := time.Now()
	fecha_actual := fecha_completa.Format("2006-01-02 15:04:05")

	// if id_proceso > 0 {
	// 	if procesado_salida {
	// 		// Ya hay archivo de salida en este id_proceso
	// 		_, err := db.Exec("insert into extractor.ext_procesados (id_modelo, fecha_desde, fecha_hasta, version, cant_registros_salida, archivo_salida, fecha_ejecucion_salida) values ($1,$2,$3,$4,$5,$6,$7)", id_modelo, fecha1, fecha2, version, cant_registros, nombre_salida, fecha_actual)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	} else {
	// 		// Hay un proceso para el mismo modelo y periodo pero no tiene archivo de control
	// 		_, err := db.Exec("update extractor.ext_procesados ep set cant_registros_control = $1, archivo_control = $2, fecha_ejecucion_control = $3 where ep.id_proceso = $4", cant_registros, nombre_salida, fecha_actual, id_proceso)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// } else {
	var id_proceso int
	err := db.QueryRow("insert into extractor.ext_procesados (id_modelo, fecha_desde, fecha_hasta, version, cant_registros_salida, archivo_salida, fecha_ejecucion_salida) values ($1,$2,$3,$4,$5,$6,$7) returning id_proceso", id_modelo, fecha1, fecha2, version, cant_registros, nombre_salida, fecha_actual).Scan(&id_proceso)
	if err != nil {
		return 0, err
	}
	return id_proceso, nil
	// }

	// return nil
}

func ProcesadosControl(db *sql.DB, id_proceso int, id_modelo int, fecha1 string, fecha2 string, version int, cant_registros int, nombre_control string) error {

	fecha_completa := time.Now()
	fecha_actual := fecha_completa.Format("2006-01-02 15:04:05")

	// Hay un proceso para el mismo modelo y periodo pero no tiene archivo de control
	// query := fmt.Sprintf("update extractor.ext_procesados ep set cant_registros_control = %v, archivo_control = %s, fecha_ejecucion_control = %s where ep.id_proceso = %v", cant_registros, nombre_control, fecha_actual, id_proceso)
	// fmt.Println(query)
	res, err := db.Exec("update extractor.ext_procesados ep set cant_registros_control = $1, archivo_control = $2, fecha_ejecucion_control = $3 where ep.id_proceso = $4", cant_registros, nombre_control, fecha_actual, id_proceso)
	if err != nil {
		return err
	} else {
		num, _ := res.RowsAffected()
		fmt.Printf("Actualizó %d registros", num)
	}
	// if id_proceso > 0 {
	// 	if procesado_control {
	// 		// Ya hay archivo de control en este id_proceso
	// 		_, err := db.Exec("insert into extractor.ext_procesados (id_modelo, fecha_desde, fecha_hasta, version, cant_registros_control, archivo_control, fecha_ejecucion_control) values ($1,$2,$3,$4,$5,$6,$7)", id_modelo, fecha1, fecha2, version, cant_registros, nombre_control, fecha_actual)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	} else {
	// 	}
	// } else {
	// 	_, err := db.Exec("insert into extractor.ext_procesados (id_modelo, fecha_desde, fecha_hasta, version, cant_registros_control, archivo_control, fecha_ejecucion_control) values ($1,$2,$3,$4,$5,$6,$7)", id_modelo, fecha1, fecha2, version, cant_registros, nombre_control, fecha_actual)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}
