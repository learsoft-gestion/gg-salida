package handlers

import (
	"Nueva/modelos"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func DeleteProceso(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		idProceso := vars["idProceso"]

		// Obtengo el proceso a borrar
		querySelect := "select id_modelo, fecha_desde, fecha_hasta, num_version from extractor.ext_procesados where id_proceso = $1"
		var proceso modelos.DTOproceso

		err := db.QueryRow(querySelect, idProceso).Scan(&proceso.Id_modelo, &proceso.Fecha_desde, &proceso.Fecha_hasta, &proceso.Version)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Proceso no encontrado", http.StatusNotFound)
			} else {
				fmt.Println("Error al ejecutar query: ", err.Error())
				http.Error(w, "Error en el servidor: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Borro registro de control congelado
		if proceso.Fecha_desde == proceso.Fecha_hasta {
			queryDeleteCongelado := "delete from extractor.ext_control_congelado where id_modelo = $1 and fecha = $2 and num_version = $3"
			_, err = db.Exec(queryDeleteCongelado, proceso.Id_modelo, proceso.Fecha_desde, proceso.Version)
			if err != nil {
				fmt.Println("Error al eliminar de otra tabla: ", err.Error())
				http.Error(w, "Error en el servidor: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Finalmente borro registro de proceso
		query := "delete from extractor.ext_procesados where id_proceso = $1"

		result, err := db.Exec(query, idProceso)
		if err != nil {
			fmt.Println("Error al ejecutar query: ", err.Error())
			http.Error(w, "Error en el servidor: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if cuenta, err := result.RowsAffected(); cuenta == 1 {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			respuesta := "Proceso borrado exitosamente"
			jsonResp, _ := json.Marshal(respuesta)
			w.Write(jsonResp)
			return
		} else if cuenta == 0 {
			fmt.Println("Error al borrar proceso")
			http.Error(w, "Error al borrar proceso", http.StatusBadRequest)
		} else if err != nil {
			fmt.Println("Error al borrar proceso: " + err.Error())
			http.Error(w, "Error en el servidor: "+err.Error(), http.StatusInternalServerError)
		}
	}
}
