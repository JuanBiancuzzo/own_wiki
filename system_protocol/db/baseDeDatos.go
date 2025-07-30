package db

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"strings"

	l "own_wiki/system_protocol/listas"

	"github.com/joho/godotenv"
)

//go:embed esquema.sql
var crearTablas string

func EstablecerBaseDeDatos() (*sql.DB, error) {
	_ = godotenv.Load()

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to DB: %v", err)
	}

	for {
		if err = db.Ping(); err == nil {
			fmt.Println("Successfully connected to the database!")
			break
		}
	}

	return db, nil
}

func datosParaTabla(info *InfoArchivos) []any {
	// todos son largos de algo
	return []any{
		info.MaxPath,   // path
		info.MaxNombre, // nombre
		info.MaxNombre, // apellido
		255,            // nombre editorial
		info.MaxTags,   // tags
		255,            // titulo del libro
		255,            // subtitulo del libro
		info.MaxUrl,    // url del libro
		255,            // nombre del capitulo
		255,            // nombre de una distribucion de proba
		255,            // nombre de una carrera
		255,            // nombre del plan de una carrera
		255,            // nombre de una materia
		255,            // Codigo de materia
		255,            // nombre de una materia equivalente
		255,            // Codigo de materia
		255,            // nombre del tema de una materia
		255,            // nombre de la pagina del curso
		255,            // nombre del curso
		info.MaxUrl,    // url del curso
		255,            // nombre del tema del curso
		255,            // nombre del tema de investigacion
		255,            // nombre de revista de papers
		255,            // titulo del paper
		255,            // subtitulo del paper
		info.MaxUrl,    // url paper
		255,            // nombre del tema de matematica
		255,            // nombre del bloque de matematica
		255,            // nombre del grupo legal
		255,            // nombre de la seccion legal
		255,            // abreviacion del documento legal
		255,            // nombre del articulo
		255,            // nombre del canal de youtube
		255,            // nombre del video de youtube
		info.MaxUrl,    // url del video de youtube
		255,            // nombre de la pagina web
		255,            // titulo de la pagina web
		info.MaxUrl,    // url de la pagina web
		255,            // nombre articulo wiki
		info.MaxUrl,    // url articulo wiki
		255,            // nombre del diccionario
		255,            // palabra del diccionario
		info.MaxUrl,    // url del diccionario
	}
}

/*
Que tablas necesitamos
  - Archivos - DONE
  - Tags - DONE
  - Libros (con sus capitulos)
  - Referencias
  - Data structures
  - Distribuciones
  - Documentos
  - Impresion 3d (todavia ni lo tengo definido entonces tal vez despues)
  - Librerias (todavia no completamente definido entonces tal vez despues)
  - Papers
  - Programas (todavia no completamente definido entonces tal vez despues)
  - Recetas (todavia no completamente definido entonces tal vez despues)
  - Teoremas, procposiciones y observaciones
  - Carreras
  - Materias
  - Resumenes materias
  - Cursos
  - Resumenes cursos
  - Temas de investigacion
*/
func CrearTablas(db *sql.DB, info *InfoArchivos) error {
	queryCrearTablas := fmt.Sprintf(strings.ReplaceAll(crearTablas, "(?)", "(%d)"), datosParaTabla(info)...)

	cantidadTablas := uint32(len(strings.Split(queryCrearTablas, ");")))

	pilaEliminar := l.NewPilaConCapacidad[string](cantidadTablas)
	colaQuery := l.NewColaConCapacidad[string](cantidadTablas)

	for query := range strings.SplitSeq(queryCrearTablas, ";") {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}
		colaQuery.Encolar(query)

		var _extra string
		nombreTabla := "unknown"
		if _, err := fmt.Sscanf(query, "CREATE TABLE IF NOT EXISTS %s %s\n", &nombreTabla, &_extra); err == nil {
			pilaEliminar.Apilar(nombreTabla)
		}
	}

	// Tal vez en vez de eliminarla se pueden alterar, pero eso ya implica un poco mas de esfuerzo

	fmt.Println("Limpiando sus datos")
	for !pilaEliminar.Vacia() {
		if tabla, err := pilaEliminar.Desapilar(); err != nil {
			fmt.Printf("error al intentar eliminar tabla con error: %v\n", err)

		} else if _, err = db.Exec(fmt.Sprintf("DROP TABLE %s;", tabla)); err != nil {
			fmt.Printf("error al intentar eliminar la tabla %s, con error: %v\n", tabla, err)

		} else {
			fmt.Printf("Eliminar tabla: %s\n", tabla)
		}
	}

	fmt.Printf("Creando todas las tablas\n")
	for !colaQuery.Vacia() {
		if query, err := colaQuery.Desencolar(); err != nil {
			fmt.Printf("error al intentar crear tabla con error: %v\n", err)

		} else if _, err = db.Exec(query); err != nil {
			fmt.Printf("error al intentar eliminar la tabla con la query %s, con error: %v\n", query, err)

		} else {
			fmt.Printf("Crear la tabla con la query: %s\n", query)
		}
	}

	fmt.Println("Tablas creadas")

	return nil
}
