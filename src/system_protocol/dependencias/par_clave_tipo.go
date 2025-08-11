package dependencias

import (
	"fmt"
	"strings"
)

type ParClaveTipo struct {
	Representativa bool
	Clave          string
	Tipo           string
	Necesario      bool
}

func NewClaveInt(representativo bool, clave string, necesario bool) ParClaveTipo {
	return ParClaveTipo{
		Representativa: representativo,
		Clave:          clave,
		Tipo:           "INT",
		Necesario:      necesario,
	}
}

func NewClaveBool(representativo bool, clave string, necesario bool) ParClaveTipo {
	return ParClaveTipo{
		Representativa: representativo,
		Clave:          clave,
		Tipo:           "BOOLEAN",
		Necesario:      necesario,
	}
}

func NewClaveString(representativo bool, clave string, largo uint, necesario bool) ParClaveTipo {
	tipo := fmt.Sprintf("VARCHAR(%d) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", largo)
	return ParClaveTipo{
		Representativa: representativo,
		Clave:          clave,
		Tipo:           tipo,
		Necesario:      necesario,
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
		Tipo:           fmt.Sprintf("ENUM(%s)", strings.Join(valoresRep, ", ")),
		Necesario:      necesario,
	}
}
