package bass_de_datos

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type conexion struct {
	sql  *sql.DB
	lock *sync.RWMutex
	pool chan *conexion
}

func newConexion(archivoBdd string, pool chan *conexion) (*conexion, error) {
	bdd, err := sql.Open("sqlite3", archivoBdd)
	if err != nil {
		return nil, fmt.Errorf("error connecting to DB: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err = bdd.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("no se pudo pinear el servidor de SQLite, con error: %v", err)
	}

	var lock sync.RWMutex
	return &conexion{
		sql:  bdd,
		lock: &lock,
		pool: pool,
	}, nil
}

func (c *conexion) Devolver() {
	if c.pool == nil {
		return
	}

	c.pool <- c
}

func (c *conexion) Close() {
	c.sql.Close()
}

func (c *conexion) Exec(query string, datos ...any) (sql.Result, error) {
	c.lock.Lock()
	filaAfectada, err := c.sql.Exec(query, datos...)
	c.lock.Unlock()

	return filaAfectada, err
}

func (c *conexion) QueryRow(query string, datos ...any) *sql.Row {
	c.lock.RLock()
	fila := c.sql.QueryRow(query, datos...)
	c.lock.RUnlock()

	return fila
}

type filasSQL struct {
	filas *sql.Rows
	lock  *sync.RWMutex

	pool chan *conexion
	conn *conexion
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
	f.pool <- f.conn
}

func (c *conexion) Query(query string, datos ...any) (filasSQL, error) {
	c.lock.RLock()
	filas, err := c.sql.Query(query, datos...)

	return filasSQL{
		filas: filas,
		lock:  c.lock,
		pool:  c.pool,
		conn:  c,
	}, err
}
