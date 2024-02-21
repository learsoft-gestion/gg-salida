package conexiones

import (
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/lib/pq"
)

func ConectarBase(dbName string, config string, system string) (*sql.DB, error) {

	var connectionString string
	database := fmt.Sprintf("%s_%s", dbName, config)
	credenciales, err := CargarConfiguracion(database)
	if err != nil {
		return nil, err
	}

	if system == "sqlserver" {
		connectionString = fmt.Sprintf("Server=%s;Database=%s;User ID=%s;Password=%s", credenciales.Host, credenciales.Name, credenciales.Username, credenciales.Password)
	} else {
		connectionString = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s TimeZone=UTC+3", credenciales.Host, credenciales.Port, credenciales.Username, credenciales.Password, credenciales.Name)
	}

	db, err := sql.Open(system, connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Conexion exitosa a la base de datos %s de direccion: %s\n", credenciales.Name, credenciales.Host)
	return db, nil
}
