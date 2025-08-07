package dependencias

import "hash/maphash"

type IntFK uint64

type ForeignKey struct {
	Tabla     string
	HashDatos IntFK
}

func NewForeignKey(tabla string, hashDatos IntFK) ForeignKey {
	return ForeignKey{
		Tabla:     tabla,
		HashDatos: hashDatos,
	}
}

type Hash struct {
	Seed maphash.Seed
}

func NewHash() *Hash {
	return &Hash{
		Seed: maphash.MakeSeed(),
	}
}

func (h *Hash) HasearDatos(datos []byte) IntFK {
	return IntFK(maphash.Bytes(h.Seed, datos))
}
