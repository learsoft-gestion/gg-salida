package handlers

import (
	"Nueva/modelos"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func Migrador(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query, err := getQueryExtArchivos(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, "Error al mover el archivo "+err.Error(), http.StatusInternalServerError)
			return
		}

		var procesos []modelos.Migrador

		for rows.Next() {
			var p modelos.Migrador
			err := rows.Scan(&p.ID, &p.Procesado, &p.FechaProcesado, &p.ArchivoEntrada, &p.ArchivoFinal,
				&p.Empresa, &p.Periodo, &p.Convenio, &p.Estado, &p.Descripcion, &p.RutaEntrada, &p.RutaFinal)
			if err != nil {
				http.Error(w, "Error al escanear fila "+err.Error(), http.StatusInternalServerError)
				return
			}
			procesos = append(procesos, p)
		}

		jsonData, err := json.Marshal(procesos)
		if err != nil {
			http.Error(w, "Error al convertir a JSON "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	}
}

func getQueryExtArchivos(r *http.Request) (string, error) {
	query := `SELECT * FROM extractor.ext_archivos WHERE 1 = 1 `

	idDesde := r.FormValue("idDesde")
	idHasta := r.FormValue("idHasta")
	procesado := r.FormValue("procesado")
	fechaDesde := r.FormValue("fechaDesde")
	fechaHasta := r.FormValue("fechaHasta")
	descripcion := r.FormValue("descripcion")

	if idDesde != "" && idHasta != "" {
		idDesde, err1 := strconv.Atoi(idDesde)
		idHasta, err2 := strconv.Atoi(idHasta)
		if err1 != nil || err2 != nil {
			return "", fmt.Errorf("error al parsear los valores 'idDesde' y 'idHasta'")
		}

		query += fmt.Sprintf(` and id_numero between %d and %d `, idDesde, idHasta)
	} else if idDesde != "" || idHasta != "" {
		if idDesde != "" {
			idDesde, _ := strconv.Atoi(idDesde)
			query += fmt.Sprintf(` and id_numero >= %d `, idDesde)
		} else {
			idHasta, _ := strconv.Atoi(idHasta)
			query += fmt.Sprintf(` and id_numero <= %d `, idHasta)
		}
	}
	if procesado != "" {
		procesado, err := strconv.ParseBool(procesado)
		if err != nil {
			return "", fmt.Errorf("error al parsear el valor de 'procesado': %s", err.Error())
		}

		query += fmt.Sprintf(` and procesado = %t `, procesado)
	}
	if fechaDesde != "" && fechaHasta != "" {
		query += fmt.Sprintf(` and fecha_procesado between to_timestamp('%s', 'DD-MM-YYYY') and to_timestamp('%s', 'DD-MM-YYYY') + INTERVAL '1 day' `, fechaDesde, fechaHasta)
	} else if fechaDesde != "" || fechaHasta != "" {
		if fechaDesde != "" {
			query += fmt.Sprintf(` and fecha_procesado >= to_timestamp('%s', 'DD-MM-YYYY') `, fechaDesde)
		} else {
			query += fmt.Sprintf(` and fecha_procesado <= to_timestamp('%s', 'DD-MM-YYYY') + INTERVAL '1 day' `, fechaHasta)
		}
	}
	if descripcion != "" {
		descripcion, err := strconv.ParseBool(descripcion)
		if err != nil {
			return "", fmt.Errorf("error al parsear el valor de 'descripcion': %s", err.Error())
		}
		if descripcion {
			query += ` and descripcion = 'Procesado correctamente' `
		} else {
			query += ` and descripcion <> 'Procesado correctamente' `
		}
	}
	query += ` order by 1`

	return query, nil
}
