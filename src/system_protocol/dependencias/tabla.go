package dependencias

import (
	"fmt"
	b "own_wiki/system_protocol/bass_de_datos"
	"strings"
)

type TipoTabla byte

const (
	DEPENDIENTE_NO_DEPENDIBLE   = 0b00
	DEPENDIENTE_DEPENDIBLE      = 0b10
	INDEPENDIENTE_NO_DEPENDIBLE = 0b01
	INDEPENDIENTE_DEPENDIBLE    = 0b11
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

type ReferenciaTabla struct {
	Clave string
	Tabla DescripcionTabla
}

func NewReferenciaTabla(clave string, tabla DescripcionTabla) ReferenciaTabla {
	return ReferenciaTabla{
		Clave: clave,
		Tabla: tabla,
	}
}

type DescripcionTabla struct {
	NombreTabla string
	TipoTabla   TipoTabla
	ClavesTipo  []ParClaveTipo
	Referencias []ReferenciaTabla

	necesarioQuery bool
	query          string
	insertar       string
}

func ConstruirTabla(nombreTabla string, tipoTabla TipoTabla, elementosUnicos bool, clavesTipo []ParClaveTipo, referencias []ReferenciaTabla) DescripcionTabla {
	insertarParam := []string{}
	insertarValues := []string{}
	queryParam := []string{}
	for _, claveTipo := range clavesTipo {
		insertarParam = append(insertarParam, claveTipo.Clave)
		insertarValues = append(insertarValues, "?")
		queryParam = append(queryParam, fmt.Sprintf("%s = ?", claveTipo.Clave))
	}
	for _, referencia := range referencias {
		insertarParam = append(insertarParam, referencia.Clave)
		insertarValues = append(insertarValues, "0")
	}

	return DescripcionTabla{
		NombreTabla: nombreTabla,
		TipoTabla:   tipoTabla,
		ClavesTipo:  clavesTipo,
		Referencias: referencias,

		necesarioQuery: elementosUnicos,
		insertar: fmt.Sprintf(
			"INSERT INTO %s (%s) VALUES (%s)",
			nombreTabla,
			strings.Join(insertarParam, ", "),
			strings.Join(insertarValues, ", "),
		),
		query: fmt.Sprintf(
			"SELECT id FROM %s WHERE %s",
			nombreTabla,
			strings.Join(queryParam, " AND "),
		),
	}
}

func (dt DescripcionTabla) CrearTablaRelajada(bdd *b.Bdd) error {
	parametros := []string{}
	for i, parClaveTipo := range dt.ClavesTipo {
		extra := ","
		if i+1 == len(dt.ClavesTipo) && len(dt.Referencias) == 0 {
			extra = ""
		}
		parametros = append(parametros, fmt.Sprintf("%s %s%s", parClaveTipo.Clave, parClaveTipo.Tipo, extra))
	}

	for i, referencia := range dt.Referencias {
		extra := ","
		if i+1 == len(dt.Referencias) {
			extra = ""
		}
		parametros = append(parametros, fmt.Sprintf("%s INT%s", referencia.Clave, extra))
	}

	tabla := fmt.Sprintf(
		"CREATE TABLE %s (\nid INT AUTO_INCREMENT PRIMARY KEY,\n\t%s\n);",
		dt.NombreTabla,
		strings.Join(parametros, "\n\t"),
	)

	if err := bdd.CrearTabla(tabla); err != nil {
		return fmt.Errorf("no se pudo crear la tabla \n%s\n, con error: %v", tabla, err)
	}
	return nil
}

// TODO
func (dt DescripcionTabla) RestringirTabla(bdd *b.Bdd) error {
	return nil
}

func (dt DescripcionTabla) Hash(hash *Hash, datos ...any) (IntFK, error) {
	if len(datos) != len(dt.ClavesTipo) {
		return 0, fmt.Errorf("en la tabla %s, al hashear %T, no tenia la misma estructura que la esperada", dt.NombreTabla, datos)
	}

	datosRepresentativos := []any{}
	for i, claveTipo := range dt.ClavesTipo {
		if claveTipo.Representativa {
			datosRepresentativos = append(datosRepresentativos, datos[i])
		}
	}

	return hash.HasearDatos(datosRepresentativos...), nil
}

func (dt DescripcionTabla) Existe(bdd *b.Bdd, datos ...any) (bool, error) {
	if !dt.necesarioQuery {
		return false, nil
	}

	_, err := bdd.Obtener(dt.query, datos...)
	return err == nil, nil
}

func (dt DescripcionTabla) ObtenerDependencias() []DescripcionTabla {
	tablas := []DescripcionTabla{}

	for _, referencia := range dt.Referencias {
		tablas = append(tablas, referencia.Tabla)
	}

	return tablas
}
func EsTipoDependiente(tipo TipoTabla) bool {
	return tipo == DEPENDIENTE_DEPENDIBLE || tipo == DEPENDIENTE_NO_DEPENDIBLE
}

func EsTipoDependible(tipo TipoTabla) bool {
	return tipo == DEPENDIENTE_DEPENDIBLE || tipo == INDEPENDIENTE_DEPENDIBLE
}
