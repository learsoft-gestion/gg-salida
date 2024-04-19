package src

import (
	"encoding/json"
	"os"
)

type Config struct {
	EntradaDirectorio string `json:"entradaDirectorio"`
	SalidaDirectorio  string `json:"salidaDirectorio"`
}

var Configuracion Config

func CargarConfiguracion() error {
	archivoConfig, err := os.Open("./config/migrador.json")
	if err != nil {
		return err
	}
	defer archivoConfig.Close()

	decodificador := json.NewDecoder(archivoConfig)
	return decodificador.Decode(&Configuracion)
}
