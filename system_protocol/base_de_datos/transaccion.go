package base_de_datos

import (
	"database/sql"
)

type Transaccion struct {
	transaccion *sql.Tx
}

func NewTransaccion(bdd *sql.DB) (Transaccion, error) {
	if tx, err := bdd.Begin(); err != nil {
		return Transaccion{}, err
	} else {
		return Transaccion{
			transaccion: tx,
		}, nil
	}
}

func (t Transaccion) Commit() error {
	return t.transaccion.Commit()
}

func (t Transaccion) RollBack() error {
	return t.transaccion.Rollback()
}

func (t Transaccion) Sentencia(sentencia Sentencia) Sentencia {
	return sentencia.NewSentenciaDeTransaccion(t.transaccion)
}
