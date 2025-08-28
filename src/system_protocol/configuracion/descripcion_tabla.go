package configuracion

import (
	d "own_wiki/system_protocol/dependencias"
)

type DescripcionVariable struct {
	Clave       string
	Descripcion any
}

type DescVariableSimple struct {
	Tipo           d.TipoVariableSimple
	Representativo bool
	Necesario      bool
}

func NewVariableSimple(tipo d.TipoVariableSimple, representativo bool, clave string, necesario bool) DescripcionVariable {
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

func NewVariableString(representativo bool, clave string, largo uint, necesario bool) DescripcionVariable {
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

func NewVariableEnum(representativo bool, clave string, valores []string, necesario bool) DescripcionVariable {
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

func NewVariableReferencia(representativo bool, clave string, tablas []string) DescripcionVariable {
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

func NewVariableArrayReferencias(clave, claveSelf, tablaCreada string) DescripcionVariable {
	return DescripcionVariable{
		Clave: clave,
		Descripcion: DescVariableArrayReferencia{
			ClaveSelf:   claveSelf,
			TablaCreada: tablaCreada,
		},
	}
}

type DescripcionTabla struct {
	Nombre             string
	ElementosRepetidos bool
	Variables          []DescripcionVariable
}

func NewDescripcionTabla(nombreTabla string, elementosRepetidos bool, variables []DescripcionVariable) DescripcionTabla {
	return DescripcionTabla{
		Nombre:             nombreTabla,
		ElementosRepetidos: elementosRepetidos,
		Variables:          variables,
	}
}
