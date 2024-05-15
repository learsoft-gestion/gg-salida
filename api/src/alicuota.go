package src

import (
	"Nueva/modelos"
	"database/sql"
	"fmt"
)

func Alicuota(db *sql.DB, id_convenio int, fechaDesde, fechaHasta string) ([]modelos.AlicuotaBack, error) {

	var alicuotas []modelos.AlicuotaBack
	var valorAlicuota string
	var valoresAlicuotas []modelos.AlicuotaBack

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
		alicuotas = append(alicuotas, modelos.AlicuotaBack{NombreAli: nombreAli, ReplaceAli: replaceAli})
	}

	for _, alicuota := range alicuotas {
		db.QueryRow("select extractor.obt_alicuota($1,$2,$3,$4)", id_convenio, alicuota.NombreAli, fechaDesde, fechaHasta).Scan(&valorAlicuota)
		valoresAlicuotas = append(valoresAlicuotas, modelos.AlicuotaBack{ValorAli: valorAlicuota, ReplaceAli: alicuota.ReplaceAli})
	}

	return valoresAlicuotas, nil
}
