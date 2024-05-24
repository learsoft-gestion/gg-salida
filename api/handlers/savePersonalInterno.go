package handlers

import (
	"Nueva/modelos"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

func SavePersonalInterno(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var personalInterno modelos.PersonalInterno
		err := json.NewDecoder(r.Body).Decode(&personalInterno)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, "Error decodificando JSON", http.StatusBadRequest)
			return
		}

		if r.Method == "POST" {
			// Valido formato de XX-XXXXXXXX-X para cuil
			regexp := regexp.MustCompile(`^(20|23|24|27)-\d{8}-\d$`)
			if personalInterno.Cuil != "" && !regexp.MatchString(personalInterno.Cuil) {
				fmt.Println("Formato de cuil inválido: " + personalInterno.Cuil)
				http.Error(w, "Debe ingresar un cuil válido.", http.StatusBadRequest)
				return
			}
			query := "INSERT INTO extractor.ext_personal_interno (cuil) values ($1)"

			_, err := db.Exec(query, personalInterno.Cuil)
			if err != nil {
				fmt.Println("Error al insertar nuevo legajo:", err.Error())
				http.Error(w, "Error en el servidor: "+err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			respuesta := "Alícuota creada exitosamente"
			jsonResp, _ := json.Marshal(respuesta)
			w.Write(jsonResp)
		} else if r.Method == "DELETE" {
			query := "delete from extractor.ext_personal_interno a where cuil = $1"

			result, err := db.Exec(query, personalInterno.Cuil)
			if err != nil {
				fmt.Println("Error al ejecutar query: ", err.Error())
				http.Error(w, "Error en el servidor: "+err.Error(), http.StatusInternalServerError)
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
				http.Error(w, "La alícuota tiene valores asociados", http.StatusBadRequest)
			} else if err != nil {
				fmt.Println("Error al borrar alícuota: " + err.Error())
				http.Error(w, "Error en el servidor: "+err.Error(), http.StatusInternalServerError)
			}
		}
	}
}
