package src

import (
	"Nueva/modelos"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
)

func CargarExcel(db *sql.DB, idLogDetalle int, proceso modelos.Proceso, data []modelos.Registro, nombreSalida string) (string, error) {
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

	// Construir path
	// path := fmt.Sprintf("../templates/%s", proceso.Archivo_modelo)

	// Abrir archivo Excel
	// f, err := excelize.OpenFile(path)
	// if err != nil {
	// 	ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
	// 	return "", err
	// }

	// Nombre de la hoja
	// sheetName := f.GetSheetName(0)

	// Obtener las celdas de la primera fila
	// filas, err := f.GetRows(sheetName)
	// if err != nil {
	// 	ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
	// 	return "", err
	// }
	fileNuevo := excelize.NewFile()

	// var estilos []int
	// for i := range filas[0] {
	// 	colLetter, _ := excelize.ColumnNumberToName(i + 1)
	// 	cell := colLetter + "1"
	// 	styles, err := f.GetCellStyle(sheetName, cell)
	// 	if err != nil {
	// 		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
	// 		return "", err
	// 	}
	// 	estilo, err := f.GetStyle(styles)
	// 	if err != nil {
	// 		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
	// 		return "", err
	// 	}
	// 	fileNuevo.NewStyle(estilo)
	// 	fileNuevo.SetCellStyle(sheetName, colLetter, "1", styles)
	// 	// estilos = append(estilos, styles)
	// }

	// Extraer valores de la primera fila
	// var primeraFila []string
	// primeraFila = append(primeraFila, filas[0]...)
	sheetName := "Hoja1"
	fileNuevo.SetSheetName("Sheet1", sheetName)

	// Escribir encabezados en el Excel
	for _, campo := range plantilla.Campos {
		cell := campo.Columna + "1"
		// fmt.Println("Nombre del campo: ", campo.Titulo)
		fileNuevo.SetCellValue(sheetName, cell, campo.Titulo)
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

				// Validaciones
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
				if campo.Tipo != "string" && campo.Tipo != "float" && campo.Tipo != "fecha" && campo.Tipo != "fijo" {
					return "", fmt.Errorf("JSON: tipo desconocido para %s", campo.Titulo)
				}

				switch v := val.(type) {
				case int:
					value += fmt.Sprintf("%d", v)
				case float64:
					value += fmt.Sprintf("%.2f", v)
				case string:
					if campo.Nombre == "CAT_REDUCIDO" {
						fmt.Println("CAT_REDUCIDO: ", v)
					}
					if strings.ToLower(campo.Formato) == "cuil sin guion" {
						value += strings.ReplaceAll(v, "-", "")
					} else if campo.Formato == "DD/MM/YYYY" {
						value += formatearFecha(v, campo.Formato)
					} else {
						value += v
					}
				case []int:
					value += fmt.Sprintf("%v", v)
				case []byte:
					if strings.ToLower(campo.Tipo) == "float" && strings.ToLower(campo.Formato) == "coma" {
						value += strings.Replace(string(v), ".", ",", -1)
					} else if strings.ToLower(campo.Tipo) == "float" && strings.ToLower(campo.Formato) == "millares con coma" {
						value += formatearFloat(string(v))
					} else {
						value += string(v)
					}
				case nil:
					value = fmt.Sprintf("##%s##", campo.Nombre)
				default:
					value = fmt.Sprintf("%v", v)
				}
				if campo.Tipo == "fijo" {
					value = campo.Formato
				}
			}
			// colLetter, _ := excelize.ColumnNumberToName(colIndex + 1)
			// cell := colLetter + "1"
			cell := campo.Columna + fmt.Sprintf("%v", i+2)
			fileNuevo.SetCellValue(sheetName, cell, value)
		}
	}

	// err := fileNuevo.SetRowStyle(sheetName, 1, 10, estilos[1])
	// if err != nil {
	// 	return err
	// }

	// Escribir datos en el archivo Excel
	// for rowIndex, registro := range data { // row ----> 1 Gabi CABA 27
	// 	for colIndex, columna := range primeraFila {
	// 		// columnaMin := strings.ToLower(columna)
	// 		colLetter, _ := excelize.ColumnNumberToName(colIndex + 1)
	// 		cell := colLetter + strconv.Itoa(rowIndex+2)
	// 		valor, ok := registro.Valores[strings.ToUpper(columna)]
	// 		if ok {
	// 			f.SetCellValue(sheetName, cell, valor)
	// 		} else {
	// 			f.SetCellValue(sheetName, cell, columna)
	// 		}
	// 		// fmt.Printf("sheetName: %s cell: %s value: %s\n", sheetName, cell, valor)
	// 	}
	// }

	// Guardar archivo
	if err = fileNuevo.SaveAs(nombreSalida); err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", err
	}

	return nombreSalida, nil
}

func formatearFloat(s string) string {
	partes := strings.Split(s, ".")
	str := formatearMillares(partes[0])
	return str + "," + partes[1]
}

// 1112223
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
	if formato == "DD/MM/YYYY" {
		partes = append(partes, s[:4], s[4:6], s[6:])
		strFinal = partes[2] + "/" + partes[1] + "/" + partes[0]
	} else if formato == "MM/YYYY" {
		partes = append(partes, s[:4], s[4:])
		strFinal = partes[0] + "/" + partes[1]
	}
	return strFinal
}
