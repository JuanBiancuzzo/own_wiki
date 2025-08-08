package dependencias

import "hash/maphash"

type IntFK int32

type ForeignKey struct {
	Key              string
	TablaDestino     string
	HashDatosDestino IntFK
}

func NewForeignKey(key string, tabla string, hashDatos IntFK) ForeignKey {
	return ForeignKey{
		Key:              key,
		TablaDestino:     tabla,
		HashDatosDestino: hashDatos,
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
	rep64 := maphash.Bytes(h.Seed, datos)
	rep32 := rep64 >> 32 // Se mantiene con los ultimos 32 bits
	return IntFK(rep32)
}
