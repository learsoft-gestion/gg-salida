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

	err = leerCrearExcel("./templates/camioneros.xlsx", personas, "./salida/salida2.xlsx")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Datos insertados en salida2.xlsx")
}

func crearExcel(columnas []string, data []map[string]interface{}, estilos []int, nombreArchivo string) error {
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
	fmt.Println("Estilos[1]: ", estilos[1])
	err := f.SetRowStyle(sheetName, 1, 10, estilos[1])
	if err != nil {
		return err
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

func leerExcelTemplate(path string) ([]string, []int, error) {
	// Abrir archivo Excel
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, nil, err
	}

	// Nombre de la hoja
	sheetName := f.GetSheetName(0)

	// Obtener las celdas de la primera fila
	filas, err := f.GetRows(sheetName)
	if err != nil {
		return nil, nil, err
	}

	// var estilos []*excelize.Style
	var estilos []int
	for i := range filas[0] {
		colLetter, _ := excelize.ColumnNumberToName(i + 1)
		cell := colLetter + "1"
		styles, err := f.GetCellStyle(sheetName, cell)
		if err != nil {
			return nil, nil, err
		}
		// estilo, _ := f.GetStyle(styles)
		estilos = append(estilos, styles)
	}
	fmt.Println(estilos)

	// Extraer valores de la primera fila
	var primeraFila []string
	primeraFila = append(primeraFila, filas[0]...)

	return primeraFila, estilos, nil
}

func leerCrearExcel(path string, data []map[string]interface{}, nombre string) error {
	// Abrir archivo Excel
	f, err := excelize.OpenFile(path)
	if err != nil {
		return err
	}

	// Nombre de la hoja
	sheetName := f.GetSheetName(0)

	// Obtener las celdas de la primera fila
	filas, err := f.GetRows(sheetName)
	if err != nil {
		return err
	}
	fileNuevo := excelize.NewFile()

	// var estilos []int
	for i := range filas[0] {
		colLetter, _ := excelize.ColumnNumberToName(i + 1)
		cell := colLetter + "1"
		styles, err := f.GetCellStyle(sheetName, cell)
		if err != nil {
			return err
		}
		estilo, err := f.GetStyle(styles)
		if err != nil {
			return err
		}
		fileNuevo.NewStyle(estilo)
		fileNuevo.SetCellStyle(sheetName, colLetter, "1", styles)
		// estilos = append(estilos, styles)
	}

	// Extraer valores de la primera fila
	var primeraFila []string
	primeraFila = append(primeraFila, filas[0]...)

	fileNuevo.SetSheetName("Sheet1", sheetName)

	// Escribir encabezados en el Excel
	for colIndex, columna := range primeraFila {
		columnaMin := strings.ToLower(columna)
		colLetter, _ := excelize.ColumnNumberToName(colIndex + 1)
		cell := colLetter + "1"
		f.SetCellValue(sheetName, cell, columnaMin)
	}

	// err := fileNuevo.SetRowStyle(sheetName, 1, 10, estilos[1])
	// if err != nil {
	// 	return err
	// }

	// Escribir datos en el archivo Excel
	for rowIndex, persona := range data { // row ----> 1 Gabi CABA 27
		for colIndex, columna := range primeraFila {
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
	return f.SaveAs(nombre)
}
