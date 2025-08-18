package dependencias

import (
	"fmt"
	"strconv"
	"strings"
)

type TipoVariable byte

func (tv TipoVariable) ValorDeString(valor string) (any, error) {
	switch tv {
	case TV_INT:
		fallthrough
	case TV_REFERENCIA:
		return strconv.Atoi(valor)

	case TV_BOOL:
		switch strings.ToLower(valor) {
		case "true":
			return true, nil
		case "false":
			return false, nil
		default:
			return nil, fmt.Errorf("no es un valor bool, ya que es %s", valor)
		}

	case TV_STRING:
		fallthrough
	case TV_ENUM:
		fallthrough
	case TV_DATE:
		return valor, nil
	}

	return nil, fmt.Errorf("no tiene el tipo correcto")
}

func (tv TipoVariable) ReferenciaValor() (any, error) {
	switch tv {
	case TV_INT:
		fallthrough
	case TV_REFERENCIA:
		var numero int
		return &numero, nil

	case TV_BOOL:
		var booleano bool
		return &booleano, nil

	case TV_STRING:
		fallthrough
	case TV_ENUM:
		fallthrough
	case TV_DATE:
		var texto string
		return &texto, nil

	}

	return nil, fmt.Errorf("no tiene el tipo correcto, tiene: %v", tv)
}

const (
	TV_INT = iota
	TV_BOOL
	TV_STRING
	TV_ENUM
	TV_DATE
	TV_REFERENCIA
)

type ParClaveTipo struct {
	Representativa bool
	Clave          string
	TipoSQL        string
	Necesario      bool

	tipo TipoVariable
}

func NewClaveInt(representativo bool, clave string, necesario bool) ParClaveTipo {
	return ParClaveTipo{
		Representativa: representativo,
		Clave:          clave,
		TipoSQL:        "INT",
		Necesario:      necesario,
		tipo:           TV_INT,
	}
}

func NewClaveBool(representativo bool, clave string, necesario bool) ParClaveTipo {
	return ParClaveTipo{
		Representativa: representativo,
		Clave:          clave,
		TipoSQL:        "BOOLEAN",
		Necesario:      necesario,
		tipo:           TV_BOOL,
	}
}

func NewClaveString(representativo bool, clave string, largo uint, necesario bool) ParClaveTipo {
	tipo := fmt.Sprintf("VARCHAR(%d) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", largo)
	return ParClaveTipo{
		Representativa: representativo,
		Clave:          clave,
		TipoSQL:        tipo,
		Necesario:      necesario,
		tipo:           TV_STRING,
	}
}

func NewClaveEnum(representativo bool, clave string, valores []string, necesario bool) ParClaveTipo {
	valoresRep := []string{}
	for _, valor := range valores {
		valoresRep = append(valoresRep, fmt.Sprintf("\"%s\"", valor))
	}

	return ParClaveTipo{
		Representativa: representativo,
		Clave:          clave,
		TipoSQL:        fmt.Sprintf("ENUM(%s)", strings.Join(valoresRep, ", ")),
		Necesario:      necesario,
		tipo:           TV_ENUM,
	}
}

func NewClaveDate(representativo bool, clave string, necesario bool) ParClaveTipo {
	return ParClaveTipo{
		Representativa: representativo,
		Clave:          clave,
		TipoSQL:        "DATE",
		Necesario:      necesario,
		tipo:           TV_DATE,
	}
}
