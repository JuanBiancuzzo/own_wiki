package procesar

import (
	"fmt"
	"strconv"
)

type TipoNota string

const (
	TN_FACULTAD      = "Facultad"
	TN_COLECCION     = "Coleccion"
	TN_CURSO         = "Curso"
	TN_INVESTIGACION = "Investigacion"
	TN_PROYECTO      = "Proyecto"
)

type Etapa string

const (
	ETAPA_SIN_EMPEZAR = "SinEmpezar"
	ETAPA_EMPEZADO    = "Empezado"
	ETAPA_AMPLIAR     = "Ampliar"
	ETAPA_TERMINADO   = "Terminado"
)

type TipoDistribucion string

const (
	DISTRIBUCION_DISCRETA     = "Discreta"
	DISTRIBUCION_CONTINUA     = "Continua"
	DISTRIBUCION_MULTIVARIADA = "Multivariada"
)

type ParteCuatrimestre string

const (
	CUATRIMESTRE_PRIMERO = "Primero"
	CUATRIMESTRE_SEGUNDO = "Segundo"
)

type TipoMateria string

const (
	MATERIA_REAL        = "Materia"
	MATERIA_EQUIVALENTE = "Equivalente"
)

type TipoCurso string

const (
	CURSO_ONLINE     = "Online"
	CURSO_PRESENCIAL = "Presencial"
)

type TipoEscritorPaper string

const (
	PAPER_EDITOR = "Editor"
	PAPER_AUTOR  = "Autor"
)

func NumeroODefault(representacion string, valorDefault int) int {
	if nuevoValor, err := strconv.Atoi(representacion); err == nil {
		return nuevoValor
	} else {
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

func ObtenerTipoDistribucion(representacion string) (TipoDistribucion, error) {
	var tipoDistribucion TipoDistribucion
	switch representacion {
	case "discreta":
		tipoDistribucion = DISTRIBUCION_DISCRETA
	case "continua":
		tipoDistribucion = DISTRIBUCION_CONTINUA
	case "multivariada":
		tipoDistribucion = DISTRIBUCION_MULTIVARIADA
	default:
		return DISTRIBUCION_DISCRETA, fmt.Errorf("el tipo de distribucion (%s) no es uno de los esperados", representacion)
	}

	return tipoDistribucion, nil
}

func ObtenerCuatrimestreParte(representacionCuatri string) (int, ParteCuatrimestre, error) {
	var anio int
	var cuatriNum int
	var cuatri ParteCuatrimestre

	if _, err := fmt.Sscanf(representacionCuatri, "%dC%d", &anio, &cuatriNum); err != nil {
		return anio, cuatri, fmt.Errorf("el tipo de anio-cuatri (%s) no es uno de los esperados", representacionCuatri)
	}

	switch cuatriNum {
	case 1:
		cuatri = CUATRIMESTRE_PRIMERO
	case 2:
		cuatri = CUATRIMESTRE_SEGUNDO
	default:
		return anio, cuatri, fmt.Errorf("el cuatri dado por %d no es posible representar", cuatriNum)
	}

	return anio, cuatri, nil
}
