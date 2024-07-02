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

func CargarExcel(db *sql.DB, idLogDetalle int, proceso modelos.Proceso, data []modelos.Registro, nombreSalida string, path string, tipo_ejecucion string, infoText string, version int) (string, error) {

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

	estilos := ObtenerEstilos(fileNuevo)

	if tipo_ejecucion == "control" {

		var saltos []int

		columnas_combinadas := combinarColumnas(data)

		// fmt.Println("\nColumnas combinadas: ", columnas_combinadas)

		for i, registro := range data {
			colLetter := ObtenerLetra(i + 3) // Almacena la columna para este registro
			grupoPrevio := ""
			grupoActual := ""
			filaInicialGrupo := 3

			// fmt.Println(len(registro.Columnas))
			for j, campo := range columnas_combinadas {

				if i == 0 {
					grupoActual = ExtraerGrupo(campo)

					// Eliminar marca en el nombre del campo
					campo = EliminarPrefijo(campo, fmt.Sprintf("<%s>", grupoActual))
					campo = EliminarPrefijo(campo, "*C")
					campo = EliminarPrefijo(campo, "*")

					if grupoPrevio != "" && grupoPrevio != grupoActual {
						// El grupo cambia, fusiono las celdas de la columna "A"
						// fmt.Printf("Inicio: %v, Fin: %v\n", filaInicialGrupo, j+1)
						fileNuevo.MergeCell(sheetName, fmt.Sprintf("A%d", filaInicialGrupo), fmt.Sprintf("A%d", j+1))
						fileNuevo.SetCellValue(sheetName, fmt.Sprintf("A%d", filaInicialGrupo), strings.ToUpper(grupoPrevio))
						fileNuevo.SetCellStyle(sheetName, fmt.Sprintf("A%d", filaInicialGrupo), fmt.Sprintf("A%d", filaInicialGrupo), estilos.StyleVertical)
						filaInicialGrupo = j + 2
						saltos = append(saltos, j+2+len(saltos))
					}

					// Escribir claves
					cellKey := "B" + strconv.Itoa(j+2)
					fileNuevo.SetCellValue(sheetName, cellKey, campo)
					fileNuevo.SetCellStyle(sheetName, cellKey, cellKey, estilos.StyleEncabezadoControl) // Fondo gris para el campo
					fileNuevo.SetRowHeight(sheetName, j+2, 25)                                          // Amplia alto de la fila
				}

				// Escribir valores
				value := registro.Valores[strings.ToUpper(columnas_combinadas[j])]
				cellValue := colLetter + strconv.Itoa(j+2)

				if strings.ToUpper(campo) == "PERIODOLIQ" {
					if i == len(data)-1 && registro.Valores["NUM_VERSION"] == nil {
						value = formatearPeriodoLiq(value.(string)) + fmt.Sprintf(" (%d)", version)
					} else {
						value = formatearPeriodoLiq(value.(string)) + fmt.Sprintf(" (%d)", registro.Valores["NUM_VERSION"].(int))
						delete(registro.Valores, "NUM_VERSION")
					}
					fileNuevo.SetCellStyle(sheetName, cellValue, cellValue, estilos.StyleColumnaControl)
				}

				switch v := value.(type) {
				case []uint8:
					valueStr := string(v)
					valueFloat, err := strconv.ParseFloat(valueStr, 64)
					if err != nil {
						fmt.Println(err.Error())
					}
					fileNuevo.SetCellValue(sheetName, cellValue, valueFloat)
					fileNuevo.SetCellStyle(sheetName, cellValue, cellValue, estilos.StyleMoneda)
				case string:
					fileNuevo.SetCellValue(sheetName, cellValue, v)
					fileNuevo.SetCellStyle(sheetName, cellValue, cellValue, estilos.StyleAligned)
				case int64:
					fileNuevo.SetCellValue(sheetName, cellValue, v)
					fileNuevo.SetCellStyle(sheetName, cellValue, cellValue, estilos.StyleAligned)
				case float64:
					if strings.Contains(strings.ToUpper(campo), "LEGAJO") {
						fileNuevo.SetCellValue(sheetName, cellValue, v)
						fileNuevo.SetCellStyle(sheetName, cellValue, cellValue, estilos.StyleAligned)
					} else {
						fileNuevo.SetCellValue(sheetName, cellValue, v)
						fileNuevo.SetCellStyle(sheetName, cellValue, cellValue, estilos.StyleMoneda)
					}
				case nil:
					fileNuevo.SetCellValue(sheetName, cellValue, "")
					fileNuevo.SetCellStyle(sheetName, cellValue, cellValue, estilos.StyleAligned)
				default:
					fileNuevo.SetCellValue(sheetName, cellValue, "defValue")
					fmt.Printf("Tipo de dato en %s: %T\n", campo, value)

				}

				if strings.Contains(columnas_combinadas[j], "*C") {
					fileNuevo.SetCellStyle(sheetName, "B"+strconv.Itoa(j+2), "B"+strconv.Itoa(j+2), estilos.StyleControlCeleste)             // Negrita y fondo gris para el total
					fileNuevo.SetCellStyle(sheetName, colLetter+strconv.Itoa(j+2), colLetter+strconv.Itoa(j+2), estilos.StyleControlCeleste) // Negrita y fondo gris para el total
				} else if strings.Contains(columnas_combinadas[j], "*") {
					fileNuevo.SetCellStyle(sheetName, "B"+strconv.Itoa(j+2), "B"+strconv.Itoa(j+2), estilos.StyleTotalesControl)             // Negrita y fondo gris para el total
					fileNuevo.SetCellStyle(sheetName, colLetter+strconv.Itoa(j+2), colLetter+strconv.Itoa(j+2), estilos.StyleTotalesControl) // Negrita y fondo gris para el total
				}

				if i == 0 {
					grupoPrevio = grupoActual

					if j == len(columnas_combinadas)-1 {
						// Fusionar celdas del ultimo grupo
						if grupoPrevio != "" {
							// fmt.Printf("Inicio: %v, Fin: %v\n", filaInicialGrupo, j+2)
							fileNuevo.MergeCell(sheetName, fmt.Sprintf("A%d", filaInicialGrupo), fmt.Sprintf("A%d", j+2))
							fileNuevo.SetCellValue(sheetName, fmt.Sprintf("A%d", filaInicialGrupo), strings.ToUpper(grupoPrevio))
							fileNuevo.SetCellStyle(sheetName, fmt.Sprintf("A%d", filaInicialGrupo), fmt.Sprintf("A%d", filaInicialGrupo), estilos.StyleVertical)
						}
					}
				}

			}

			fileNuevo.SetCellStyle(sheetName, colLetter+"2", colLetter+"2", estilos.StyleColumnaControl) // Fondo negro en fila 2
			fileNuevo.SetColWidth(sheetName, "B", "B", 30)                                               // Amplia ancho de columna

		}

		// Insertar separacion entre filas
		for _, fila := range saltos {
			fileNuevo.InsertRows(sheetName, fila, 1)
		}

		// Escribir campo fijos
		fileNuevo.SetCellValue(sheetName, "A2", "Detalle")
		fileNuevo.SetCellValue(sheetName, "B2", "Conceptos")

		//Estilos fijos
		fileNuevo.SetCellStyle(sheetName, "A2", "A2", estilos.StyleColumnaControl)
		fileNuevo.SetCellStyle(sheetName, "B2", "B2", estilos.StyleColumnaControl)

		if infoText != "" {

			// Crear nueva hoja llamada "INFO"
			_, err := fileNuevo.NewSheet("INFO")
			if err != nil {
				ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
				return "", err
			}

			fileNuevo.SetColWidth("INFO", "A", "A", 20)
			fileNuevo.SetColWidth("INFO", "B", "B", 90)
			fileNuevo.SetColWidth("INFO", "C", "C", 90)

			// Extraer los segmentos del texto
			segments, err := extractGroupsAndSegments(infoText)
			if err != nil {
				ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
				return "", err
			}

			charPerLine := 100
			// Escribir los grupos y sus contenidos en la hoja "INFO"
			for i, pair := range segments {
				groupNameCell := fmt.Sprintf("A%d", i+1)
				segmentCell := fmt.Sprintf("B%d", i+1)
				fileNuevo.SetCellValue("INFO", groupNameCell, pair[0])
				fileNuevo.SetCellStyle("INFO", groupNameCell, groupNameCell, estilos.StyleColumnaInfo)
				if strings.Contains(pair[1], "::") {
					// Tiene traduccion
					partes := strings.Split(pair[1], "::")
					if len(partes) < 2 {
						fmt.Println("INFO no tenia ::")
						ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
						return "", err
					}
					fileNuevo.SetCellValue("INFO", segmentCell, partes[0])
					fileNuevo.SetCellStyle("INFO", segmentCell, segmentCell, estilos.StyleValorInfo)
					fileNuevo.SetCellValue("INFO", fmt.Sprintf("C%d", i+1), partes[1])
					fileNuevo.SetCellStyle("INFO", fmt.Sprintf("C%d", i+1), fmt.Sprintf("C%d", i+1), estilos.StyleValorInfo)
				} else {
					fileNuevo.SetCellValue("INFO", segmentCell, pair[1])
					fileNuevo.SetCellStyle("INFO", segmentCell, segmentCell, estilos.StyleValorInfo)
				}

				// Calcular y establecer el alto de la fila
				rowHeight := calculateRowHeight(pair[1], charPerLine)
				if rowHeight < 15.0 {
					fileNuevo.SetRowHeight("INFO", i+1, 15.0)
				} else {
					fileNuevo.SetRowHeight("INFO", i+1, rowHeight)
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

	if strings.ToLower(plantilla.Cabecera.Sentido_encabezado) == "vertical" {
		// Escribir verticalmente encabezados en el Excel
		for _, campo := range plantilla.Campos {
			cell := "B" + campo.Columna
			fileNuevo.SetCellValue(sheetName, cell, campo.Titulo)
		}
	} else if strings.ToLower(plantilla.Cabecera.Encabezados) != "no" {
		// Escribir horizontalmente encabezados en el Excel
		for _, campo := range plantilla.Campos {
			cell := campo.Columna + "1"
			fileNuevo.SetCellValue(sheetName, cell, campo.Titulo)
			fileNuevo.SetCellStyle(sheetName, cell, cell, estilos.StyleDefaultCabecera)
		}
		if strings.ToLower(plantilla.Cabecera.Estilo) == "nomina" {
			fileNuevo.SetRowStyle(sheetName, 1, 1, estilos.StyleEncabezadoNomina)
			fileNuevo.SetRowHeight(sheetName, 1, 50)
		}
	}

	// Escribir datos en el archivo Excel
	for i, registro := range data {
		for _, campo := range plantilla.Campos {
			var value string
			var cell string
			// Validaciones

			if campo.Nombre == "" && campo.Tipo != "suma" {
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
			} else if strings.ToLower(campo.Tipo) == "orden" {
				value = fmt.Sprintf("%d", i+1)
			} else {
				campo.Nombre = strings.ToUpper(campo.Nombre)
				val := registro.Valores[campo.Nombre]

				if campo.Columna == "" {
					return "", fmt.Errorf("JSON: el campo %s no tiene columna", campo.Nombre)
				}
				if campo.Inicio != 0 || campo.Fin != 0 {
					return "", fmt.Errorf("JSON: el campo %s no debe tener 'inicio' ni 'fin'", campo.Nombre)
				}
				if campo.Tipo == "fecha" {
					if campo.Formato != "DD/MM/YYYY" && campo.Formato != "MM/YYYY" && campo.Formato != "DD-MM-YYYY" && campo.Formato != "YYYYMMDD" {
						return "", fmt.Errorf("JSON: formato desconocido para %s", campo.Nombre)
					}
				}
				if campo.Formato == "DD/MM/YYYY" || campo.Formato == "DD-MM-YYYY" {
					if campo.Tipo != "fecha" {
						return "", fmt.Errorf("JSON: el campo %s debe ser de tipo fecha", campo.Nombre)
					}
				}

				// if strings.ToLower(campo.Nombre) == "afiliado" {
				// 	fmt.Printf("Campo: %s Valor: %v Tipo: %s\n", campo.Nombre, val, reflect.TypeOf(val))
				// }
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
									if strings.TrimSpace(element.Id) == strings.TrimSpace(v) {
										value += element.Nombre
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
					} else if campo.Formato == "MM/YYYY" {
						value = formatearFecha(v, campo.Formato)
					} else if strings.ToLower(campo.Tipo) == "condicional" {
						numRegex := regexp.MustCompile(`^\s{8}$`)
						condiciones := strings.Split(campo.Formato, "/")
						if numRegex.MatchString(v) {
							value = condiciones[1]
						} else if strings.TrimSpace(v) == "0" {
							value = condiciones[1]
						} else {
							value = condiciones[0]
						}
					} else {
						value += strings.TrimSpace(v)
					}
				case int64:
					if strings.ToLower(campo.Tipo) == "condicional" {
						condiciones := strings.Split(campo.Formato, "/")
						if v == 0 {
							value = condiciones[1]
						} else {
							value = condiciones[0]
						}
					} else {
						value = strings.TrimSpace(fmt.Sprintf("%v", v))
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
					if strings.ToLower(campo.Tipo) == "condicional" {
						condiciones := strings.Split(campo.Formato, "/")
						value = condiciones[1]
					} else {
						value = ""
					}

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
				cell = colLetter + campo.Columna
				if strings.ToLower(campo.Tipo) == "moneda" {
					valor, _ := strconv.ParseFloat(value, 64)
					fileNuevo.SetCellValue(sheetName, cell, valor)
					fileNuevo.SetCellStyle(sheetName, cell, cell, estilos.StyleMoneda)
				} else if strings.ToLower(campo.Tipo) == "int" {
					valor, _ := strconv.Atoi(value)
					fileNuevo.SetCellValue(sheetName, cell, valor)
					fileNuevo.SetCellStyle(sheetName, colLetter, campo.Columna, estilos.StyleNumero)
				} else if strings.ToLower(campo.Tipo) == "numero decimal" {
					valor, _ := strconv.ParseFloat(value, 64)
					fileNuevo.SetCellStyle(sheetName, colLetter, campo.Columna, estilos.StyleNumeroDecimal)
					fileNuevo.SetCellValue(sheetName, cell, valor)
				} else {
					fileNuevo.SetCellValue(sheetName, cell, value)
				}
			} else {
				if strings.ToLower(plantilla.Cabecera.Encabezados) != "no" {
					cell = campo.Columna + fmt.Sprintf("%v", i+2)
				} else {
					cell = campo.Columna + fmt.Sprintf("%v", i+1)
				}
				if strings.ToLower(campo.Tipo) == "moneda" {
					valor, _ := strconv.ParseFloat(value, 64)
					fileNuevo.SetCellValue(sheetName, cell, valor)
					fileNuevo.SetCellStyle(sheetName, cell, cell, estilos.StyleMoneda)
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
					fileNuevo.SetCellStyle(sheetName, cell, cell, estilos.StyleNumero)
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
					fileNuevo.SetCellStyle(sheetName, cell, cell, estilos.StyleNumeroDecimal)
				} else {
					fileNuevo.SetCellValue(sheetName, cell, value)
					if campo.Ancho > 0 {
						fileNuevo.SetColWidth(sheetName, campo.Columna, campo.Columna, float64(campo.Ancho))
					}
					fileNuevo.SetCellStyle(sheetName, cell, cell, estilos.StyleDefault)
				}
			}
		}
		// Para el ultimo campo
		// if i == len(data)-1 {
		// 	if len(options) > 0 {
		// 		fmt.Println("Se creo la lista con rango: ", fmt.Sprintf("K2:K%v", i+2))
		// 		// dv.SetSqref(fmt.Sprintf("K2:K%v", i+2))
		// 		fileNuevo.AddDataValidation(sheetName, dv)
		// 	}
		// }
	}

	var options []string
	for _, va := range plantilla.Variables {
		if va.Nombre == "lista desplegable" {
			options = va.Valores
		}
	}
	// validationString := "\"" + strings.Join(options, ";") + "\""
	dv := excelize.NewDataValidation(true)
	dv.Type = "list"
	dv.Sqref = "K2:K82"
	dv.SetDropList(options)
	dv.ShowDropDown = true
	err := fileNuevo.AddDataValidation(sheetName, dv)
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", err
	}

	// Guardar archivo
	if err = fileNuevo.SaveAs(nombreSalida); err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		return "", err
	}

	return nombreSalida, nil
}

func ExtraerGrupo(campo string) string {
	inicio := strings.Index(campo, "<")
	fin := strings.Index(campo, ">")
	if inicio != -1 && fin != -1 && fin > inicio {
		return campo[inicio+1 : fin]
	}
	return ""
}

func EliminarPrefijo(campo, prefijo string) string {
	return strings.Replace(campo, prefijo, "", -1)
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
	} else if formato == "DDMMYYYY" {
		partes = append(partes, s[:4], s[4:6], s[6:])
		strFinal = partes[2] + partes[1] + partes[0]
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

// Función para procesar el texto y extraer los grupos y sus contenidos
func extractGroupsAndSegments(text string) ([][2]string, error) {
	re := regexp.MustCompile(`\[(.*?)\]`)
	matches := re.FindAllStringSubmatch(text, -1)

	if len(matches) == 0 {
		return nil, fmt.Errorf("no groups found in text")
	}

	var result [][2]string
	for i, match := range matches {
		groupName := match[1]
		start := strings.Index(text, match[0]) + len(match[0])
		end := len(text)
		if i+1 < len(matches) {
			end = strings.Index(text, matches[i+1][0])
		}
		segment := strings.TrimSpace(text[start:end])
		result = append(result, [2]string{groupName, segment})
	}
	return result, nil
}

// Función para calcular el alto de la fila basado en el contenido de la celda
func calculateRowHeight(content string, charsPerLine int) float64 {
	// Separar el contenido en líneas considerando los saltos de línea intencionales
	lines := strings.Split(content, "\n")
	totalLines := 0

	for _, line := range lines {
		// Contar el número de líneas necesarias para este fragmento de texto
		lineLength := len(line)
		lineCount := (lineLength / charsPerLine) + 1
		totalLines += lineCount
	}

	// Asumir 15 píxeles por línea, ajustar según tus necesidades
	rowHeight := 15.0 * float64(totalLines)
	return rowHeight
}
