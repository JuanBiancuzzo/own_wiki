package estructura

import (
	"database/sql"
	"fmt"
)

const QUERY_EDITORIAL = "SELECT id FROM editoriales WHERE editorial = ?"
const INSERTAR_EDITORIAL = "INSERT INTO editoriales (editorial) VALUES (?)"

const INSERTAR_LIBRO = "INSERT INTO libros (titulo, subtitulo, idEditorial, anio, edicion, volumen, url, idArchivo) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
const INSERTAR_CAPITULO = "INSERT INTO capitulos (capitulo, nombre, paginaInicial, paginaFinal, idLibro, idArchivo) VALUES (?, ?, ?, ?, ?, ?)"
const INSERTAR_AUTOR_LIBRO = "INSERT INTO autoresLibro (idLibro, idPersona) VALUES (?, ?)"
const INSERTAR_EDITOR_CAPITULO = "INSERT INTO editoresCapitulo (idCapitulo, idPersona) VALUES (?, ?)"

type Capitulo struct {
	Capitulo      int
	Nombre        string
	Editores      []*Persona
	PaginaInicial int
	PaginaFinal   int
}

func NewCapitulo(capitulo string, nombre string, editores []*Persona, paginaInicial string, paginaFinal string) *Capitulo {
	return &Capitulo{
		Capitulo:      NumeroODefault(capitulo, 0),
		Nombre:        nombre,
		Editores:      editores,
		PaginaInicial: NumeroODefault(paginaInicial, 0),
		PaginaFinal:   NumeroODefault(paginaFinal, 0),
	}
}

func (c *Capitulo) Insertar(idLibro int64, idArchivo int64) []any {
	return []any{
		c.Capitulo,
		c.Nombre,
		c.PaginaInicial,
		c.PaginaFinal,
		idLibro,
		idArchivo,
	}
}

type Libro struct {
	PathArchivo string
	Titulo      string
	Subtitulo   string
	Editorial   string
	Anio        int
	Edicion     int
	Volumen     int
	Url         string
	Autores     []*Persona
	Capitulos   []*Capitulo
}

func NewLibro(pathArchivo string, titulo string, subtitulo string, editorial string, anio string, edicion string, volumen string, url string, autores []*Persona, capitulos []*Capitulo) *Libro {
	return &Libro{
		PathArchivo: pathArchivo,
		Titulo:      titulo,
		Subtitulo:   subtitulo,
		Editorial:   editorial,
		Anio:        NumeroODefault(anio, 0),
		Edicion:     NumeroODefault(edicion, 1),
		Volumen:     NumeroODefault(volumen, 0),
		Url:         url,
		Autores:     autores,
		Capitulos:   capitulos,
	}
}

func (l *Libro) Insertar(idEditorial int64, idArchivo int64) []any {
	return []any{
		l.Titulo,
		l.Subtitulo,
		idEditorial,
		l.Anio,
		l.Edicion,
		l.Volumen,
		l.Url,
		idArchivo,
	}
}

func (l *Libro) CargarDatos(bdd *sql.DB, canal chan string) bool {
	canal <- "Insertar Libro"

	var existe bool
	var idArchivo int64
	var idLibro int64

	if idArchivo, existe = Obtener(
		func() *sql.Row { return bdd.QueryRow(QUERY_ARCHIVO, l.PathArchivo) },
	); !existe {
		return false

	} else if idEditorial, err := ObtenerOInsertar(
		func() *sql.Row { return bdd.QueryRow(QUERY_EDITORIAL, l.Editorial) },
		func() (sql.Result, error) { return bdd.Exec(INSERTAR_EDITORIAL, l.Editorial) },
	); err != nil {
		canal <- fmt.Sprintf("error al hacer una querry de la editorial %s con error: %v", l.Editorial, err)

	} else if idLibro, err = Insertar(func() (sql.Result, error) {
		return bdd.Exec(INSERTAR_LIBRO, l.Insertar(idEditorial, idArchivo)...)
	}); err != nil {
		canal <- fmt.Sprintf("error al insertar un libro, con error: %v", err)

	}

	for _, autor := range l.Autores {
		if idAutor, err := ObtenerOInsertar(
			func() *sql.Row { return bdd.QueryRow(QUERY_PERSONAS, autor.Insertar()...) },
			func() (sql.Result, error) { return bdd.Exec(INSERTAR_PERSONA, autor.Insertar()...) },
		); err != nil {
			canal <- fmt.Sprintf("error al hacer una querry del autor: %s %s con error: %v", autor.Nombre, autor.Apellido, err)

		} else if _, err := bdd.Exec(INSERTAR_AUTOR_LIBRO, idLibro, idAutor); err != nil {
			canal <- fmt.Sprintf("error al insertar par idLibro-idAutor, con error: %v", err)
		}
	}

	for _, capitulo := range l.Capitulos {
		if idCapitulo, err := Insertar(func() (sql.Result, error) {
			return bdd.Exec(INSERTAR_CAPITULO, capitulo.Insertar(idLibro, idArchivo)...)
		}); err != nil {
			canal <- fmt.Sprintf("error al insertar un capitulo, con error: %v", err)
		} else {
			for _, autor := range capitulo.Editores {
				if idAutor, err := ObtenerOInsertar(
					func() *sql.Row { return bdd.QueryRow(QUERY_PERSONAS, autor.Insertar()...) },
					func() (sql.Result, error) { return bdd.Exec(INSERTAR_PERSONA, autor.Insertar()...) },
				); err != nil {
					canal <- fmt.Sprintf("error al hacer una querry del autor: %s %s con error: %v", autor.Nombre, autor.Apellido, err)

				} else if _, err := bdd.Exec(INSERTAR_EDITOR_CAPITULO, idCapitulo, idAutor); err != nil {
					canal <- fmt.Sprintf("error al insertar par idLibro-idAutor, con error: %v", err)
				}
			}
		}
	}

	return true
}
