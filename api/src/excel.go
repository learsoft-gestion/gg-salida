package src

import (
	"Nueva/modelos"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func CargarExcel(db *sql.DB, idLogDetalle int, proceso modelos.Proceso, data []modelos.Registro, nombreSalida string) error {
	// Construir path
	path := fmt.Sprintf("../templates/%s.xlsx", proceso.Nombre)

	// Abrir archivo Excel
	f, err := excelize.OpenFile(path)
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
	}

	// Nombre de la hoja
	sheetName := f.GetSheetName(0)

	// Obtener las celdas de la primera fila
	filas, err := f.GetRows(sheetName)
	if err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
	}
	fileNuevo := excelize.NewFile()

	// var estilos []int
	for i := range filas[0] {
		colLetter, _ := excelize.ColumnNumberToName(i + 1)
		cell := colLetter + "1"
		styles, err := f.GetCellStyle(sheetName, cell)
		if err != nil {
			ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
		}
		estilo, err := f.GetStyle(styles)
		if err != nil {
			ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
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
	for rowIndex, registro := range data { // row ----> 1 Gabi CABA 27
		for colIndex, columna := range primeraFila {
			// columnaMin := strings.ToLower(columna)
			colLetter, _ := excelize.ColumnNumberToName(colIndex + 1)
			cell := colLetter + strconv.Itoa(rowIndex+2)
			valor, ok := registro.Valores[columna]
			if ok {
				f.SetCellValue(sheetName, cell, valor)
			}
			// fmt.Printf("sheetName: %s cell: %s value: %s\n", sheetName, cell, valor)
		}
	}

	// Guardar archivo
	if err = f.SaveAs(nombreSalida); err != nil {
		ManejoErrores(db, idLogDetalle, proceso.Nombre, err)
	}

	return nil
}
