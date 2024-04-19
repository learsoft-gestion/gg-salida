package handlers

import (
	"Nueva/modelos"
	"database/sql"
	"encoding/json"
	"net/http"
)

var Clientes []modelos.Cliente

func GetClientes(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Clientes = nil
		query := "select razon_social, cuit from datos.clientes"
		nombre_cliente := r.URL.Query().Get("cliente")
		cuit_cliente := r.URL.Query().Get("cuit")
		if len(nombre_cliente) > 0 {
			query += " where razon_social like '%" + nombre_cliente + "%'"
		}
		if len(cuit_cliente) > 0 {
			if len(nombre_cliente) == 0 {
				query += " where cuit like '%" + cuit_cliente + "%'"
			} else {
				query += " and cuit like '%" + cuit_cliente + "%'"
			}
		}

		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for rows.Next() {
			var cliente modelos.Cliente
			if err = rows.Scan(&cliente.Nombre, &cliente.Cuit); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			Clientes = append(Clientes, cliente)

		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(Clientes); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
