package dependencias

type DescripcionVariable struct {
	Clave       string
	Descripcion any
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
