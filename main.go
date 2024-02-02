package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

type Persona struct {
	Id     string `json:"id"`
	Nombre string `json:"nombre"`
	Edad   string `json:"edad"`
}

func main() {

	// Leer archivo .txt para obtener los nombres de las columnas
	// columnas, err := leerTemplate("./templates/sindicato2.txt")
	// if err != nil {
	// 	panic(err.Error())
	// }

	// Leer archivo excel template
	cols, err := leerExcelTemplate("./templates/camioneros.xlsx")
	if err != nil {
		panic(err.Error())
	}
	// fmt.Println("Columnas del template: ", cols)

	// Leer archivo json
	var personas []map[string]interface{}
	contenido, err := os.ReadFile("./datos.json")
	if err != nil {
		panic(err.Error())
	}
	err = json.Unmarshal(contenido, &personas)
	if err != nil {
		panic(err.Error())
	}

	// // Crear archivo Excel
	err = crearExcel(cols, personas, "./salida/salida5.xlsx")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Datos insertados en salida5.xlsx")
}

func leerTemplate(path string) ([]string, error) {
	// Abrir el archivo de templates
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(content), ","), nil
}

func crearExcel(columnas []string, data []map[string]interface{}, nombreArchivo string) error {
	f := excelize.NewFile()
	sheetName := "MiHoja1"
	f.SetSheetName("Sheet1", sheetName)

	// Escribir encabezados en el Excel
	for colIndex, columna := range columnas {
		columnaMin := strings.ToLower(columna)
		colLetter, _ := excelize.ColumnNumberToName(colIndex + 1)
		cell := colLetter + "1"
		f.SetCellValue(sheetName, cell, columnaMin)
	}

	// Escribir datos en el archivo Excel
	for rowIndex, persona := range data { // row ----> 1 Gabi CABA 27
		for colIndex, columna := range columnas {
			columnaMin := strings.ToLower(columna)
			colLetter, _ := excelize.ColumnNumberToName(colIndex + 1)
			cell := colLetter + strconv.Itoa(rowIndex+2)
			valor, ok := persona[columnaMin]
			if ok {
				f.SetCellValue(sheetName, cell, valor)
			}
			// fmt.Printf("sheetName: %s cell: %s value: %s\n", sheetName, cell, valor)
		}
	}

	// Guardar archivo
	return f.SaveAs(nombreArchivo)
}

func leerExcelTemplate(path string) ([]string, error) {
	// Abrir archivo Excel
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}

	// Nombre de la hoja
	sheetName := f.GetSheetName(0)

	// Obtener las celdas de la primera fila
	filas, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	// Extraer valores de la primera fila
	var primeraFila []string
	primeraFila = append(primeraFila, filas[0]...)

	return primeraFila, nil
}
