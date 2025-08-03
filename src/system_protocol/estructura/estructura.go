package estructura

import (
	"fmt"
	"strconv"
)

type Etapa string

const (
	ETAPA_SIN_EMPEZAR = "SinEmpezar"
	ETAPA_EMPEZADO    = "Empezado"
	ETAPA_AMPLIAR     = "Ampliar"
	ETAPA_TERMINADO   = "Terminado"
)

func NumeroODefault(representacion string, valorDefault int) int {
	if nuevoValor, err := strconv.Atoi(representacion); err == nil {
		return nuevoValor
	} else {
		return valorDefault
	}
}

func BooleanoODefault(representacion string, valorDefault bool) bool {
	switch representacion {
	case "true":
		return true
	case "false":
		return false
	default:
		return valorDefault
	}
}

func EtapaODefault(representacion string, valorDefault Etapa) Etapa {
	var etapa Etapa
	switch representacion {
	case "sin-empezar":
		etapa = ETAPA_SIN_EMPEZAR
	case "empezado":
		etapa = ETAPA_EMPEZADO
	case "ampliar":
		etapa = ETAPA_AMPLIAR
	case "terminado":
		etapa = ETAPA_TERMINADO
	default:
		etapa = valorDefault
	}
	return etapa
}

func ObtenerEtapa(representacionEtapa string) (Etapa, error) {
	var etapa Etapa
	switch representacionEtapa {
	case "sin-empezar":
		etapa = ETAPA_SIN_EMPEZAR
	case "empezado":
		etapa = ETAPA_EMPEZADO
	case "ampliar":
		etapa = ETAPA_AMPLIAR
	case "terminado":
		etapa = ETAPA_TERMINADO
	default:
		return ETAPA_SIN_EMPEZAR, fmt.Errorf("el tipo de etapa (%s) no es uno de los esperados", representacionEtapa)
	}

	return etapa, nil
}
