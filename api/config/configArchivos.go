package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	EntradaDirectorio string `json:"entrada"`
	SalidaDirectorio  string `json:"salida"`
}

var Configuracion Config

func CargarConfiguracion() error {
	archivoConfig, err := os.Open("./config/configArch.json")
	if err != nil {
		return err
	}
	defer archivoConfig.Close()

	decodificador := json.NewDecoder(archivoConfig)
	return decodificador.Decode(&Configuracion)
}
