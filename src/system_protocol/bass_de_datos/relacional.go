package bass_de_datos

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func EstablecerConexionRelacional(carpetaOutput string, canalMensajes chan string) (*sql.DB, error) {
	bdd, err := sql.Open("sqlite3", fmt.Sprintf("%s/baseDeDatos.db", carpetaOutput))
	if err != nil {
		return nil, fmt.Errorf("error connecting to DB: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err = bdd.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("no se pudo pinear el servidor de SQLite, con error: %v", err)
	}
	canalMensajes <- "Se conecto correctamente a SQLite"

	return bdd, nil
}

func CerrarBddRelacional(bdd *sql.DB) {
	if bdd == nil {
		return
	}

	bdd.Close()
}
