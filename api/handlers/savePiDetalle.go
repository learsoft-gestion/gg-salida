package handlers

import (
	"Nueva/modelos"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func SavePiDetalle(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var piDetalle modelos.PiDetalle
		err := json.NewDecoder(r.Body).Decode(&piDetalle)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, "Error decodificando JSON", http.StatusBadRequest)
			return
		}

		if r.Method == "POST" {
			query := "INSERT INTO extractor.ext_pi_detalle (cuil, fecha_ingreso, remuneracion_total, categoria, descuenta_cuota_sindical) values ($1, $2, $3, $4, $5) RETURNING cuil"

			_, err := db.Exec(query, piDetalle.Cuil, piDetalle.FechaIngreso, piDetalle.RemTotal, piDetalle.Categoria, piDetalle.DescuentaCuotaSindical)
			if err != nil {
				fmt.Println("Error al insertar nuevo detalle:", err.Error())
				http.Error(w, "Error en el servidor: "+err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			respuesta := "Detalle creado exitosamente"
			jsonResp, _ := json.Marshal(respuesta)
			w.Write(jsonResp)
		} else if r.Method == "PATCH" {
			query := "UPDATE extractor.ext_pi_detalle SET cuil = $1, fecha_ingreso = $2, remuneracion_total = $3, categoria = $4, descuenta_cuota_sindical = $5 where cuil = $6"

			result, err := db.Exec(query, piDetalle.Cuil, piDetalle.FechaIngreso, piDetalle.RemTotal, piDetalle.Categoria, piDetalle.DescuentaCuotaSindical, piDetalle.Cuil)
			if err != nil {
				fmt.Println("Error al ejecutar query: ", err.Error())
				http.Error(w, "Error en el servidor: "+err.Error(), http.StatusInternalServerError)
				return
			}

			if cuenta, err := result.RowsAffected(); cuenta == 1 {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				respuesta := "Detalle actualizado exitosamente"
				jsonResp, _ := json.Marshal(respuesta)
				w.Write(jsonResp)
				return
			} else if err != nil {
				fmt.Println("Error al actualizar detalle: " + err.Error())
				http.Error(w, "Error en el servidor: "+err.Error(), http.StatusInternalServerError)
			}
		} else if r.Method == "DELETE" {
			query := "delete from extractor.ext_pi_detalle where cuil = $1"

			result, err := db.Exec(query, piDetalle.Cuil)
			if err != nil {
				fmt.Println("Error al ejecutar query: ", err.Error())
				http.Error(w, "Error en el servidor: "+err.Error(), http.StatusInternalServerError)
				return
			}

			if cuenta, err := result.RowsAffected(); cuenta == 1 {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				respuesta := "Detalle borrado exitosamente"
				jsonResp, _ := json.Marshal(respuesta)
				w.Write(jsonResp)
				return
			} else if cuenta == 0 {
				fmt.Println("Error al borrar detalle")
				http.Error(w, "Error al borrar Detalle", http.StatusBadRequest)
			} else if err != nil {
				fmt.Println("Error al borrar Detalle: " + err.Error())
				http.Error(w, "Error en el servidor: "+err.Error(), http.StatusInternalServerError)
			}
		}
	}
}
