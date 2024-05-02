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

func CargarExcel(db *sql.DB, idLogDetalle int, proceso modelos.Proceso, data []modelos.Registro, nombreSalida string, path string, tipo_ejecucion string) (string, error) {

	var plantilla modelos.Plantilla

	if tipo_ejecucion != "control" {

		// fmt.Println("Path: ", path)
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
	}

	fileNuevo := excelize.NewFile()

	sheetName := "Hoja1"
	fileNuevo.SetSheetName("Sheet1", sheetName)

	fileNuevo.SetColWidth(sheetName, "A", "CA", 15)
	// fuente_valores := &excelize.Font{
	// 	Size: 9,
	// }
	// borde := []excelize.Border{
	// 	{
	// 		Type:  "top",
	// 		Color: "000000",
	// 		Style: 1,
	// 	},
	// 	{
	// 		Type:  "bottom",
	// 		Color: "000000",
	// 		Style: 1,
	// 	},
	// 	{
	// 		Type:  "left",
	// 		Color: "000000",
	// 		Style: 1,
	// 	},
	// 	{
	// 		Type:  "right",
	// 		Color: "000000",
	// 		Style: 1,
	// 	},
	// },

	styleMoneda, _ := fileNuevo.NewStyle(&excelize.Style{
		NumFmt: 177,
		Font: &excelize.Font{
			Size: 9,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   false,
		},
		Border: []excelize.Border{
			{
				Type:  "top",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "bottom",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "left",
				Color: "000000",
				Style: 1,
			},
			{
				Type:  "right",
				Color: "000000",
				Style: 1,
			},
		},
	})
	styleNumero, _ := fileNuevo.NewStyle(&excelize.Style{NumFmt: 1, Border: []excelize.Border{{Type: "top", Color: "000000", Style: 1}, {Type: "bottom", Color: "000000", Style: 1}, {Type: "left", Color: "000000", Style: 1}, {Type: "right", Color: "000000", Style: 1}}})
	styleNumeroDecimal, _ := fileNuevo.NewStyle(&excelize.Style{NumFmt: 2, Border: []excelize.Border{{Type: "top", Color: "000000", Style: 1}, {Type: "bottom", Color: "000000", Style: 1}, {Type: "left", Color: "000000", Style: 1}, {Type: "right", Color: "000000", Style: 1}}})
	styleEncabezadoNomina, _ := fileNuevo.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 10, Color: "#FF0000"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#a7a7a7"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", WrapText: true},
		Border:    []excelize.Border{{Type: "top", Color: "000000", Style: 1}, {Type: "bottom", Color: "000000", Style: 1}, {Type: "left", Color: "000000", Style: 1}, {Type: "right", Color: "000000", Style: 1}},
	})
	styleEncabezadoControl, _ := fileNuevo.NewStyle(&excelize.Style{Font: &excelize.Font{Size: 10}, Border: []excelize.Border{{Type: "top", Color: "000000", Style: 1}, {Type: "bottom", Color: "000000", Style: 1}, {Type: "left", Color: "000000", Style: 1}, {Type: "right", Color: "000000", Style: 1}}, Fill: excelize.Fill{Type: "pattern", Color: []string{"#a7a7a7"}, Pattern: 1}, Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", WrapText: true}})
	styleColumnaControl, _ := fileNuevo.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Color: "#FFD3A7"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#000000"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
	})
	styleTotalesControl, _ := fileNuevo.NewStyle(&excelize.Style{
		NumFmt:    177,
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#a7a7a7"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", ReadingOrder: 0, Indent: 0, RelativeIndent: 0, ShrinkToFit: false, TextRotation: 0, WrapText: false},
		Border:    []excelize.Border{{Type: "top", Color: "000000", Style: 1}, {Type: "bottom", Color: "000000", Style: 1}, {Type: "left", Color: "000000", Style: 1}, {Type: "right", Color: "000000", Style: 1}},
	})
	styleAligned, _ := fileNuevo.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", ReadingOrder: 0, Indent: 0, RelativeIndent: 0, ShrinkToFit: false, TextRotation: 0, WrapText: false},
		Border:    []excelize.Border{{Type: "top", Color: "000000", Style: 1}, {Type: "bottom", Color: "000000", Style: 1}, {Type: "left", Color: "000000", Style: 1}, {Type: "right", Color: "000000", Style: 1}},
	})
	styleDefault, _ := fileNuevo.NewStyle(&excelize.Style{
		Border: []excelize.Border{{Type: "top", Color: "000000", Style: 1}, {Type: "bottom", Color: "000000", Style: 1}, {Type: "left", Color: "000000", Style: 1}, {Type: "right", Color: "000000", Style: 1}},
	})
	styleVertical, _ := fileNuevo.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", TextRotation: 90, WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#C6E0B4"}, Pattern: 1},
		Border:    []excelize.Border{{Type: "top", Color: "000000", Style: 1}, {Type: "bottom", Color: "000000", Style: 1}, {Type: "left", Color: "000000", Style: 1}, {Type: "right", Color: "000000", Style: 1}},
	})

	if tipo_ejecucion == "control" {

		for i, registro := range data {
			colLetter := ObtenerLetra(i + 3) // Almacena la columna para este registro
			colTotalNum := 0                 // Almacena la fila del total de liquidacion
			colTotalNum2 := 0                // Almacena la fial del total de KTNA

			for j, campo := range registro.Columnas {

				// Escribir claves
				cellKey := "B" + strconv.Itoa(j+2)
				fileNuevo.SetCellValue(sheetName, cellKey, campo)
				fileNuevo.SetCellStyle(sheetName, cellKey, cellKey, styleEncabezadoControl) // Fondo gris para el campo
				fileNuevo.SetRowHeight(sheetName, j+2, 25)                                  // Amplia alto de la fila

				// Escribir valores
				value := registro.Valores[strings.ToUpper(campo)]
				cellValue := colLetter + strconv.Itoa(j+2)

				if strings.ToUpper(campo) == "PERIODOLIQ" {
					value = formatearPeriodoLiq(value.(string))
					fileNuevo.SetCellStyle(sheetName, cellValue, cellValue, styleColumnaControl)
				}

				switch v := value.(type) {
				case []uint8:
					valueStr := string(v)
					valueFloat, err := strconv.ParseFloat(valueStr, 64)
					if err != nil {
						fmt.Println(err.Error())
					}
					fileNuevo.SetCellValue(sheetName, cellValue, valueFloat)
					fileNuevo.SetCellStyle(sheetName, cellValue, cellValue, styleMoneda)
				case string:
					fileNuevo.SetCellValue(sheetName, cellValue, v)
					fileNuevo.SetCellStyle(sheetName, cellValue, cellValue, styleAligned)
				case int64:
					fileNuevo.SetCellValue(sheetName, cellValue, v)
					fileNuevo.SetCellStyle(sheetName, cellValue, cellValue, styleAligned)
				default:
					fileNuevo.SetCellValue(sheetName, cellValue, "defValue")
					fmt.Printf("Tipo de dato en %s: %T\n", campo, value)

				}

				if strings.Contains(campo, "KTNA") && strings.Contains(campo, "TOTAL") {

					if !strings.Contains(registro.Columnas[j+1], "KTNA") {
						for k := colTotalNum + 1; k <= j+2; k++ {
							cell := fmt.Sprintf("A%d", k)
							fileNuevo.SetCellValue(sheetName, cell, "TOTAL DESCONTADO KTNA")
							fileNuevo.SetCellStyle(sheetName, cell, cell, styleVertical)
							// colTotal = fmt.Sprintf("%d", j+2)
						}
						colTotalNum2 = j + 3 // es j+2 + 1 ya que inserto una fila antes de los campos de ktna

						// fmt.Println("Proxima iteracion ya no es ktna: ", registro.Columnas[j+1])
						// fmt.Printf("Merge de %s a %s para ktna.\n", "A"+strconv.Itoa(colTotalNum), "A"+strconv.Itoa(j+2))
						fileNuevo.MergeCell(sheetName, "A"+strconv.Itoa(colTotalNum+1), "A"+strconv.Itoa(j+2))
					}

					// fmt.Printf("Celda de KTNA TOTAL: %s\n", "B"+strconv.Itoa(j+2))
					fileNuevo.SetCellStyle(sheetName, "B"+strconv.Itoa(j+2), "B"+strconv.Itoa(j+2), styleTotalesControl)             // Negrita y fondo gris para el total
					fileNuevo.SetCellStyle(sheetName, colLetter+strconv.Itoa(j+2), colLetter+strconv.Itoa(j+2), styleTotalesControl) // Negrita y fondo gris para el total

				} else if strings.Contains(campo, "TOTAL") && !strings.Contains(campo, "KTNA") {

					if strings.Contains(registro.Columnas[j+1], "KTNA") { // Si el siguiente campo es de ktna
						for z := 3; z <= j+2; z++ {
							cell := fmt.Sprintf("A%d", z)
							fileNuevo.SetCellValue(sheetName, cell, "Liquidacion")
							fileNuevo.SetCellStyle(sheetName, cell, cell, styleVertical)
						}
						colTotalNum = j + 2

						// fmt.Println("Proxima iteración es ktna: ", registro.Columnas[j+1])
						// fmt.Printf("Merge de %s a %s para liquidacion.\n", "A3", "A"+strconv.Itoa(j+2))
						fileNuevo.MergeCell(sheetName, "A3", "A"+strconv.Itoa(j+2))
					}

					// fmt.Printf("Celda de TOTAL: %s\n", "B"+strconv.Itoa(j+2))
					fileNuevo.SetCellStyle(sheetName, "B"+strconv.Itoa(j+2), "B"+strconv.Itoa(j+2), styleTotalesControl)             // Negrita y fondo gris para el total
					fileNuevo.SetCellStyle(sheetName, colLetter+strconv.Itoa(j+2), colLetter+strconv.Itoa(j+2), styleTotalesControl) // Negrita y fondo gris para el total

				}

			}

			for j, campo := range registro.Columnas { // Borrar prefijos "KTNA"
				if strings.Contains(campo, "KTNA") {
					campo = strings.Replace(campo, "KTNA ", "", -1)
					fileNuevo.SetCellValue(sheetName, "B"+strconv.Itoa(j+2), campo)
				}
			}

			fileNuevo.InsertRows(sheetName, colTotalNum+1, 1)  // Fila vacia
			fileNuevo.InsertRows(sheetName, colTotalNum2+1, 1) // Fila vacia

			fileNuevo.SetCellStyle(sheetName, colLetter+"2", colLetter+"2", styleColumnaControl) // Fondo negro en fila 2
			fileNuevo.SetColWidth(sheetName, "B", "B", 30)                                       // Amplia ancho de columna

		}

		// Escribir campo fijos
		fileNuevo.SetCellValue(sheetName, "A2", "Detalle")
		fileNuevo.SetCellValue(sheetName, "B2", "Conceptos")

		//Estilos fijos
		fileNuevo.SetCellStyle(sheetName, "A2", "A2", styleColumnaControl)
		fileNuevo.SetCellStyle(sheetName, "B2", "B2", styleColumnaControl)

		// Guardar archivo
		if err := fileNuevo.SaveAs(nombreSalida); err != nil {
			ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
			return "", err
		}

		return nombreSalida, nil
	}

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
			fileNuevo.SetCellStyle(sheetName, cell, cell, styleDefault)
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
				} else if strings.ToLower(campo.Tipo) == "numero" {
					if len(value) > 0 {
						valor, err := strconv.ParseFloat(value, 64)
						if err != nil {
							return "", err
						}
						fileNuevo.SetCellValue(sheetName, cell, valor)
					} else {
						fileNuevo.SetCellValue(sheetName, cell, value)
					}
					fileNuevo.SetCellStyle(sheetName, cell, cell, styleNumeroDecimal)
				} else {
					fileNuevo.SetCellValue(sheetName, cell, value)
					if campo.Ancho > 0 {
						fileNuevo.SetColWidth(sheetName, campo.Columna, campo.Columna, float64(campo.Ancho))
					}
					fileNuevo.SetCellStyle(sheetName, cell, cell, styleDefault)
				}
			}
		}
	}

	// Guardar archivo
	if err := fileNuevo.SaveAs(nombreSalida); err != nil {
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

func formatearPeriodoLiq(s string) string {
	strMes := s[4:]
	intMes, err := strconv.Atoi(strMes)
	if err != nil {
		fmt.Println("Error en formatearPeriodoLiq: ", err.Error())
	}
	meses := []string{"ene", "feb", "mar", "abr", "may", "jun", "jul", "ago", "sep", "oct", "nov", "dic"}
	for i, mes := range meses {
		if (intMes - 1) == i {
			return fmt.Sprintf("%s-%s", mes, s[:4])
		}
	}
	return s
}
