package handlers

import (
	"Nueva/conexiones"
	"Nueva/modelos"
	"Nueva/src"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func ConsultaCongelados(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Conexion al origen de datos
		sql, err := conexiones.ConectarBase("recibos", "prod", "sqlserver")
		if err != nil {
			fmt.Println("Error al intentar conectarse a la base de datos de sqlserver")
			http.Error(w, "Error interno del servidor: "+err.Error(), http.StatusBadRequest)
			return
		}

		var datos modelos.DTOdatos
		err = json.NewDecoder(r.Body).Decode(&datos)
		if err != nil {
			http.Error(w, "Error decodificando JSON: "+err.Error(), http.StatusBadRequest)
			return
		}
		datos.Fecha = src.FormatoFecha(datos.Fecha)
		datos.Fecha2 = src.FormatoFecha(datos.Fecha2)

		queryModelos := "SELECT em.id_modelo, em.id_empresa_adm, ea.razon_social as nombre_empresa, ea.reducido as nombre_empresa_reducido, c.id_convenio as id_convenio, c.nombre as nombre_convenio, em.id_concepto, em.id_tipo, em.nombre, c.filtro as filtro_convenio, em.filtro_personas, em.filtro_recibos, em.formato_salida, em.archivo_modelo, em.archivo_nomina, em.columna_estado, em.id_query, em.select_control, em.select_salida FROM extractor.ext_modelos em JOIN extractor.ext_empresas_adm ea ON em.id_empresa_adm = ea.id_empresa_adm JOIN extractor.ext_convenios c ON em.id_convenio = c.id_convenio where vigente and em.id_modelo = $1"
		// fmt.Println("Query modelos: ", queryModelos)
		rows, err := db.Query(queryModelos, datos.Id_modelo)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, "Error al ejecutar query: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer rows.Close()

		var proceso modelos.Proceso
		for rows.Next() {

			err = rows.Scan(&proceso.Id_modelo, &proceso.Id_empresa, &proceso.Nombre_empresa, &proceso.Nombre_empresa_reducido, &proceso.Id_convenio, &proceso.Nombre_convenio, &proceso.Id_concepto, &proceso.Id_tipo, &proceso.Nombre, &proceso.Filtro_convenio, &proceso.Filtro_personas, &proceso.Filtro_recibos, &proceso.Formato_salida, &proceso.Archivo_modelo, &proceso.Archivo_nomina, &proceso.Columna_estado, &proceso.Id_query, &proceso.Select_control, &proceso.Select_salida)
			if err != nil {
				fmt.Println(err.Error())
				http.Error(w, "Error al escanear proceso: "+err.Error(), http.StatusBadRequest)
				return
			}

			proceso.Id_procesado = datos.Id_procesado
		}

		// Ejecutar nomina
		respuesta_nomina := Nomina(db, sql, datos, proceso, false)

		// Ejecutar control
		respuesta_control := Control(db, sql, datos, proceso, false)

		if respuesta_nomina.Archivos_nomina == nil || respuesta_control.Archivos_control == nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonResp, _ := json.Marshal(modelos.Respuesta{Archivos_salida: nil, Archivos_nomina: respuesta_nomina.Archivos_nomina, Archivos_control: respuesta_control.Archivos_control})
			w.Write(jsonResp)
		} else if respuesta_nomina.Archivos_nomina[0] == "No se han encontrado registros" {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			respuesta := modelos.Respuesta{
				Mensaje:         "No se han encontrado registros",
				Archivos_nomina: respuesta_nomina.Archivos_nomina,
			}
			jsonResp, _ := json.Marshal(respuesta)
			w.Write(jsonResp)
		} else if respuesta_nomina.Archivos_nomina != nil && respuesta_control.Archivos_control != nil {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			respuesta := modelos.Respuesta{
				Mensaje:          "Informe generado exitosamente",
				Archivos_nomina:  respuesta_nomina.Archivos_nomina,
				Archivos_control: respuesta_control.Archivos_control,
			}
			jsonResp, _ := json.Marshal(respuesta)
			w.Write(jsonResp)
		}

	}
}
