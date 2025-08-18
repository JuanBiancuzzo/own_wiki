package dependencias

import (
	"fmt"
	"strings"
)

type Condicion struct {
	expresionCondicion  string
	clavesOrdenadas     []string
	claveRepresentativa map[string]string
}

func NewCondicion(parClaveRepresentacion map[string]string) Condicion {
	expresionesClaves := make([]string, len(parClaveRepresentacion))
	clavesOrdenadas := make([]string, len(parClaveRepresentacion))

	contador := 0
	for clave := range parClaveRepresentacion {
		clavesOrdenadas[contador] = clave
		expresionesClaves[contador] = fmt.Sprintf("%s = ?", clave)
		contador++
	}

	return Condicion{
		expresionCondicion:  strings.Join(expresionesClaves, " AND "),
		clavesOrdenadas:     clavesOrdenadas,
		claveRepresentativa: parClaveRepresentacion,
	}
}

func (c Condicion) Expresion(datosRepresentativos map[string]string, variablePorClave map[string]TipoVariable) (string, []any, error) {
	cantidadDatos := len(c.clavesOrdenadas)
	datos := make([]any, cantidadDatos)

	for i, clave := range c.clavesOrdenadas {
		claveRepresentativa := c.claveRepresentativa[clave]

		if representacion, ok := datosRepresentativos[claveRepresentativa]; !ok {
			return "", datos, fmt.Errorf("no se paso los datos necesarios para la condicion dado por la clave %s", clave)

		} else if tipo, ok := variablePorClave[clave]; !ok {
			return "", datos, fmt.Errorf("no se consiguio el tipo dado por la clave %s", clave)

		} else if dato, err := tipo.ValorDeString(representacion); err != nil {
			return "", datos, fmt.Errorf("no se pudo obtener el valor real dada la representacion (%s), con error: %v", representacion, err)

		} else {
			datos[i] = dato
		}
	}

	if cantidadDatos > 0 {
		return fmt.Sprintf("WHERE %s", c.expresionCondicion), datos, nil
	}
	return "", datos, nil
}
