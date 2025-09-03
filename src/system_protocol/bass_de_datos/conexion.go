package bass_de_datos

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type conexion struct {
	sql  *sql.DB
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

	bdd.SetMaxOpenConns(1)

	return &conexion{
		sql:  bdd,
		pool: pool,
	}, nil
}

func (c *conexion) Close() {
	c.sql.Close()
}

type resultadoSQL struct {
	LastInsertId func() (int64, error)
	RowsAffected func() (int64, error)
}

func newResultado(resultado sql.Result) resultadoSQL {
	lastInsertedId, errInsert := resultado.LastInsertId()
	rowsAffected, errAffected := resultado.RowsAffected()
	return resultadoSQL{
		LastInsertId: func() (int64, error) { return lastInsertedId, errInsert },
		RowsAffected: func() (int64, error) { return rowsAffected, errAffected },
	}
}

func (c *conexion) Exec(query string, datos ...any) (resultadoSQL, error) {
	resultado, err := c.sql.Exec(fmt.Sprintf("PRAGMA busy_timeout=10000;\n%s;", query), datos...)
	if err != nil {
		c.pool <- c
		return resultadoSQL{}, err
	}
	nuevoResultado := newResultado(resultado)
	c.pool <- c
	return nuevoResultado, nil
}

type filaSQL struct {
	fila *sql.Row

	pool chan *conexion
	conn *conexion
}

func (f filaSQL) Scan(datos ...any) error {
	err := f.fila.Scan(datos...)
	f.pool <- f.conn
	return err
}

func (c *conexion) QueryRow(query string, datos ...any) *filaSQL {
	return &filaSQL{
		fila: c.sql.QueryRow(query, datos...),
		pool: c.pool,
		conn: c,
	}
}

type filasSQL struct {
	filas *sql.Rows

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
	f.pool <- f.conn
}

func (c *conexion) Query(query string, datos ...any) (filasSQL, error) {
	filas, err := c.sql.Query(query, datos...)
	if err != nil {
		c.pool <- c
		return filasSQL{}, err
	}

	return filasSQL{
		filas: filas,
		pool:  c.pool,
		conn:  c,
	}, err
}
