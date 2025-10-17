package dependencias

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

func (dt DescripcionTabla) ObtenerVariable(clave string) (DescripcionVariable, bool) {
	if clave == "id" {
		return NewDescVariableSimple(TVS_INT, true, clave, true), true
	}

	for _, variable := range dt.Variables {
		if variable.Clave == clave {
			return variable, true
		}
	}

	return DescripcionVariable{}, false
}
