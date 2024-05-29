package src

import (
	"Nueva/modelos"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

// Ejecuta la query de ext_query y trae los registros que necesito para escribir los archivos de salida y control.
func Extractor(db, sql *sql.DB, proceso modelos.Proceso, fecha string, fecha2 string, idLogDetalle int, tipo_ejecucion string) ([]modelos.Registro, error) {

	// Reemplazo de fecha en query
	var queryFinal string
	var personal_interno string

	if proceso.Columna_estado != "" {
		queryFinal = strings.Replace(proceso.Query, "$COLUMNA_ESTADO$", proceso.Columna_estado, -1)
	} else {
		queryFinal = strings.Replace(proceso.Query, "$COLUMNA_ESTADO$", "case when p.prdmesanio=p.prdmesanio then 'ok' end", -1)

	}
	queryFinal = strings.Replace(queryFinal, "$PERIODO$", fecha, -1)
	queryFinal = strings.Replace(queryFinal, "$PERIODO2$", fecha2, -1)
	queryFinal = strings.Replace(queryFinal, "$FILTRO_CONVENIO$", proceso.Filtro_convenio, -1)
	queryFinal = strings.Replace(queryFinal, "$CONVENIO$", strconv.Itoa(proceso.Id_convenio), -1)
	queryFinal = strings.Replace(queryFinal, "$EMPRESA$", strconv.Itoa(proceso.Id_empresa), -1)
	if proceso.Filtro_personas != "" {
		queryFinal = strings.Replace(queryFinal, "$FILTRO_PERSONAS$", proceso.Filtro_personas, -1)
	} else {
		parts := strings.Split(queryFinal, "$FILTRO_PERSONAS$")
		// fmt.Println("Query: \n", queryFinal)
		queryFinal = strings.TrimSpace(parts[0]) + "\n" + strings.TrimSpace(parts[1])
	}
	if proceso.Filtro_recibos != "" {
		queryFinal = strings.Replace(queryFinal, "$FILTRO_RECIBOS$", proceso.Filtro_recibos, -1)
	} else {
		parts := strings.Split(queryFinal, "$FILTRO_RECIBOS$")
		if len(parts) == 1 && tipo_ejecucion == "salida" {
			fmt.Println("QUERY: ", queryFinal)
			// fmt.Println("Proceso Query: ", proceso.Query)
		} else {
			queryFinal = strings.TrimSpace(parts[0]) + "\n" + strings.TrimSpace(parts[1])
		}
	}
	err := db.QueryRow("select extractor.obt_tabla_pi()").Scan(&personal_interno)
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		fmt.Println("Error al obtener personal interno del extractor")
		return nil, err
	} else {
		queryFinal = strings.Replace(queryFinal, "$QUERY_PI$", personal_interno, -1)
	}

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
	// fmt.Println(columnas)

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
			// if colName == "OK" {
			// 	fmt.Printf("%v", *(valores[i].(*interface{})))
			// }
			registroMapa[colName] = *(valores[i].(*interface{}))
		}

		id := *valores[0].(*interface{})
		idString := fmt.Sprintf("%v", id)

		// if strings.ToLower(tipo_ejecucion) == "salida" {
		// 	if registroMapa["OK"] != nil {
		// 		// fmt.Println(registroMapa["OK"])
		// 		registro := modelos.Registro{
		// 			Ids:     idString,
		// 			Valores: registroMapa,
		// 		}
		// 		registros = append(registros, registro)
		// 		// } else {
		// 		// 	fmt.Println(registroMapa["OK"])
		// 	}
		// } else {
		registro := modelos.Registro{
			Ids:      idString,
			Columnas: columnas,
			Valores:  registroMapa,
		}
		registros = append(registros, registro)
		// }
	}

	return registros, nil
}

func AddToSet(slice []modelos.Option, element modelos.Option) []modelos.Option {
	for _, el := range slice {
		if el.Id == element.Id && el.Nombre == element.Nombre {
			return slice
		}
	}
	return append(slice, modelos.Option{Id: element.Id, Nombre: element.Nombre})
}
func AddToSetConceptos(slice []modelos.Concepto, element modelos.Concepto) []modelos.Concepto {
	for _, el := range slice {
		if el.Id == element.Id && el.Nombre == element.Nombre {
			return slice
		}
	}
	return append(slice, modelos.Concepto{Id: element.Id, Nombre: element.Nombre})
}

func AddToSlice(slice []string, element string) []string {
	for _, el := range slice {
		if el == element {
			return slice
		}
	}
	return append(slice, element)
}

func FormatoFecha(s string) string {
	// MM/YYYY ------> YYYYMM
	parts := strings.Split(s, "/")
	return parts[1] + parts[0]
}
