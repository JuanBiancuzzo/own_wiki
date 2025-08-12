package bass_de_datos

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"
)

/*
Que tablas necesitamos
  - Archivos - DONE
  - Tags - DONE
  - Libros (con sus capitulos) - DONE
  - Distribuciones - DONE
  - Papers - DONE
  - Carreras - DONE
  - Materias - DONE
  - Resumenes materias - DONE
  - Cursos - DONE
  - Resumenes cursos - DONE
  - Referencias
  - Data structures
  - Documentos
  - Teoremas, procposiciones y observaciones
  - Temas de investigacion
  - Impresion 3d (todavia ni lo tengo definido entonces tal vez despues)
  - Librerias (todavia no completamente definido entonces tal vez despues)
  - Programas (todavia no completamente definido entonces tal vez despues)
  - Recetas (todavia no completamente definido entonces tal vez despues)
*/

func EstablecerConexionRelacional(canalMensajes chan string) (*sql.DB, error) {
	dbUser := os.Getenv("MYSQL_USER")
	dbPass := os.Getenv("MYSQL_PASSWORD")
	dbHost := os.Getenv("MYSQL_HOST")
	dbPort := os.Getenv("MYSQL_PORT")
	dbName := os.Getenv("MYSQL_DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)

	bdd, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to DB: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err = bdd.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("no se pudo pinear el servidor de MySQL, con error: %v", err)
	}
	canalMensajes <- "Se conecto correctamente a MySQL"

	return bdd, nil
}

func CerrarBddRelacional(bdd *sql.DB) {
	if bdd == nil {
		return
	}

	bdd.Close()
}
