package handlers

import (
	"Nueva/modelos"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func SavePiCabecera(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var piCabecera modelos.PiCabecera
		err := json.NewDecoder(r.Body).Decode(&piCabecera)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, "Error decodificando JSON", http.StatusBadRequest)
			return
		}

		if r.Method == "POST" {
			query := "INSERT INTO extractor.ext_pi_cabecera (id_empresa_adm, periodo) values ($1, $2) RETURNING id_pi"

			result := db.QueryRow(query, piCabecera.IdEmpresaAdm, piCabecera.Periodo)
			var lastInsertID int
			err := result.Scan(&lastInsertID)
			if err != nil {
				fmt.Println("Error al obtener el último ID insertado:", err.Error())
				http.Error(w, "Error en el servidor: "+err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			respuesta := struct {
				Mensaje string `json:"mensaje"`
				ID      int    `json:"id"`
			}{
				Mensaje: "Cabecera creada exitosamente",
				ID:      lastInsertID,
			}
			jsonResp, _ := json.Marshal(respuesta)
			w.Write(jsonResp)
		} else if r.Method == "PATCH" {
			query := "UPDATE extractor.ext_pi_cabecera SET id_empresa_adm = $1, periodo = $2 where id_pi = $3"

			result, err := db.Exec(query, piCabecera.IdEmpresaAdm, piCabecera.Periodo, piCabecera.IdPi)
			if err != nil {
				fmt.Println("Error al ejecutar query: ", err.Error())
				http.Error(w, "Error en el servidor: "+err.Error(), http.StatusInternalServerError)
				return
			}

			if cuenta, err := result.RowsAffected(); cuenta == 1 {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				respuesta := "Cabecera actualizada exitosamente"
				jsonResp, _ := json.Marshal(respuesta)
				w.Write(jsonResp)
				return
			} else if err != nil {
				fmt.Println("Error al actualizar alícuota: " + err.Error())
				http.Error(w, "Error en el servidor: "+err.Error(), http.StatusInternalServerError)
			}
		} else if r.Method == "DELETE" {
			query := "delete from extractor.ext_pi_cabecera c where c.id_pi = $1 and not exist (select 1 from extractor.ext_pi_detalle d where d.id_pi = c.id_pi)"

			result, err := db.Exec(query, piCabecera.IdPi)
			if err != nil {
				fmt.Println("Error al ejecutar query: ", err.Error())
				http.Error(w, "Error en el servidor: "+err.Error(), http.StatusInternalServerError)
				return
			}

			if cuenta, err := result.RowsAffected(); cuenta == 1 {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				respuesta := "Cabecera borrada exitosamente"
				jsonResp, _ := json.Marshal(respuesta)
				w.Write(jsonResp)
				return
			} else if cuenta == 0 {
				fmt.Println("Error al borrar alícuota: Tiene valores asociados")
				http.Error(w, "La cabecera tiene valores asociados", http.StatusBadRequest)
			} else if err != nil {
				fmt.Println("Error al borrar cabecera: " + err.Error())
				http.Error(w, "Error en el servidor: "+err.Error(), http.StatusInternalServerError)
			}
		}
	}
}
