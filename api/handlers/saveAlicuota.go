package handlers

import (
	"Nueva/modelos"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func SaveAlicuota(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var alicuota modelos.Alicuota
		err := json.NewDecoder(r.Body).Decode(&alicuota)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, "Error decodificando JSON", http.StatusBadRequest)
			return
		}

		if r.Method == "PATCH" {
			query := "UPDATE extractor.ext_alicuotas SET nombre = $1, descripcion = $2 where id_alicuota = $3"

			result, err := db.Exec(query, alicuota.Nombre, alicuota.Descripcion, alicuota.IdAlicuota)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("Error al ejecutar query: ", err.Error())
				http.Error(w, "Error en el servidor", http.StatusBadRequest)
				return
			}

			if cuenta, err := result.RowsAffected(); cuenta == 1 {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				respuesta := "Alícuota actualizada exitosamente"
				jsonResp, _ := json.Marshal(respuesta)
				w.Write(jsonResp)
				return
			} else if err != nil {
				fmt.Println("Error al actualizar alícuota: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				http.Error(w, "Error en el servidor", http.StatusInternalServerError)
			}
		} else if r.Method == "POST" {
			query := "INSERT INTO extractor.ext_alicuotas (id_convenio, nombre, descripcion) values ($1, $2, $3) RETURNING id_alicuota"

			result := db.QueryRow(query, alicuota.IdConvenio, alicuota.Nombre, alicuota.Descripcion)
			var lastInsertID int
			err := result.Scan(&lastInsertID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("Error al obtener el último ID insertado:", err.Error())
				http.Error(w, "Error en el servidor", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			respuesta := struct {
				Mensaje string `json:"mensaje"`
				ID      int    `json:"id"`
			}{
				Mensaje: "Alícuota creada exitosamente",
				ID:      lastInsertID,
			}
			jsonResp, _ := json.Marshal(respuesta)
			w.Write(jsonResp)
		} else if r.Method == "DELETE" {
			query := "delete from extractor.ext_alicuotas a where a.id_alicuota = $1 and not exists (select 1 from extractor.ext_valores_alicuotas va where a.id_alicuota = va.id_alicuota)"

			result, err := db.Exec(query, alicuota.IdAlicuota)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("Error al ejecutar query: ", err.Error())
				http.Error(w, "Error en el servidor", http.StatusBadRequest)
				return
			}

			if cuenta, err := result.RowsAffected(); cuenta == 1 {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				respuesta := "Alícuota borrada exitosamente"
				jsonResp, _ := json.Marshal(respuesta)
				w.Write(jsonResp)
				return
			} else if cuenta == 0 {
				fmt.Println("Error al borrar alícuota: Tiene valores asociados")
				w.WriteHeader(http.StatusInternalServerError)
				http.Error(w, "La alícuota tiene valores asociados", http.StatusBadRequest)
			} else if err != nil {
				fmt.Println("Error al borrar alícuota: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				http.Error(w, "Error en el servidor", http.StatusInternalServerError)
			}
		}
	}
}
