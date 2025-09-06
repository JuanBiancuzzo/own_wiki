package base_de_datos

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const NOMBRE_BDD = "baseDeDatos.db"

type Bdd struct {
	conn *sql.DB
}

func NewBdd(carpetaOutput string, canalMensajes chan string) (*Bdd, error) {
	conn, err := sql.Open("sqlite3", fmt.Sprintf("%s/%s?_journal_mode=WAL", carpetaOutput, NOMBRE_BDD))
	if err != nil {
		return nil, fmt.Errorf("error connecting to DB: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err = conn.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("no se pudo pinear el servidor de SQLite, con error: %v", err)
	}

	return &Bdd{
		conn: conn,
	}, nil
}

func (bdd *Bdd) Checkpoint() error {
	_, err := bdd.conn.Exec("PRAGMA wal_checkpoint(TRUNCATE);")
	return err
}

func (bdd *Bdd) Close() {
	if err := bdd.Checkpoint(); err != nil {
		fmt.Printf("Salió mal el checkpoint, con error: %v", err)
	}
	bdd.conn.Close()
}

func (bdd *Bdd) CrearTabla(query string, datos ...any) error {
	_, err := bdd.exec(query, datos...)
	return err
}

func (bdd *Bdd) EliminarTabla(nombreTabla string) error {
	_, err := bdd.exec(fmt.Sprintf("DROP TABLE %s", nombreTabla))
	return err
}

func (bdd *Bdd) Existe(query string, datos ...any) (bool, error) {
	lectura := make([]any, len(datos))
	fila := bdd.QueryRow(query, datos...)

	if err := fila.Scan(lectura...); err != nil {
		return false, nil
	}

	return true, nil
}

func (bdd *Bdd) Obtener(query string, datos ...any) (int64, error) {
	var id int64
	fila := bdd.QueryRow(query, datos...)
	if fila == nil {
		return id, fmt.Errorf("error al obtener query")
	}

	if err := fila.Scan(&id); err != nil {
		return id, fmt.Errorf("error al intentar query la bdd, con error: %v", err)
	}

	return id, nil
}

func (bdd *Bdd) InsertarId(query string, datos ...any) (int64, error) {
	if filaAfectada, err := bdd.exec(query, datos...); err != nil {
		return 0, fmt.Errorf("error al insertar con query (ejecutando exec), con error: %v", err)

	} else if id, err := filaAfectada.LastInsertId(); err != nil {
		return 0, fmt.Errorf("error al obtener id from query, con error: %v", err)

	} else {
		return id, nil
	}
}

func (bdd *Bdd) Update(query string, datos ...any) error {
	_, err := bdd.exec(query, datos...)
	return err
}

func (bdd *Bdd) Eliminar(query string, datos ...any) error {
	_, err := bdd.exec(query, datos...)
	return err
}

func (bdd *Bdd) QueryRow(query string, datos ...any) *sql.Row {
	return bdd.conn.QueryRow(query, datos...)
}

func (bdd *Bdd) Query(query string, datos ...any) (*sql.Rows, error) {
	return bdd.conn.Query(query, datos...)
}

func (bdd *Bdd) exec(query string, datos ...any) (sql.Result, error) {
	return bdd.conn.Exec(query, datos...)
}
