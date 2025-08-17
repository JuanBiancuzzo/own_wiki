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

func NewCondicion(listaClaves []string) Condicion {
	expresionesClaves := make([]string, len(listaClaves))
	claveRepresentativa := make(map[string]string)

	for i, clave := range listaClaves {
		expresionesClaves[i] = fmt.Sprintf("%s = ?", clave)

	}

	return Condicion{
		expresionCondicion:  strings.Join(expresionesClaves, " AND "),
		clavesOrdenadas:     listaClaves,
		claveRepresentativa: claveRepresentativa,
	}
}

func (c Condicion) Expresion(datosRepresentativos map[string]string, variablePorClave map[string]TipoVariable) (string, []any, error) {
	datos := make([]any, len(c.clavesOrdenadas))
	for i, clave := range c.clavesOrdenadas {
		claveRepresentativa := c.claveRepresentativa[clave]

		if representacion, ok := datosRepresentativos[claveRepresentativa]; !ok {
			return "", datos, fmt.Errorf("no se paso los datos necesarios para la condicion dado por la clave %s", clave)

		} else if tipo, ok := variablePorClave[clave]; !ok {
			return "", datos, fmt.Errorf("no se consiguio el tipo dado por la clave %s", clave)

		} else if dato, err := tipo.ValorDeString(representacion); err != nil {
			return "", datos, fmt.Errorf("no se pudo obtener el valor real dada la representacion (%s)", representacion)

		} else {
			datos[i] = dato
		}
	}

	return c.expresionCondicion, datos, nil
}
