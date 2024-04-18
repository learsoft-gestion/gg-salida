package src

import (
	"Nueva/modelos"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func CargarExcel(db *sql.DB, idLogDetalle int, proceso modelos.Proceso, data []modelos.Registro, nombreSalida string, path string) (string, error) {
	// Leer archivo de plantilla
	var plantilla modelos.Plantilla
	// path := "../templates/" + proceso.Archivo_control
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

	fileNuevo := excelize.NewFile()

	sheetName := "Hoja1"
	fileNuevo.SetSheetName("Sheet1", sheetName)

	fileNuevo.SetColWidth(sheetName, "A", "CA", 15)

	styleMoneda, _ := fileNuevo.NewStyle(&excelize.Style{NumFmt: 177})
	styleNumero, _ := fileNuevo.NewStyle(&excelize.Style{NumFmt: 1})
	styleNumeroDecimal, _ := fileNuevo.NewStyle(&excelize.Style{NumFmt: 2})
	// styleDefault, _ := fileNuevo.NewStyle(&excelize.Style{Alignment: al})
	styleEncabezadoNomina, _ := fileNuevo.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:  10,
			Color: "#FF0000",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#a7a7a7"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
			WrapText:   true,
		},
	})
	// styleColumnaControl, _ := fileNuevo.NewStyle(&excelize.Style{
	// 	Font: &excelize.Font{
	// 		Color: "#FFFFFF",
	// 	},
	// 	Fill: excelize.Fill{
	// 		Type:    "pattern",
	// 		Color:   []string{"#000000"},
	// 		Pattern: 1,
	// 	},
	// 	Alignment: &excelize.Alignment{
	// 		Horizontal: "center",
	// 		Vertical:   "center",
	// 		WrapText:   true,
	// 	},
	// })

	if strings.ToLower(plantilla.Cabecera.Sentido_encabezado) == "vertical" {
		// Escribir verticalmente encabezados en el Excel
		for _, campo := range plantilla.Campos {
			cell := "B" + campo.Columna
			fileNuevo.SetCellValue(sheetName, cell, campo.Titulo)
		}
	} else {
		// Escribir horizontalmente encabezados en el Excel
		for _, campo := range plantilla.Campos {
			cell := campo.Columna + "1"
			fileNuevo.SetCellValue(sheetName, cell, campo.Titulo)
		}
		if strings.ToLower(plantilla.Cabecera.Estilo) == "nomina" {
			fileNuevo.SetRowStyle(sheetName, 1, 1, styleEncabezadoNomina)
			fileNuevo.SetRowHeight(sheetName, 1, 50)
		}
	}

	// Escribir datos en el archivo Excel
	for i, registro := range data {
		for _, campo := range plantilla.Campos {
			var value string

			// Validaciones

			if campo.Nombre == "" {
				value += campo.Formato
			} else {
				campo.Nombre = strings.ToUpper(campo.Nombre)
				val := registro.Valores[campo.Nombre]

				if campo.Columna == "" {
					return "", fmt.Errorf("JSON: el campo %s no tiene columna", campo.Titulo)
				}
				if campo.Inicio != 0 || campo.Fin != 0 {
					return "", fmt.Errorf("JSON: el campo %s no debe tener 'inicio' ni 'fin'", campo.Titulo)
				}
				if campo.Tipo == "fecha" {
					if campo.Formato != "DD/MM/YYYY" && campo.Formato != "DD-MM-YYYY" && campo.Formato != "YYYYMMDD" {
						return "", fmt.Errorf("JSON: formato desconocido para %s", campo.Titulo)
					}
				}
				if campo.Formato == "DD/MM/YYYY" || campo.Formato == "DD-MM-YYYY" {
					if campo.Tipo != "fecha" {
						return "", fmt.Errorf("JSON: el campo %s debe ser de tipo fecha", campo.Titulo)
					}
				}

				switch v := val.(type) {
				case int:
					value += strings.TrimSpace(fmt.Sprintf("%d", v))
				case float64:
					value += strings.TrimSpace(fmt.Sprintf("%.2f", v))
				case string:
					if strings.ToLower(campo.Formato) == "cuil sin guion" {
						value = strings.ReplaceAll(v, "-", "")
					} else if strings.ToLower(campo.Tipo) == "lookup" {
						// El dato lo saco del .json
						for _, variable := range plantilla.Variables {
							if strings.ToUpper(variable.Nombre) == campo.Nombre {
								for _, element := range variable.Datos {
									if element.Id == v[4:] {
										value += fmt.Sprintf("%s-%v", element.Nombre, v[2:4])
									}
								}
							}
						}
					} else if campo.Formato == "DD/MM/YYYY" {
						numRegex := regexp.MustCompile(`^\s{8}$`)
						val := strings.TrimSpace(v)
						if len(val) == 8 {
							if numRegex.MatchString(v) {
								value += strings.TrimSpace(v)
							} else {
								value = formatearFecha(v, campo.Formato)
							}
						}
					} else if strings.ToLower(campo.Tipo) == "condicional" {
						numRegex := regexp.MustCompile(`^\s{8}$`)
						condiciones := strings.Split(campo.Formato, "/")
						if numRegex.MatchString(v) {
							value = condiciones[1]
						} else {
							value = condiciones[0]
						}
					} else {
						value += strings.TrimSpace(v)
					}
				case []int:
					value += strings.TrimSpace(fmt.Sprintf("%v", v))
				case []byte:
					if strings.ToLower(campo.Tipo) == "float" && strings.ToLower(campo.Formato) == "coma" {
						value += strings.Replace(string(v), ".", ",", -1)
					} else if strings.ToLower(campo.Tipo) == "float" && strings.ToLower(campo.Formato) == "millares con coma" {
						value += formatearFloat(string(v))
					} else {
						value += strings.TrimSpace(string(v))
					}
				case nil:
					value = ""
				default:
					value = strings.TrimSpace(fmt.Sprintf("%v", v))
				}
				if campo.Tipo == "fijo" {
					value = campo.Formato
				}
			}
			value = strings.TrimSpace(value)

			if strings.ToLower(plantilla.Cabecera.Sentido_encabezado) == "vertical" {
				colLetter := ObtenerLetra(i + 3)
				cell := colLetter + campo.Columna
				if strings.ToLower(campo.Tipo) == "moneda" {
					valor, _ := strconv.ParseFloat(value, 64)
					fileNuevo.SetCellValue(sheetName, cell, valor)
					fileNuevo.SetCellStyle(sheetName, cell, cell, styleMoneda)
				} else if strings.ToLower(campo.Tipo) == "int" {
					valor, _ := strconv.Atoi(value)
					fileNuevo.SetCellValue(sheetName, cell, valor)
					fileNuevo.SetCellStyle(sheetName, colLetter, campo.Columna, styleNumero)
				} else if strings.ToLower(campo.Tipo) == "numero decimal" {
					fileNuevo.SetCellStyle(sheetName, colLetter, campo.Columna, styleNumeroDecimal)
					fileNuevo.SetCellValue(sheetName, cell, value)
				} else {
					fileNuevo.SetCellValue(sheetName, cell, value)
				}
			} else {
				cell := campo.Columna + fmt.Sprintf("%v", i+2)
				if strings.ToLower(campo.Tipo) == "moneda" {
					valor, _ := strconv.ParseFloat(value, 64)
					fileNuevo.SetCellValue(sheetName, cell, valor)
					fileNuevo.SetCellStyle(sheetName, cell, cell, styleMoneda)
				} else if strings.ToLower(campo.Tipo) == "int" {
					if len(value) > 0 {
						valor, err := strconv.ParseFloat(value, 64)
						if err != nil {
							return "", err
						}
						fileNuevo.SetCellValue(sheetName, cell, valor)
					} else {
						fileNuevo.SetCellValue(sheetName, cell, value)
					}
					fileNuevo.SetCellStyle(sheetName, cell, cell, styleNumero)
				} else {
					fileNuevo.SetCellValue(sheetName, cell, value)
					if campo.Ancho > 0 {
						fileNuevo.SetColWidth(sheetName, campo.Columna, campo.Columna, float64(campo.Ancho))
					}
				}
			}
		}
	}

	// Guardar archivo
	if err = fileNuevo.SaveAs(nombreSalida); err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", err
	}

	return nombreSalida, nil
}

func ObtenerLetra(numero int) string {
	// Asumiendo que estamos trabajando con el alfabeto inglés (a-z)
	alfabeto := "abcdefghijklmnopqrstuvwxyz"

	longitudAlfabeto := len(alfabeto)
	numRepeticiones := numero / longitudAlfabeto
	indice := numero % longitudAlfabeto

	// Si el índice es 0, corresponde a la última letra del alfabeto
	if indice == 0 {
		indice = longitudAlfabeto
		// Si el índice es 0, debemos reducir el número de repeticiones
		numRepeticiones--
	}

	primeraLetra := string(alfabeto[indice-1])

	// Si hay repeticiones, obtener la segunda letra
	var segundaLetra string
	if numRepeticiones > 0 {
		segundaLetra = string(alfabeto[numRepeticiones-1])
	}

	// Combinar las letras
	letras := strings.Repeat(segundaLetra, numRepeticiones) + primeraLetra

	return letras
}

// 1000.1 ---> 1000,1
func formatearFloat(s string) string {
	partes := strings.Split(s, ".")
	str := formatearMillares(partes[0])
	return str + "," + partes[1]
}

// 1112223 ----> 111.222,3
func formatearMillares(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	return formatearMillares(s[:n-3]) + "." + s[n-3:]
}

// YYYYMMDD -----> DDMMYYYY -----> DD/MM/YYYY
func formatearFecha(s string, formato string) string {
	var strFinal string
	var partes []string

	// Verificar formato y transformar
	if formato == "DD/MM/YYYY" {
		partes = append(partes, s[:4], s[4:6], s[6:])
		strFinal = partes[2] + "/" + partes[1] + "/" + partes[0]
	} else if formato == "MM/YYYY" {
		partes = append(partes, s[:4], s[4:])
		strFinal = partes[1] + "/" + partes[0]
	}
	return strFinal
}
