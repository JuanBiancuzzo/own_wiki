package dependencias

import b "own_wiki/system_protocol/bass_de_datos"

type TipoTabla byte

const (
	DEPENDIENTE_NO_DEPENDIBLE   = 0b00
	DEPENDIENTE_DEPENDIBLE      = 0b10
	INDEPENDIENTE_NO_DEPENDIBLE = 0b01
	INDEPENDIENTE_DEPENDIBLE    = 0b11
)

func EsTipoDependiente(tipo TipoTabla) bool {
	return tipo&0b01 != 0b01
}

func EsTipoDependible(tipo TipoTabla) bool {
	return tipo&0b10 == 0b10
}

type Tabla interface {
	// Nombre de la tabla
	Nombre() string
	// Query para insertar los datos
	Query(bdd *b.Bdd, datos ...any) (int64, error)

	// Es decir, que no aparezcan NOT NULL
	CrearTablaRelajada(bdd *b.Bdd, info *b.InfoArchivos) error
	RestringirTabla(bdd *b.Bdd) error

	ObtenerDependencias() []Tabla
}
