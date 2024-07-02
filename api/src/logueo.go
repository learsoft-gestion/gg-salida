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

	if nombre_salida == "Error" {
		_, err = ProcesadosNomina(db, id_proceso, cant_registros, "Error", 0, "", "")
		if err != nil {
			return 0, err
		}
		err = ProcesadosControl(db, id_proceso, "Error")
		if err != nil {
			return 0, err
		}
	}

	return id_proceso, nil

}

func ProcesadosNomina(db *sql.DB, id_proceso int, cant_registros int, nombre_nomina string, id_modelo int, fecha1, fecha2 string) (int, error) {

	var fecha_completa time.Time

	location, err := time.LoadLocation("America/Argentina/Buenos_Aires")
	if err != nil {
		fmt.Println("Error al encontrar el timezone")
		fecha_completa = time.Now()
	} else {
		fecha_completa = time.Now().In(location)
	}
	fecha_actual := fecha_completa.Format("2006-01-02 15:04:05")

	if id_modelo > 0 {
		var id_consulta int
		err = db.QueryRow("insert into extractor.ext_procesados (id_modelo, fecha_desde, fecha_hasta, archivo_nomina, fecha_ejecucion) values ($1,$2,$3,$4,$5) returning id_proceso", id_modelo, fecha1, fecha2, nombre_nomina, fecha_actual).Scan(&id_consulta)
		if err != nil {
			return 0, err
		}

		if nombre_nomina == "Error" {
			_, err = db.Exec("insert into extractor.ext_procesados (id_modelo, fecha_desde, fecha_hasta, archivo_nomina, fecha_ejecucion) values ($1,$2,$3,$4,$5) returning id_proceso", id_modelo, fecha1, fecha2, "Error", fecha_actual)
			if err != nil {
				return 0, err
			}
			err = ProcesadosControl(db, id_consulta, "Error")
			if err != nil {
				return 0, err
			}
		}

		return id_consulta, nil
	}

	res, err := db.Exec("update extractor.ext_procesados ep set cant_registros_nomina = $1, archivo_nomina = $2, fecha_ejecucion = $3 where ep.id_proceso = $4", cant_registros, nombre_nomina, fecha_actual, id_proceso)
	if err != nil {
		return 0, err
	} else {
		actualizados, err := res.RowsAffected()
		if err != nil {
			return 0, err
		} else if actualizados < 1 {
			fmt.Println("Id_procesado: ", id_proceso)
			return 0, fmt.Errorf("no actualizó ningun registro en ext_procesados para el procesamiento de nomina")
		}
	}

	if nombre_nomina == "Error" {
		_, err = db.Exec("update extractor.ext_procesados ep set archivo_salida = $1 where ep.id_proceso = $2", "Error", id_proceso)
		if err != nil {
			return 0, err
		}
		err = ProcesadosControl(db, id_proceso, "Error")
		if err != nil {
			return 0, err
		}
	}

	return 0, nil
}

func ProcesadosControl(db *sql.DB, id_proceso int, nombre_control string) error {

	_, err := db.Exec("update extractor.ext_procesados ep set archivo_control = $1 where ep.id_proceso = $2", nombre_control, id_proceso)
	if err != nil {
		return err
	}

	if nombre_control == "Error" {
		_, err = db.Exec("update extractor.ext_procesados ep set archivo_salida = $1, archivo_nomina = $2 where ep.id_proceso = $3", "Error", "Error", id_proceso)
		if err != nil {
			return err
		}
	}

	return nil
}

// func ConsultadosNomina(db *sql.DB, id_modelo int, fecha1 string, fecha2 string, nombre_nomina string) (int, error) {

// 	var fecha_completa time.Time

// 	location, err := time.LoadLocation("America/Argentina/Buenos_Aires")
// 	if err != nil {
// 		fmt.Println("Error al encontrar el timezone")
// 		fecha_completa = time.Now()
// 	} else {
// 		fecha_completa = time.Now().In(location)
// 	}
// 	fecha_actual := fecha_completa.Format("2006-01-02 15:04:05")

// 	var id_consulta int
// 	err = db.QueryRow("insert into extractor.ext_consultados (id_modelo, fecha_desde, fecha_hasta, archivo_nomina, fecha_ejecucion) values ($1,$2,$3,$4,$5) returning id_consulta", id_modelo, fecha1, fecha2, nombre_nomina, fecha_actual).Scan(&id_consulta)
// 	if err != nil {
// 		return 0, err
// 	}

// 	if nombre_nomina == "Error" {
// 		_, err = db.Exec("insert into extractor.ext_consultados (id_modelo, fecha_desde, fecha_hasta, archivo_nomina, fecha_ejecucion) values ($1,$2,$3,$4,$5) returning id_consulta", id_modelo, fecha1, fecha2, "Error", fecha_actual)
// 		if err != nil {
// 			return 0, err
// 		}
// 		err = ConsultadosControl(db, id_consulta, "Error")
// 		if err != nil {
// 			return 0, err
// 		}
// 	}

// 	return id_consulta, nil
// }

// func ConsultadosControl(db *sql.DB, id_consulta int, nombre_control string) error {

// 	_, err := db.Exec("update extractor.ext_consultados ep set archivo_control = $1 where ep.id_consulta = $2", nombre_control, id_consulta)
// 	if err != nil {
// 		return err
// 	}

// 	if nombre_control == "Error" {
// 		_, err = db.Exec("update extractor.ext_procesados ep set archivo_nomina = $1 where ep.id_consulta = $3", "Error", id_consulta)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
