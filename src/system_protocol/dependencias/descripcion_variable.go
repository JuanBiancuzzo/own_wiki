package dependencias

import (
	"fmt"
	"strconv"
	"strings"
)

type DescripcionVariable struct {
	Clave       string
	Descripcion any
}

func (v DescripcionVariable) ValorPorRepresentacion(representacion string) (any, error) {
	switch descripcion := v.Descripcion.(type) {
	case DescVariableSimple:
		switch descripcion.Tipo {
		case TVS_INT:
			return strconv.Atoi(representacion)
		case TVS_BOOL:
			switch strings.ToLower(strings.TrimSpace(representacion)) {
			case "true":
				return true, nil
			case "false":
				return false, nil
			default:
				return false, fmt.Errorf("la representacion no es un bool ya que es: %v", representacion)
			}
		case TVS_DATE:
			return representacion, nil
		}

	case DescVariableString:
		return representacion, nil
	case DescVariableEnum:
		return representacion, nil
	}
	return nil, fmt.Errorf("no es una variable que pueda tener un valor de una rerpresentacion")
}

func (v DescripcionVariable) ObtenerReferencia() (any, error) {
	switch informacion := v.Descripcion.(type) {
	case DescVariableSimple:
		switch informacion.Tipo {
		case TVS_INT:
			var numero int
			return &numero, nil
		case TVS_BOOL:
			var booleano bool
			return &booleano, nil
		case TVS_DATE:
			var representacion string
			return &representacion, nil
		}

	case DescVariableString:
		var representacion string
		return &representacion, nil
	case DescVariableEnum:
		var representacion string
		return &representacion, nil
	}

	return nil, fmt.Errorf("no es una variable que pueda tener un valor de una rerpresentacion")
}

func (v DescripcionVariable) Desreferenciar(referencia any) (any, error) {
	switch informacion := v.Descripcion.(type) {
	case DescVariableSimple:
		switch informacion.Tipo {
		case TVS_INT:
			if numeroRef, ok := referencia.(*int); !ok {
				return nil, fmt.Errorf("se esperaba que fuera un numero, pero no lo es")
			} else {
				return *numeroRef, nil
			}
		case TVS_BOOL:
			if booleanRef, ok := referencia.(*bool); !ok {
				return nil, fmt.Errorf("se esperaba que fuera un boolean, pero no lo es")
			} else {
				return *booleanRef, nil
			}
		case TVS_DATE:
			if stringRef, ok := referencia.(*string); !ok {
				return nil, fmt.Errorf("se esperaba que fuera un date, pero no lo es")
			} else {
				return *stringRef, nil
			}
		}

	case DescVariableString:
		if stringRef, ok := referencia.(*string); !ok {
			return nil, fmt.Errorf("se esperaba que fuera un string, pero no lo es")
		} else {
			return *stringRef, nil
		}
	case DescVariableEnum:
		if stringRef, ok := referencia.(*string); !ok {
			return nil, fmt.Errorf("se esperaba que fuera un enum, pero no lo es")
		} else {
			return *stringRef, nil
		}
	}

	return nil, fmt.Errorf("no es una variable que pueda tener un valor de una rerpresentacion")
}

type DescVariableSimple struct {
	Tipo           TipoVariableSimple
	Representativo bool
	Necesario      bool
}

func NewDescVariableSimple(tipo TipoVariableSimple, representativo bool, clave string, necesario bool) DescripcionVariable {
	return DescripcionVariable{
		Clave: clave,
		Descripcion: DescVariableSimple{
			Tipo:           tipo,
			Representativo: representativo,
			Necesario:      necesario,
		},
	}
}

type DescVariableString struct {
	Representativo bool
	Necesario      bool
	Largo          uint
}

func NewDescVariableString(representativo bool, clave string, largo uint, necesario bool) DescripcionVariable {
	return DescripcionVariable{
		Clave: clave,
		Descripcion: DescVariableString{
			Representativo: representativo,
			Necesario:      necesario,
			Largo:          largo,
		},
	}
}

type DescVariableEnum struct {
	Representativo bool
	Necesario      bool
	Valores        []string
}

func NewDescVariableEnum(representativo bool, clave string, valores []string, necesario bool) DescripcionVariable {
	return DescripcionVariable{
		Clave: clave,
		Descripcion: DescVariableEnum{
			Representativo: representativo,
			Necesario:      necesario,
			Valores:        valores,
		},
	}
}

type DescVariableReferencia struct {
	Representativo bool
	Tablas         []string
}

func NewDescVariableReferencia(representativo bool, clave string, tablas []string) DescripcionVariable {
	return DescripcionVariable{
		Clave: clave,
		Descripcion: DescVariableReferencia{
			Representativo: representativo,
			Tablas:         tablas,
		},
	}
}

type DescVariableArrayReferencia struct {
	ClaveSelf   string
	TablaCreada string
}

func NewDescVariableArrayReferencias(clave, claveSelf, tablaCreada string) DescripcionVariable {
	return DescripcionVariable{
		Clave: clave,
		Descripcion: DescVariableArrayReferencia{
			ClaveSelf:   claveSelf,
			TablaCreada: tablaCreada,
		},
	}
}
