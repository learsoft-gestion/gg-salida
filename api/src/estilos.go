package src

import (
	"Nueva/modelos"

	"github.com/xuri/excelize/v2"
)

func ObtenerEstilos(fileNuevo *excelize.File) modelos.Estilos {
	styleMoneda, _ := fileNuevo.NewStyle(&excelize.Style{
		NumFmt: 177, Font: &excelize.Font{Size: 9}, Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: false},
		Border: []excelize.Border{{Type: "top", Color: "000000", Style: 1}, {Type: "bottom", Color: "000000", Style: 1}, {Type: "left", Color: "000000", Style: 1}, {Type: "right", Color: "000000", Style: 1}},
	})
	styleNumero, _ := fileNuevo.NewStyle(&excelize.Style{NumFmt: 1, Border: []excelize.Border{{Type: "top", Color: "000000", Style: 1}, {Type: "bottom", Color: "000000", Style: 1}, {Type: "left", Color: "000000", Style: 1}, {Type: "right", Color: "000000", Style: 1}}})
	styleNumeroDecimal, _ := fileNuevo.NewStyle(&excelize.Style{NumFmt: 2, Border: []excelize.Border{{Type: "top", Color: "000000", Style: 1}, {Type: "bottom", Color: "000000", Style: 1}, {Type: "left", Color: "000000", Style: 1}, {Type: "right", Color: "000000", Style: 1}}})
	styleEncabezadoNomina, _ := fileNuevo.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 10, Color: "#FF0000"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#a7a7a7"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", WrapText: true},
		Border:    []excelize.Border{{Type: "top", Color: "000000", Style: 1}, {Type: "bottom", Color: "000000", Style: 1}, {Type: "left", Color: "000000", Style: 1}, {Type: "right", Color: "000000", Style: 1}},
	})
	styleEncabezadoControl, _ := fileNuevo.NewStyle(&excelize.Style{Font: &excelize.Font{Size: 10}, Border: []excelize.Border{{Type: "top", Color: "000000", Style: 1}, {Type: "bottom", Color: "000000", Style: 1}, {Type: "left", Color: "000000", Style: 1}, {Type: "right", Color: "000000", Style: 1}}, Fill: excelize.Fill{Type: "pattern", Color: []string{"#DCDCDC"}, Pattern: 1}, Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", WrapText: true}})
	styleColumnaControl, _ := fileNuevo.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Color: "#FFD3A7"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#000000"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
	})
	styleTotalesControl, _ := fileNuevo.NewStyle(&excelize.Style{
		NumFmt:    177,
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#DCDCDC"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", ReadingOrder: 0, Indent: 0, RelativeIndent: 0, ShrinkToFit: false, TextRotation: 0, WrapText: false},
		Border:    []excelize.Border{{Type: "top", Color: "000000", Style: 1}, {Type: "bottom", Color: "000000", Style: 1}, {Type: "left", Color: "000000", Style: 1}, {Type: "right", Color: "000000", Style: 1}},
	})
	styleControlCeleste, _ := fileNuevo.NewStyle(&excelize.Style{
		NumFmt:    177,
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#BDD7EE"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", ReadingOrder: 0, Indent: 0, RelativeIndent: 0, ShrinkToFit: false, TextRotation: 0, WrapText: false},
		Border:    []excelize.Border{{Type: "top", Color: "000000", Style: 1}, {Type: "bottom", Color: "000000", Style: 1}, {Type: "left", Color: "000000", Style: 1}, {Type: "right", Color: "000000", Style: 1}},
	})
	styleAligned, _ := fileNuevo.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", ReadingOrder: 0, Indent: 0, RelativeIndent: 0, ShrinkToFit: false, TextRotation: 0, WrapText: false},
		Border:    []excelize.Border{{Type: "top", Color: "000000", Style: 1}, {Type: "bottom", Color: "000000", Style: 1}, {Type: "left", Color: "000000", Style: 1}, {Type: "right", Color: "000000", Style: 1}},
	})
	styleDefaultCabecera, _ := fileNuevo.NewStyle(&excelize.Style{
		Border: []excelize.Border{{Type: "top", Color: "000000", Style: 1}, {Type: "bottom", Color: "000000", Style: 1}, {Type: "left", Color: "000000", Style: 1}, {Type: "right", Color: "000000", Style: 1}},
		Fill:   excelize.Fill{Type: "pattern", Color: []string{"#fdf59a"}, Pattern: 1},
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
	styleColumnaInfo, _ := fileNuevo.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#DCDCDC"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border:    []excelize.Border{{Type: "top", Color: "000000", Style: 1}, {Type: "bottom", Color: "000000", Style: 1}, {Type: "left", Color: "000000", Style: 1}, {Type: "right", Color: "000000", Style: 1}},
	})
	styleValorInfo, _ := fileNuevo.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Vertical: "center", WrapText: true},
		Border:    []excelize.Border{{Type: "top", Color: "000000", Style: 1}, {Type: "bottom", Color: "000000", Style: 1}, {Type: "left", Color: "000000", Style: 1}, {Type: "right", Color: "000000", Style: 1}},
	})

	return modelos.Estilos{StyleMoneda: styleMoneda, StyleNumero: styleNumero, StyleNumeroDecimal: styleNumeroDecimal, StyleEncabezadoNomina: styleEncabezadoNomina, StyleEncabezadoControl: styleEncabezadoControl, StyleColumnaControl: styleColumnaControl, StyleTotalesControl: styleTotalesControl, StyleControlCeleste: styleControlCeleste, StyleAligned: styleAligned, StyleDefaultCabecera: styleDefaultCabecera, StyleDefault: styleDefault, StyleVertical: styleVertical, StyleColumnaInfo: styleColumnaInfo, StyleValorInfo: styleValorInfo}
}
