package bass_de_datos

import (
	"database/sql"
	"fmt"
	"sync"
)

type Bdd struct {
	sql  *sql.DB
	lock *sync.RWMutex
}

func NewBdd(sql *sql.DB) *Bdd {
	var lock sync.RWMutex
	return &Bdd{
		sql:  sql,
		lock: &lock,
	}
}

func (bdd *Bdd) CrearTabla(query string, datos ...any) error {
	_, err := bdd.sql.Exec(query, datos...)
	return err
}

func (bdd *Bdd) EliminarTabla(nombreTabla string) error {
	_, err := bdd.sql.Exec(fmt.Sprintf("DROP TABLE %s", nombreTabla))
	return err
}

func (bdd *Bdd) Existe(query string, datos ...any) (bool, error) {
	lectura := make([]any, len(datos))
	bdd.lock.RLock()
	fila := bdd.sql.QueryRow(query, datos...)
	bdd.lock.RUnlock()

	if err := fila.Scan(lectura...); err != nil {
		return false, nil
	}

	return true, nil
}

func (bdd *Bdd) Obtener(query string, datos ...any) (int64, error) {
	bdd.lock.RLock()
	var id int64
	fila := bdd.sql.QueryRow(query, datos...)
	bdd.lock.RUnlock()

	if err := fila.Scan(&id); err != nil {
		return id, fmt.Errorf("error al intentar query la bdd, con error: %v", err)
	}

	return id, nil
}

func (bdd *Bdd) Insertar(query string, datos ...any) (int64, error) {
	if filaAfectada, err := bdd.Exec(query, datos...); err != nil {
		return 0, fmt.Errorf("error al insertar con query, con error: %v", err)

	} else if id, err := filaAfectada.LastInsertId(); err != nil {
		return 0, fmt.Errorf("error al obtener id from query, con error: %v", err)

	} else {
		return id, nil
	}
}

func (bdd *Bdd) ObtenerOInsertar(queryObtener, queryInsertar string, datos ...any) (int64, error) {
	if id, err := bdd.Obtener(queryObtener, datos...); err == nil {
		return id, nil
	}
	return bdd.Insertar(queryInsertar, datos...)
}

func (bdd *Bdd) Exec(query string, datos ...any) (sql.Result, error) {
	bdd.lock.Lock()
	filaAfectada, err := bdd.sql.Exec(query, datos...)
	bdd.lock.Unlock()

	return filaAfectada, err
}

func (bdd *Bdd) QueryRow(query string, datos ...any) *sql.Row {
	bdd.lock.RLock()
	fila := bdd.sql.QueryRow(query, datos...)
	bdd.lock.RUnlock()

	return fila
}

type filasSQL struct {
	filas *sql.Rows
	lock  *sync.RWMutex
}

func (f filasSQL) Next() bool {
	return f.filas.Next()
}

func (f filasSQL) Scan(datos ...any) error {
	return f.filas.Scan(datos...)
}

func (f filasSQL) Close() {
	f.filas.Close()
	f.lock.RUnlock()
}

func (bdd *Bdd) Query(query string, datos ...any) (filasSQL, error) {
	bdd.lock.RLock()
	filas, err := bdd.sql.Query(query, datos...)

	return filasSQL{
		filas: filas,
		lock:  bdd.lock,
	}, err
}
