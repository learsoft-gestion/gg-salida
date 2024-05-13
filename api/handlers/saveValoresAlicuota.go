package handlers

import (
	"Nueva/modelos"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func SaveValoresAlicuota(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var valor modelos.ValorAlicuota
		err := json.NewDecoder(r.Body).Decode(&valor)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, "Error decodificando JSON", http.StatusBadRequest)
			return
		}

		if r.Method == "PATCH" {
			query := "UPDATE extractor.ext_valores_alicuotas SET vigencia_desde = $1, valor = $2 where id_valores_alicuota = $3"

			result, err := db.Exec(query, valor.VigenciaDesde, valor.Valor, valor.IdValoresAlicuota)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("Error al ejecutar query: ", err.Error())
				http.Error(w, "Error en el servidor", http.StatusBadRequest)
				return
			}

			if cuenta, err := result.RowsAffected(); cuenta == 1 {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				respuesta := "Valor actualizado exitosamente"
				jsonResp, _ := json.Marshal(respuesta)
				w.Write(jsonResp)
				return
			} else if err != nil {
				fmt.Println("Error al actualizar valor: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				http.Error(w, "Error en el servidor", http.StatusInternalServerError)
			}
		} else if r.Method == "POST" {
			query := "INSERT INTO extractor.ext_valores_alicuotas (id_alicuota, vigencia_desde, valor) values ($1, $2, $3)"

			result, err := db.Exec(query, valor.IdAlicuota, valor.VigenciaDesde, valor.Valor)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("Error al ejecutar query: ", err.Error())
				http.Error(w, "Error en el servidor", http.StatusBadRequest)
				return
			}

			if cuenta, err := result.RowsAffected(); cuenta == 1 {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				respuesta := "Valor creado exitosamente"
				jsonResp, _ := json.Marshal(respuesta)
				w.Write(jsonResp)
				return
			} else if err != nil {
				fmt.Println("Error al crear valor: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				http.Error(w, "Error en el servidor", http.StatusInternalServerError)
			}
		}
	}
}
