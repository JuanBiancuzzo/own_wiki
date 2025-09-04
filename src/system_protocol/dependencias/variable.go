package dependencias

import (
	"fmt"
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
		parametros = []string{fmt.Sprintf("%s %s", v.Clave, informacion.TipoSQL(v.Clave))}

	case VariableReferencia:
		cantidad := min(2, len(informacion.Tablas))
		parametros = make([]string, cantidad)
		if len(informacion.Tablas) > 1 {
			maximoLargo := 0
			for _, tabla := range informacion.Tablas {
				maximoLargo = max(maximoLargo, len(tabla.NombreTabla))
			}
			nombreVariable := fmt.Sprintf("tipo%s", v.Clave)
			parametros[0] = fmt.Sprintf("%s TEXT CHECK( LENGTH(%s) <= %d )", nombreVariable, nombreVariable, maximoLargo)
		}

		parametros[cantidad-1] = fmt.Sprintf("%s INT", v.Clave)

	case VariableArrayReferencia:
		// no es necesario, ya que no tiene representacion en sql
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

func NewVariableSimple(tipo TipoVariableSimple, representativo bool, clave string, necesario bool) Variable {
	return Variable{
		Clave: clave,
		Informacion: VariableSimple{
			Tipo:           tipo,
			Representativo: representativo,
			Necesario:      necesario,
		},
	}
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
	return fmt.Sprintf("VARCHAR(%d)", vs.Largo)
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

func (ve VariableEnum) TipoSQL(clave string) string {
	maximoLargo := 0
	for _, valor := range ve.Valores {
		maximoLargo = max(maximoLargo, len(valor))
	}
	return fmt.Sprintf("TEXT CHECK( LENGTH(%s) <= %d )", clave, maximoLargo)
}

type VariableReferencia struct {
	Representativo bool
	Tablas         []*Tabla
}

func NewVariableReferencia(representativo bool, clave string, tablas []*Tabla) Variable {
	return Variable{
		Clave: clave,
		Informacion: VariableReferencia{
			Representativo: representativo,
			Tablas:         tablas,
		},
	}
}

type VariableArrayReferencia struct {
	ClaveSelf   string
	TablaCreada string
}

func NewVariableArrayReferencias(clave, claveSelf, tablaCreada string) Variable {
	return Variable{
		Clave: clave,
		Informacion: VariableArrayReferencia{
			ClaveSelf:   claveSelf,
			TablaCreada: tablaCreada,
		},
	}
}
