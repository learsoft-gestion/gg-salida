package conexiones

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Config es la estructura que representa el JSON de configuración
type Config struct {
	DatabaseConfigs map[string]DbConfig `json:"databases"`
}

// DbConfig es la estructura que representa la configuración de una base de datos
type DbConfig struct {
	Conectar string `json:"conectar"`
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// var Dbc DbConfig = DbConfig{}
var Dbc Config = Config{}
var dbpass string
var content []byte
var err error

func CargarConfiguracion(database string) (DbConfig, error) {
	var dbName string
	// var config string
	// var dbConfigStruct DbConfig

	configFile := "./config/config.json"

	if strings.Contains(database, "_") {
		partes := strings.Split(database, "_")
		dbName = partes[0]
		// config = partes[1]
		// fmt.Printf("dbName = %s config = %s \n", dbName, config) // Para saber que me traigo
	}

	//...................................
	//Reading into struct type from a JSON file
	//...................................

	content, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Println("Error al leer configFile en la ruta ./config/config.json")
		return DbConfig{}, err
	}

	err = json.Unmarshal(content, &Dbc)
	if err != nil {
		return DbConfig{}, err
	}

	dbConfigStruct, ok := Dbc.DatabaseConfigs[database]
	if !ok {
		return DbConfig{}, fmt.Errorf("no se encontro la configuracion para la base de datos %s", dbName)
	}

	if strings.HasPrefix(dbConfigStruct.Password, "##NLP##") {
		dbpass = dbConfigStruct.Password[7:]
		config := Dbc.DatabaseConfigs[database]
		config.Password = Encriptar(dbpass)
		Dbc.DatabaseConfigs[database] = config
		//...................................
		//Writing struct type to a JSON file
		//...................................
		content, err = json.MarshalIndent(Dbc, "", "  ")
		if err != nil {
			return DbConfig{}, err
		}
		err = os.WriteFile(configFile, content, 0644)
		if err != nil {
			return DbConfig{}, err
		}
	} else {
		dbpass = Desencriptar(dbConfigStruct.Password)
		dbConfigStruct.Password = dbpass
	}
	return dbConfigStruct, nil
}
