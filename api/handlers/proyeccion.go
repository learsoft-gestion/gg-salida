package handlers

import (
	"Nueva/conexiones"
	"Nueva/modelos"
	"Nueva/src"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
)

func Proyeccion(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Conexion al origen de datos
		sql, err := conexiones.ConectarBase("recibos", "prod", "sqlserver")
		if err != nil {
			http.Error(w, "Error al conectarse a sqlServer en Proyeccion"+err.Error(), http.StatusInternalServerError)
			return
		}

		mes := r.URL.Query().Get("mes")
		mes = src.FormatoFecha(mes)

		rows, err := db.Query("select m.id_modelo from extractor.ext_modelos m left outer join extractor.ext_totales et on m.id_modelo = et.id_modelo and et.fecha = $1 where vigente and et.id_totales is null", mes)
		if err != nil {
			http.Error(w, "Error al ejecutar select en Proyeccion", http.StatusInternalServerError)
			return
		}

		var ids []int
		for rows.Next() {
			var id int
			if err = rows.Scan(&id); err != nil {
				http.Error(w, "Error al escanear datos en Proyeccion", http.StatusInternalServerError)
				return
			}
			ids = append(ids, id)
		}

		var placeholders []string
		for i := range ids {
			placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		}

		queryModelos := fmt.Sprintf("SELECT em.id_modelo, em.id_empresa_adm, ea.razon_social as nombre_empresa, ea.reducido as nombre_empresa_reducido, c.id_convenio as id_convenio, c.nombre as nombre_convenio, em.id_concepto, em.id_tipo, em.nombre, c.filtro as filtro_convenio, em.filtro_personas, em.filtro_recibos, em.formato_salida, em.archivo_modelo, em.archivo_nomina, em.columna_estado, em.id_query, em.select_control, em.select_salida FROM extractor.ext_modelos em JOIN extractor.ext_empresas_adm ea ON em.id_empresa_adm = ea.id_empresa_adm JOIN extractor.ext_convenios c ON em.id_convenio = c.id_convenio where vigente and em.id_modelo in (%s)", strings.Join(placeholders, ","))

		stmt, err := db.Prepare(queryModelos)
		if err != nil {
			http.Error(w, "Error al preparar query. "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer stmt.Close()
		var args []interface{}
		for _, arg := range ids {
			args = append(args, arg)
		}

		filas, err := stmt.Query(args...)
		if err != nil {
			http.Error(w, "Error al ejecutar el query: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer filas.Close()

		var modelosSlice []modelos.ModeloProyeccion
		for filas.Next() {
			var modelo modelos.ModeloProyeccion

			err = filas.Scan(&modelo.Id_modelo, &modelo.Id_empresa, &modelo.Nombre_empresa, &modelo.Nombre_empresa_reducido, &modelo.Id_convenio, &modelo.Nombre_convenio, &modelo.Id_concepto, &modelo.Id_tipo, &modelo.Nombre, &modelo.Filtro_convenio, &modelo.Filtro_personas, &modelo.Filtro_recibos, &modelo.Formato_salida, &modelo.Archivo_modelo, &modelo.Archivo_nomina, &modelo.Columna_estado, &modelo.Id_query, &modelo.Select_control, &modelo.Select_salida)
			if err != nil {
				fmt.Println(err.Error())
				http.Error(w, "Error al escanear modelo: "+err.Error(), http.StatusBadRequest)
				return
			}
			modelosSlice = append(modelosSlice, modelo)
		}

		for _, modelo := range modelosSlice {

			id_log, idLogDetalle, err := src.Logueo(db, modelo.Nombre)
			if err != nil {
				http.Error(w, "Error al loguear en Proyeccion", http.StatusInternalServerError)
				return
			}

			// Reemplazar Alicuotas
			var select_control modelos.Select_control
			err = db.QueryRow("SELECT extractor.obt_control($1, $2, $3)", modelo.Id_modelo, mes, mes).Scan(&select_control.Query)
			if err != nil {
				http.Error(w, "Error en el query de obt_control en Proyeccion", http.StatusInternalServerError)
				return
			}
			if select_control.Query.Valid {
				modelo.Select_control = select_control.Query.String
			}

			// Busco y guardo el query grande
			var query string
			db.QueryRow("SELECT texto_query FROM extractor.ext_query where id_query = $1", modelo.Id_query).Scan(&query)
			modelo.Query = strings.Replace(query, "$SELECT$", modelo.Select_control, 1)

			// Convierto el modelo.ModelosProyeccion en un modelos.Proceso
			modeloFinal := modelos.Proceso{Id_modelo: modelo.Id_modelo, Query: modelo.Query, Columna_estado: modelo.Columna_estado, Filtro_convenio: modelo.Filtro_convenio.String, Filtro_personas: modelo.Filtro_personas.String, Filtro_recibos: modelo.Filtro_recibos.String, Nombre: modelo.Nombre, Id_convenio: modelo.Id_convenio, Id_empresa: modelo.Id_empresa}

			var registros []modelos.Registro
			registros, err = src.Extractor(db, sql, modeloFinal, mes, mes, idLogDetalle, "control")
			if err != nil {
				http.Error(w, "Error en el Extractor en Proyeccion. "+err.Error(), http.StatusInternalServerError)
				return
			}

			var columna string
			for _, col := range registros[0].Columnas {
				if strings.Contains(col, "Total a pagar") {
					columna = col
				}
			}
			value := registros[0].Valores[strings.ToUpper(columna)]

			_, err = db.Exec("INSERT INTO extractor.ext_totales (fecha, id_modelo, valor) VALUES($1, $2, $3);", mes, modeloFinal.Id_modelo, value)
			if err != nil {
				http.Error(w, "Error al insertar valor en ext_totales en Proyeccion. "+err.Error(), http.StatusInternalServerError)
				return
			}

			// Logueo
			_, err = db.Exec("CALL extractor.act_log_detalle($1, 'F', $2)", idLogDetalle, "Valor de proyeccion insertado en ext_totales")
			if err != nil {
				http.Error(w, "Error en el logueo en Proyeccion. "+err.Error(), http.StatusInternalServerError)
				return
			}
			_, err = db.Exec("CALL extractor.etl_ending($1)", id_log)
			if err != nil {
				http.Error(w, "Error en el logueo en Proyeccion. "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

	}
}
