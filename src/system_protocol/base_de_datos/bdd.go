package base_de_datos

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	// _ "github.com/go-sql-driver/mysql"
)

type Bdd struct {
	conn *sql.DB
}

func NewBdd(carpetaOutput, nombreBdd string, canalMensajes chan string) (*Bdd, error) {
	/*
		_ = godotenv.Load()

		dbUser := os.Getenv("MYSQL_USER")
		dbPass := os.Getenv("MYSQL_PASSWORD")
		dbHost := os.Getenv("MYSQL_HOST")
		dbPort := os.Getenv("MYSQL_PORT")
		dbName := os.Getenv("MYSQL_DB_NAME")

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)

		conn, err := sql.Open("mysql", dsn)
		if err != nil {
			return nil, fmt.Errorf("error connecting to DB: %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if err = conn.PingContext(ctx); err != nil {
			return nil, fmt.Errorf("no se pudo pinear el servidor de MySQL, con error: %v", err)
		}

		conn.SetConnMaxIdleTime(5 * time.Minute)
		conn.SetConnMaxLifetime(1 * time.Hour)

	*/
	conn, err := sql.Open("sqlite3", fmt.Sprintf("%s/%s?_journal_mode=WAL", carpetaOutput, nombreBdd))
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

type TipoCheckpoint byte

const (
	TC_FULL = iota
	TC_PASSIVE
	TC_RESET
)

func (bdd *Bdd) Checkpoint(tipo TipoCheckpoint) (err error) {
	switch tipo {
	case TC_FULL:
		_, err = bdd.conn.Exec("PRAGMA wal_checkpoint(full);")
	case TC_PASSIVE:
		_, err = bdd.conn.Exec("PRAGMA wal_checkpoint(passive);")
	case TC_RESET:
		_, err = bdd.conn.Exec("PRAGMA wal_checkpoint(restart);")

	}
	return err
}

func (bdd *Bdd) Close() {
	if err := bdd.Checkpoint(TC_FULL); err != nil {
		fmt.Printf("Salió mal el checkpoint, con error: %v", err)
	}
	bdd.conn.Close()
}

func (bdd *Bdd) CrearTabla(query string, datos ...any) error {
	_, err := bdd.conn.Exec(query, datos...)
	return err
}

func (bdd *Bdd) EliminarTabla(nombreTabla string) error {
	_, err := bdd.conn.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", nombreTabla))
	return err
}

func (bdd *Bdd) Preparar(query string) (Sentencia, error) {
	return NewSentencia(bdd.conn, query)
}

func (bdd *Bdd) Transaccion() (Transaccion, error) {
	return NewTransaccion(bdd.conn)
}
