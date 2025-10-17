package dependencias

import (
	"encoding/binary"
	"hash/maphash"
)

type IntFK int32

type ForeignKey struct {
	Clave            string
	TablaDestino     string
	HashDatosDestino IntFK
}

func NewForeignKey(tabla, clave string, hash IntFK) ForeignKey {
	return ForeignKey{
		Clave:            clave,
		TablaDestino:     tabla,
		HashDatosDestino: hash,
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

func (h *Hash) HasearDatos(datos ...any) IntFK {
	bufInt := make([]byte, 4)
	datosBytes := []byte{}

	for _, dato := range datos {
		switch valor := dato.(type) {
		case bool:
			var numero uint32 = 0
			if valor {
				numero = 1
			}
			binary.BigEndian.PutUint32(bufInt, numero)
			datosBytes = append(datosBytes, bufInt...)

		case IntFK:
			binary.BigEndian.PutUint32(bufInt, uint32(valor))
			datosBytes = append(datosBytes, bufInt...)

		case int:
			binary.BigEndian.PutUint32(bufInt, uint32(valor))
			datosBytes = append(datosBytes, bufInt...)

		case string:
			datosBytes = append(datosBytes, []byte(valor)...)
		}
	}

	rep64 := maphash.Bytes(h.Seed, datosBytes)
	rep32 := rep64 >> 32 // Se mantiene con los ultimos 32 bits
	return IntFK(rep32)
}
