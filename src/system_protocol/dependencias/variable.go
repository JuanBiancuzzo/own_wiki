package dependencias

import (
	"fmt"
	"strconv"
	"strings"
)

type InformacionVariable any

type Variable struct {
	Clave       string
	Informacion InformacionVariable
}

func (v Variable) ObtenerParametroSQL() []string {
	var parametros []string

	switch informacion := v.Informacion.(type) {
	case VariableSimple:
		parametros = []string{fmt.Sprintf("%s %s", v.Clave, informacion.TipoSQL())}
	case VariableString:
		parametros = []string{fmt.Sprintf("%s %s", v.Clave, informacion.TipoSQL())}
	case VariableEnum:
		parametros = []string{fmt.Sprintf("%s %s", v.Clave, informacion.TipoSQL())}

	case VariableReferencia:
		cantidad := min(2, len(informacion.Tablas))
		parametros = make([]string, cantidad)
		if len(informacion.Tablas) > 1 {
			valoresRep := make([]string, len(informacion.Tablas))
			for i, tabla := range informacion.Tablas {
				valoresRep[i] = fmt.Sprintf("\"%s\"", tabla.NombreTabla)
			}
			parametros[0] = fmt.Sprintf("tipo%s ENUM(%s)", v.Clave, strings.Join(valoresRep, ", "))
		}

		parametros[cantidad-1] = fmt.Sprintf("%s INT", v.Clave)

	case VariableArrayReferencia:
		// no es necesario, ya que no tiene representacion en sql
	}

	return parametros
}

func (v Variable) ValorPorRepresentacion(representacion string) (any, error) {
	switch informacion := v.Informacion.(type) {
	case VariableSimple:
		switch informacion.Tipo {
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

	case VariableString:
		return representacion, nil
	case VariableEnum:
		return representacion, nil
	}
	return nil, fmt.Errorf("no es una variable que pueda tener un valor de una rerpresentacion")
}

func (v Variable) ObtenerReferencia() (any, error) {
	switch informacion := v.Informacion.(type) {
	case VariableSimple:
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

	case VariableString:
		var representacion string
		return &representacion, nil
	case VariableEnum:
		var representacion string
		return &representacion, nil
	}

	return nil, fmt.Errorf("no es una variable que pueda tener un valor de una rerpresentacion")
}

func (v Variable) Desreferenciar(referencia any) (any, error) {
	switch informacion := v.Informacion.(type) {
	case VariableSimple:
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

	case VariableString:
		if stringRef, ok := referencia.(*string); !ok {
			return nil, fmt.Errorf("se esperaba que fuera un string, pero no lo es")
		} else {
			return *stringRef, nil
		}
	case VariableEnum:
		if stringRef, ok := referencia.(*string); !ok {
			return nil, fmt.Errorf("se esperaba que fuera un enum, pero no lo es")
		} else {
			return *stringRef, nil
		}
	}

	return nil, fmt.Errorf("no es una variable que pueda tener un valor de una rerpresentacion")
}

type TipoVariableSimple byte

const (
	TVS_INT = iota
	TVS_BOOL
	TVS_DATE
)

type VariableSimple struct {
	Tipo           TipoVariableSimple
	Representativo bool
	Necesario      bool
}

func NewVariableInt(representativo bool, clave string, necesario bool) Variable {
	return Variable{
		Clave: clave,
		Informacion: VariableSimple{
			Tipo:           TVS_INT,
			Representativo: representativo,
			Necesario:      necesario,
		},
	}
}

func NewVariableBool(representativo bool, clave string, necesario bool) Variable {
	return Variable{
		Clave: clave,
		Informacion: VariableSimple{
			Tipo:           TVS_BOOL,
			Representativo: representativo,
			Necesario:      necesario,
		},
	}
}

func NewVariableDate(representativo bool, clave string, necesario bool) Variable {
	return Variable{
		Clave: clave,
		Informacion: VariableSimple{
			Tipo:           TVS_DATE,
			Representativo: representativo,
			Necesario:      necesario,
		},
	}
}

func (vs VariableSimple) TipoSQL() string {
	switch vs.Tipo {
	case TVS_INT:
		return "INT"
	case TVS_BOOL:
		return "BOOLEAN"
	case TVS_DATE:
		return "DATE"
	}
	return "ERROR"
}

type VariableString struct {
	Representativo bool
	Necesario      bool
	Largo          uint
}

func NewVariableString(representativo bool, clave string, largo uint, necesario bool) Variable {
	return Variable{
		Clave: clave,
		Informacion: VariableString{
			Representativo: representativo,
			Necesario:      necesario,
			Largo:          largo,
		},
	}
}

func (vs VariableString) TipoSQL() string {
	return fmt.Sprintf("VARCHAR(%d) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", vs.Largo)
}

type VariableEnum struct {
	Representativo bool
	Necesario      bool
	Valores        []string
}

func NewVariableEnum(representativo bool, clave string, valores []string, necesario bool) Variable {
	return Variable{
		Clave: clave,
		Informacion: VariableEnum{
			Representativo: representativo,
			Necesario:      necesario,
			Valores:        valores,
		},
	}
}

func (ve VariableEnum) TipoSQL() string {
	valoresRep := []string{}
	for _, valor := range ve.Valores {
		valoresRep = append(valoresRep, fmt.Sprintf("\"%s\"", valor))
	}
	return fmt.Sprintf("ENUM(%s)", strings.Join(valoresRep, ", "))
}

type VariableReferencia struct {
	Representativo bool
	Tablas         []*DescripcionTabla
}

func NewVariableReferencia(representativo bool, clave string, tablas []*DescripcionTabla) Variable {
	return Variable{
		Clave: clave,
		Informacion: VariableReferencia{
			Representativo: representativo,
			Tablas:         tablas,
		},
	}
}

type VariableArrayReferencia struct {
	Tablas []*DescripcionTabla
}

func NewVariableArrayReferencias(clave string, tablas []*DescripcionTabla) Variable {
	return Variable{
		Clave: clave,
		Informacion: VariableArrayReferencia{
			Tablas: tablas,
		},
	}
}
