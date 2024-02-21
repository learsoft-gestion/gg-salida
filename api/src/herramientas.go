package src

import (
	"Nueva/modelos"
	"database/sql"
	"fmt"
	"strings"
)

func Extractor(db, sql *sql.DB, proceso modelos.Proceso, datos modelos.DTOdatos, idLogDetalle int) ([]modelos.Registro, error) {

	// Reemplazo de fecha en query
	var queryFinal string
	if proceso.Cant_fechas > 1 {
		query := strings.Replace(proceso.Query, "$1", datos.FechaDesde, 1)
		queryFinal = strings.Replace(query, "$2", datos.FechaHasta, 1)
	} else {
		queryFinal = strings.Replace(proceso.Query, "$1", datos.FechaDesde, 1)
	}

	// Ejecucion de query y lectura de resultados
	rows, err := sql.Query(queryFinal)
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
	}
	defer rows.Close()

	columnas, err := rows.Columns()
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
	}
	columnasNum := len(columnas)

	var registros []modelos.Registro

	valores := make([]interface{}, columnasNum)
	for i := range valores {
		valores[i] = new(interface{})
	}

	for rows.Next() {

		if err := rows.Scan(valores...); err != nil {
			ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		}

		registroMapa := make(map[string]interface{})
		for i, colNombre := range columnas {
			registroMapa[colNombre] = *(valores[i].(*interface{}))
		}

		id := *valores[0].(*interface{})
		idString := fmt.Sprintf("%v", id)

		registro := modelos.Registro{
			Ids:     idString,
			Valores: registroMapa,
		}
		registros = append(registros, registro)
	}
	fmt.Println("Registro: ", registros[0])
	return registros, nil
}
