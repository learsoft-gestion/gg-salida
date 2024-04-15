package src

import (
	"Nueva/modelos"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func CargarXml(db *sql.DB, idLogDetalle int, proceso modelos.Proceso, registros []modelos.Registro, nombreSalida string) (string, error) {
	// Leer archivo de plantilla
	var plantilla modelos.Plantilla
	path := "../templates/" + proceso.Archivo_modelo
	fmt.Println("Path: ", path)
	file, err := os.Open(path)
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&plantilla)
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", err
	}

	// Crear un archivo para escribir el XML
	archivo, err := os.Create(nombreSalida)
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", err
	}
	defer archivo.Close()
	texto := fmt.Sprintf("%s\n", plantilla.Cabecera.XmlTag)
	texto += fmt.Sprintf("<%s>\n", plantilla.Cabecera.Tag)
	elemento := ""
	for _, registro := range registros {
		texto += fmt.Sprintf("\t<%s>\n", plantilla.Cabecera.Children)
		for _, campo := range plantilla.Campos {
			var value string

			// Validaciones
			if strings.ToLower(campo.Tipo) == "fijo" && campo.Formato == "" {
				return "", fmt.Errorf("JSON: tipo de dato fijo sin formato para el campo %s", campo.Nombre)
			}
			if strings.ToLower(campo.Tipo) == "float" {
				if strings.ToLower(campo.Formato) != "coma" && strings.ToLower(campo.Formato) != "punto" {
					return "", fmt.Errorf("JSON: tipo de dato float con formato erroneo para el campo %s", campo.Nombre)
				}
			}
			if campo.Nombre == "" && strings.ToLower(campo.Tipo) != "fijo" && strings.ToLower(campo.Tipo) != "suma" {
				return "", fmt.Errorf("campo sin nombre")
			}
			if campo.Nombre == "" && strings.ToLower(campo.Tipo) != "suma" {
				value += campo.Formato
			} else if strings.ToLower(campo.Tipo) == "suma" {
				partes := strings.Split(campo.Formato, ",")
				var acumulador float64
				for _, parte := range partes {
					campoSuma := strings.ToUpper(strings.TrimSpace(parte))
					valor := registro.Valores[campoSuma]
					switch v := valor.(type) {
					case []byte:
						// valorStr := string(v)
						valorFloat := valueToFloat(v)
						acumulador += valorFloat
					}
				}
				value = fmt.Sprintf("%.2f", acumulador)
			} else {
				campo.Nombre = strings.ToUpper(campo.Nombre)
				val := registro.Valores[campo.Nombre]
				switch v := val.(type) {
				case int:
					value += fmt.Sprintf("%d", v)
				case float64:
					value += fmt.Sprintf("%.2f", v)
				case string:
					if strings.ToLower(campo.Formato) == "cuil sin guion" {
						value += strings.ReplaceAll(v, "-", "")
					} else if campo.Formato == "MM/YYYY" {
						value = formatearFecha(v, campo.Formato)
					} else if strings.ToLower(campo.Tipo) == "lookup" {
						// El dato lo saco del .json
						for _, variable := range plantilla.Variables {
							if strings.ToUpper(variable.Nombre) == campo.Nombre {
								for _, element := range variable.Datos {
									if element.Nombre == v {
										value += fmt.Sprintf("%v", element.Id)
									}
								}
							}
						}
					} else {
						value += v
					}

				case []int:
					value += fmt.Sprintf("%v", v)
				case []byte:
					if strings.ToLower(campo.Tipo) == "float" && strings.ToLower(campo.Formato) != "coma" {
						value += string(v)
					} else if strings.ToLower(campo.Tipo) == "float" && strings.ToLower(campo.Formato) == "coma" {
						value += strings.Replace(string(v), ".", ",", -1)
					}
				default:
					value = fmt.Sprintf("##%s##", campo.Nombre)
				}
			}

			elemento = fmt.Sprintf("\t\t<%s>%s</%s>\n", campo.Titulo, strings.TrimSpace(value), campo.Titulo)
			texto += elemento

		}
		texto += fmt.Sprintf("\t</%s>\n", plantilla.Cabecera.Children)
	}
	texto += fmt.Sprintf("</%s>\n", plantilla.Cabecera.Tag)

	num, err := fmt.Fprintf(archivo, "%s", texto)
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", err
	}
	fmt.Printf("Archivo XML generado con %v bytes\n", num)

	return nombreSalida, nil
}
