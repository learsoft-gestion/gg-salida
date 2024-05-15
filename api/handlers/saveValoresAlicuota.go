package handlers

import (
	"Nueva/modelos"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
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

		// Valido formato numérico para valor
		if _, err := strconv.ParseFloat(valor.Valor, 64); valor.Valor != "" && err != nil {
			fmt.Println("Valor debe ser dato numérico")
			http.Error(w, "Debe ingresar un dato numérico como Valor.", http.StatusBadRequest)
			return
		}

		// Valido formato de fecha YYYYMM para vigenciaDesde
		regexp := regexp.MustCompile(`^\d{6}$`)
		if valor.VigenciaDesde != "" && !regexp.MatchString(valor.VigenciaDesde) {
			fmt.Println("Formato de fecha inválido: " + valor.VigenciaDesde)
			http.Error(w, "Debe ingresar formato de fecha válido.", http.StatusBadRequest)
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
			query := "INSERT INTO extractor.ext_valores_alicuotas (id_alicuota, vigencia_desde, valor) values ($1, $2, $3) RETURNING id_valores_alicuota"

			result := db.QueryRow(query, valor.IdAlicuota, valor.VigenciaDesde, valor.Valor)
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
				Mensaje: "Valor creado exitosamente",
				ID:      lastInsertID,
			}
			jsonResp, _ := json.Marshal(respuesta)
			w.Write(jsonResp)
		} else if r.Method == "DELETE" {
			query := "delete from extractor.ext_valores_alicuotas where id_valores_alicuota = $1"

			result, err := db.Exec(query, valor.IdValoresAlicuota)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Println("Error al ejecutar query: ", err.Error())
				http.Error(w, "Error en el servidor", http.StatusBadRequest)
				return
			}

			if cuenta, err := result.RowsAffected(); cuenta == 1 {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				respuesta := "Valor borrado exitosamente"
				jsonResp, _ := json.Marshal(respuesta)
				w.Write(jsonResp)
				return
			} else if err != nil {
				fmt.Println("Error al borrar valor: " + err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				http.Error(w, "Error en el servidor", http.StatusInternalServerError)
			}
		}
	}
}
