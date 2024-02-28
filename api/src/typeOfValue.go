package src

import (
	"strconv"
)

// Verifica si el valor es NO nulo y de tipo slice de bytes []Uint8
func esSliceDeBytes(valor interface{}) bool {
	if valor == nil {
		return false
	} else if _, ok := valor.([]byte); ok {
		return true
	}
	// tipo := reflect.TypeOf(valor)
	// return tipo.Kind() == reflect.Slice && tipo.Elem().Kind() == reflect.Uint8
	return false
}

// Convierte un valor int o []Uint8 en float64
func valueToFloat(valor interface{}) float64 {

	if esSliceDeBytes(valor) {
		// Si es slice de bytes significa que es un numero muy grande e incluso decimal en algunos casos. Lo paso a float64
		valueBytes := valor.([]uint8)
		valueStr := string(valueBytes)
		valueOrigen, _ := strconv.ParseFloat(valueStr, 64)
		return valueOrigen

	} else {
		claveOrigenInt, ok := valor.(int)
		if ok {
			valueOrigen := float64(claveOrigenInt)
			return valueOrigen
		} else {
			// Si no es un slice de bytes ni un int significa que es un int64 por lo tanto lo paso a float64 para poder seguir con la comparacion.
			claveOrigenInt64 := valor.(int64)
			valueOrigen := float64(claveOrigenInt64)
			return valueOrigen
		}
	}
}
