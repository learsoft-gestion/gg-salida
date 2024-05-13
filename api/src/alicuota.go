package src

import (
	"Nueva/modelos"
	"database/sql"
	"fmt"
)

func Alicuota(db *sql.DB, id_convenio int, fechaDesde, fechaHasta string) ([]modelos.Alicuota, error) {

	var alicuotas []modelos.Alicuota
	var valorAlicuota string
	var valoresAlicuotas []modelos.Alicuota

	queryAlicuotas := fmt.Sprintf("select ea.nombre, '$ALICUOTA_'||ea.nombre||'$' from extractor.ext_alicuotas ea where id_convenio = %v", id_convenio)
	filas, err := db.Query(queryAlicuotas)
	if err != nil {
		return nil, err
	}

	for filas.Next() {
		var nombreAli string
		var replaceAli string
		err = filas.Scan(&nombreAli, &replaceAli)
		if err != nil {
			return nil, err
		}
		alicuotas = append(alicuotas, modelos.Alicuota{NombreAli: nombreAli, ReplaceAli: replaceAli})
	}

	for _, alicuota := range alicuotas {
		db.QueryRow("select extractor.obt_alicuota($1,$2,$3,$4)", id_convenio, alicuota.NombreAli, fechaDesde, fechaHasta).Scan(&valorAlicuota)
		valoresAlicuotas = append(valoresAlicuotas, modelos.Alicuota{ValorAli: valorAlicuota, ReplaceAli: alicuota.ReplaceAli})
	}

	return valoresAlicuotas, nil
}
