package handlers

import (
	"Nueva/config"
	"Nueva/modelos"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type ExtArchivo struct {
	IdNumero       int
	Procesado      bool
	FechaProcesado sql.NullTime
	ArchivoEntrada string
	ArchivoFinal   string
	Empresa        string
	Periodo        string
	Convenio       string
	Estado         sql.NullString
	Descripcion    sql.NullString
}

func ProcesarArchivo(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := config.CargarConfiguracion()
		if err != nil {
			log.Fatal(err)
		}

		// Traer body
		var datos modelos.Option
		err = json.NewDecoder(r.Body).Decode(&datos)
		if err != nil {
			http.Error(w, "Error decodificando JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		query := "SELECT id_numero, procesado, fecha_procesado, archivo_entrada, archivo_final, empresa, periodo, convenio, estado, descripcion FROM extractor.ext_archivos WHERE id_numero = $1"
		row := db.QueryRow(query, datos.Id)

		var extArchivo ExtArchivo
		err = row.Scan(&extArchivo.IdNumero, &extArchivo.Procesado, &extArchivo.FechaProcesado, &extArchivo.ArchivoEntrada,
			&extArchivo.ArchivoFinal, &extArchivo.Empresa, &extArchivo.Periodo, &extArchivo.Convenio,
			&extArchivo.Estado, &extArchivo.Descripcion)
		if err != nil {
			log.Fatal(err)
		}

		// Iniciar transacción
		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}

		// Procesar archivo dentro de la transacción
		msg, err := procesarArchivo(extArchivo, tx)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}

		// Confirmar la transacción
		err = tx.Commit()
		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		respuesta := modelos.Respuesta{
			Mensaje: msg,
		}
		jsonResp, _ := json.Marshal(respuesta)
		w.Write(jsonResp)
	}
}

func procesarArchivo(registro ExtArchivo, tx *sql.Tx) (string, error) {
	rutaEntrada, err := buscarArchivo(registro.ArchivoEntrada)
	if err != nil {
		return "", err
	}
	if rutaEntrada == "" {
		log.Printf("Archivo %s no encontrado\n", registro.ArchivoEntrada)
		err = actualizarRegistroNoEncontrado(tx, registro.IdNumero)
		if err != nil {
			return "", err
		}
		return "Archivo no encontrado", nil
	}
	rutaSalida, err := buscarOCrearDirectorio(registro)
	rutaSalida = filepath.Join(rutaSalida, registro.ArchivoFinal)
	moverArchivo(rutaEntrada, rutaSalida)
	if err != nil {
		return "", err
	}
	err = actualizarRegistroEncontrado(tx, registro.IdNumero)
	if err != nil {
		return "", err
	}
	return "Archivo procesado exitosamente", nil
}

func buscarArchivo(nombreArchivo string) (string, error) {
	var ruta string

	err := filepath.Walk(config.Configuracion.EntradaDirectorio, func(rutaActual string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && info.Name() == nombreArchivo {
			ruta = rutaActual
		}

		return nil
	})

	return ruta, err
}

func actualizarRegistroNoEncontrado(db *sql.Tx, id int) error {
	_, err := db.Exec("UPDATE extractor.ext_archivos SET procesado = $1, fecha_procesado = $2, estado = $3, descripcion = $4 WHERE id_numero = $5",
		true, time.Now(), "E", "Archivo no encontrado", id)
	return err
}

func actualizarRegistroEncontrado(db *sql.Tx, id int) error {
	_, err := db.Exec("UPDATE extractor.ext_archivos SET procesado = $1, fecha_procesado = $2, estado = $3, descripcion = $4 WHERE id_numero = $5",
		true, time.Now(), "F", "Procesado correctamente", id)
	return err
}

func buscarOCrearDirectorio(registro ExtArchivo) (string, error) {
	salida := config.Configuracion.SalidaDirectorio
	empresa := registro.Empresa
	periodo := registro.Periodo
	convenio := registro.Convenio

	rutaDirectorio := filepath.Join(salida, empresa, periodo, convenio)

	_, err := os.Stat(rutaDirectorio)
	if os.IsNotExist(err) {
		err := os.MkdirAll(rutaDirectorio, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("no se pudo crear el directorio: %v", err)
		}
	}

	return rutaDirectorio, nil
}

func moverArchivo(origen, destino string) error {
	// Abrir el archivo de origen
	archivoOrigen, err := os.Open(origen)
	if err != nil {
		log.Println("No se pudo abrir el archivo", err)
		return err
	}

	// Crear el archivo de destino
	archivoDestino, err := os.Create(destino)
	if err != nil {
		log.Println("No se pudo crear archivo en destino", err)
		return err
	}
	defer archivoDestino.Close()

	// Copiar el contenido del archivo de origen al archivo de destino
	_, err = io.Copy(archivoDestino, archivoOrigen)
	if err != nil {
		log.Println("No se pudo copiar archivo", err)
		return err
	}

	archivoOrigen.Close()

	// Eliminar el archivo de origen
	err = os.Remove(origen)
	if err != nil {
		log.Println("No se pudo eliminar archivo en origen", err)
		return err
	}

	log.Printf("Archivo movido de %s a %s\n", origen, destino)
	return nil
}
