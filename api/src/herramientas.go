package src

import (
	"Nueva/modelos"
	"database/sql"
	"fmt"
	"strings"
)

func Extractor(db, sql *sql.DB, proceso modelos.Proceso, datos modelos.DTOdatos, idLogDetalle int) ([]modelos.Registro, error) {

	// Reemplazo de fecha en query
	queryFinal := strings.Replace(proceso.Query, "$PERIODO$", datos.Fecha, 1)

	if proceso.Filtro_personas != "" {
		queryFinal = strings.Replace(queryFinal, "$FILTRO_PERSONAS$", proceso.Filtro_personas, 1)
	} else {
		queryFinal = strings.Replace(queryFinal, "$FILTRO_PERSONAS$\n", "", -1)
	}
	if proceso.Filtro_recibos != "" {
		queryFinal = strings.Replace(queryFinal, "$FILTRO_RECIBOS$", proceso.Filtro_recibos, 1)
	} else {
		parts := strings.Split(queryFinal, "$FILTRO_RECIBOS$")
		queryFinal = strings.TrimSpace(parts[0]) + "\n" + strings.TrimSpace(parts[1])
		// queryFinal = strings.Replace(queryFinal, "$FILTRO_RECIBOS$", "", -1)
	}

	// queryFinal = strings.TrimFunc()
	// fmt.Println("Query: \n", queryFinal)

	// Ejecucion de query y lectura de resultados
	rows, err := sql.Query(queryFinal)
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		fmt.Println("Error al hacer la query del extractor")
		return nil, err
	}
	defer rows.Close()

	columnas, err := rows.Columns()
	if err != nil {
		fmt.Println("Error al hacer la query del extractor")
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return nil, err
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
			return nil, err
		}

		registroMapa := make(map[string]interface{})
		for i, colNombre := range columnas {
			colName := strings.ToUpper(colNombre)
			registroMapa[colName] = *(valores[i].(*interface{}))
		}

		id := *valores[0].(*interface{})
		idString := fmt.Sprintf("%v", id)

		registro := modelos.Registro{
			Ids:     idString,
			Valores: registroMapa,
		}
		registros = append(registros, registro)
	}
	// fmt.Println("Registro: ", registros[0])
	return registros, nil
}
