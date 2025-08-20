package dependencias

import (
	"fmt"
	"strings"
)

type informacionVariable any

type Variable struct {
	Clave       string
	Informacion informacionVariable
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
	}

	return parametros
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
