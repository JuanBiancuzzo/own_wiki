package dependencias

type IntFK uint64

type ForeignKey struct {
	Tabla     string
	HashDatos IntFK
}
