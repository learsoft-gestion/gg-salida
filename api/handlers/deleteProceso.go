package handlers

import (
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
