package src

import (
	"Nueva/modelos"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Formato struct {
	Campo           string
	PosicionInicial int
	PosicionFinal   int
}

func CargarTxt(db *sql.DB, idLogDetalle int, proceso modelos.Proceso, data []modelos.Registro, nombreSalida string) (string, error) {
	// Leer archivo de plantilla
	var plantilla modelos.Plantilla
	path := "./templates/" + proceso.Archivo_modelo
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

	// Abrir archivo para escritura
	archivo, err := os.Create(nombreSalida)
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", err
	}
	defer archivo.Close()

	texto := ""
	blancos := strings.Repeat(" ", 300)
	for _, dato := range data {
		for i, campo := range plantilla.Campos {
			var value string

			if strings.ToLower(campo.Tipo) == "suma" {
				partes := strings.Split(campo.Formato, ",")
				var acumulador float64
				for _, parte := range partes {
					campoSuma := strings.ToUpper(strings.TrimSpace(parte))
					valor := dato.Valores[campoSuma]
					switch v := valor.(type) {
					case []byte:
						// valorStr := string(v)
						valorFloat := valueToFloat(v)
						acumulador += valorFloat
					}
				}
				value = fmt.Sprintf("%.2f", acumulador)
			} else {
				// Validaciones
				if plantilla.Cabecera.Formato == "" {
					return "", fmt.Errorf("debe agregarle un formato a la cabecera del template")
				}
				if campo.Titulo != "" {
					return "", fmt.Errorf("JSON: el campo %s no debe tener Titulo", campo.Nombre)
				}
				if strings.ToLower(campo.Tipo) == "fijo" && campo.Formato == "" {
					return "", fmt.Errorf("JSON: tipo de dato fijo sin formato para el campo %s", campo.Nombre)
				}
				if strings.ToLower(campo.Tipo) == "float" {
					if strings.ToLower(campo.Formato) != "coma" && strings.ToLower(campo.Formato) != "punto" && !strings.Contains(strings.ToLower(campo.Formato), "rellenado") {
						return "", fmt.Errorf("JSON: tipo de dato float con formato erroneo para el campo %s", campo.Nombre)
					}
				}
				if campo.Nombre == "" && strings.ToLower(campo.Tipo) != "fijo" {
					return "", fmt.Errorf("campo sin nombre")
				}
				if campo.Nombre == "" {
					value += campo.Formato
				} else {
					campo.Nombre = strings.ToUpper(campo.Nombre)
					val := dato.Valores[campo.Nombre]
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
						} else if campo.Formato == "DD/MM/YYYY" {
							value = formatearFecha(v, campo.Formato)
						} else if campo.Formato == "DDMMYYYY" {
							value = formatearFecha(v, campo.Formato)
						} else if campo.Tipo == "condicional" {
							condiciones := strings.Split(campo.Formato, "/")
							if string(v) != "0" {
								value = condiciones[0]
							} else {
								value = condiciones[1]
							}
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
							if value == "" {
								value = "12"
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
						} else if strings.ToLower(campo.Tipo) == "condicional" {

							condiciones := strings.Split(campo.Formato, "/")
							if string(v) != "0.00" {
								value = condiciones[0]
							} else {
								value = condiciones[1]
							}
							// if string(v) != "0.00" {
							// 	value = campo.Option1
							// } else {
							// 	value = campo.Option2
							// }
						} else if strings.ToLower(campo.Tipo) == "float sin caracteres" {
							value += strings.Replace(string(v), ".", "", -1) // Numero sin puntos
						} else {
							value += string(v)
						}
					case int64:
						condiciones := strings.Split(campo.Formato, "/")
						if int64(v) == 0 {
							value = condiciones[0]
						} else {
							value = condiciones[1]
						}
					default:
						// fmt.Printf("Campo: %s Valor: %v Tipo: %s\n", campo.Nombre, v, reflect.TypeOf(v))
						value = "#"
					}
				}
			}

			// fmt.Printf("Campo: %s Valor: %s\n", campo.Nombre, value)

			if plantilla.Cabecera.Formato == "fijo" {

				longitud_campo := campo.Fin - campo.Inicio + 1

				if len(value) < longitud_campo {
					diferencia := longitud_campo - len(value)
					if strings.Contains(strings.ToLower(campo.Formato), "rellenado") {
						partesFormato := strings.Split(campo.Formato, "/")
						if len(partesFormato) < 2 {
							fmt.Println("Formato de relleno incorrecto")
							return "", err
						} else {
							value = strings.Repeat(partesFormato[1], diferencia-1) + value
						}
					} else if value == "#" {
						value += strings.Repeat("#", diferencia)
					} else {
						value += strings.Repeat(" ", diferencia)
					}

				}

				fmt.Println(campo.Formato)
				// Iterar sobre blancos para agregar letra por letra
				arreglo := []rune(blancos)
				palabra := []rune(value)
				for i := 0; i < len(palabra); i++ {
					arreglo[campo.Inicio+i] = palabra[i]
				}
				blancos = string(arreglo)
			} else if plantilla.Cabecera.Formato == "variable" {
				if i == 0 {
					blancos = strings.TrimSpace(value)
				} else {
					blancos += plantilla.Cabecera.Separador + value
				}
			}

		}
		texto += strings.TrimSpace(blancos) + "\n"

	}

	// fmt.Printf("Cadena: %s\n", texto)
	_, err = fmt.Fprintf(archivo, "%s", strings.TrimSpace(texto))
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", err
	}

	return nombreSalida, nil
}
