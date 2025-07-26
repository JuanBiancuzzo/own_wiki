package fs

import (
	"database/sql"
	"fmt"
	"os"
	"own_wiki/system_protocol/db"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
)

var insertar_archivo = "INSERT INTO archivos (path) VALUES (?)"

var insertar_tag = "INSERT INTO tags (tag, idArchivo) VALUES (?, ?)"

var query_editorial = "SELECT id FROM editoriales WHERE editorial = ?"
var insertar_editorial = "INSERT INTO editoriales (editorial) VALUES (?)"

var query_personas = "SELECT id FROM personas WHERE nombre = ? AND apellido = ?"
var insertar_persona = "INSERT INTO personas (nombre, apellido) VALUES (?, ?)"

var insertar_libro = "INSERT INTO libros (titulo, subtitulo, anio, idEditorial, edicion, volumen, url, idArchivo) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"

var insertar_capitulo = "INSERT INTO capitulos (capitulo, nombre, paginaInicial, paginaFinal, idLibro, idArchivo) VALUES (?, ?, ?, ?, ?, ?)"

var insertar_autor_libro = "INSERT INTO autoresLibro (idLibro, idPersona) VALUES (?, ?)"
var insertar_editor_capitulo = "INSERT INTO editoresCapitulo (idCapitulo, idPersona) VALUES (?, ?)"

type Archivo struct {
	Path      string
	Metadata  *Frontmatter
	Contenido string
	IdArchivo int64
}

func NewArchivo(path string) *Archivo {
	return &Archivo{
		Path:      path,
		Metadata:  nil,
		Contenido: "",
		IdArchivo: 0,
	}
}

func (a *Archivo) Interprestarse(infoArchivos *db.InfoArchivos, canal chan string) {
	if !strings.Contains(a.Path, ".md") {
		return
	}

	bytes, err := os.ReadFile(a.Path)
	if err != nil {
		canal <- fmt.Sprintf("Error al leer %s obteniendo el error: %v", a.Path, err)
		return
	}
	contenido := string(bytes)

	if strings.Index(contenido, "---") != 0 {
		return
	}

	blob := contenido[3 : 3+strings.Index(contenido[3:], "---")]
	decodificador := yaml.NewDecoder(strings.NewReader(blob))

	var metadata Frontmatter
	err = decodificador.Decode(&metadata)
	if err != nil {
		canal <- fmt.Sprintf("Error al decodificar en %s la metadata, con el error: %v\n", a.Path, err)
		return
	}

	a.Metadata = &metadata
	a.Contenido = contenido[3+strings.Index(contenido[3:], "---")+len("---"):]

	for _, tag := range a.Metadata.Tags {
		infoArchivos.MaxTags = max(infoArchivos.MaxTags, uint32(len(tag)))
	}

	for _, autor := range a.Metadata.Autores {
		infoArchivos.MaxNombre = max(infoArchivos.MaxNombre, uint32(len(autor.Nombre)))
		infoArchivos.MaxApellido = max(infoArchivos.MaxApellido, uint32(len(autor.Apellido)))
	}

	infoArchivos.MaxNombreLibro = max(infoArchivos.MaxNombreLibro, uint32(len(a.Metadata.TituloObra)))
	infoArchivos.MaxNombreLibro = max(infoArchivos.MaxNombreLibro, uint32(len(a.Metadata.SubtituloObra)))
	for _, capitulo := range a.Metadata.Capitulos {
		infoArchivos.MaxNombreLibro = max(infoArchivos.MaxNombreLibro, uint32(len(capitulo.NombreCapitulo)))
	}

	infoArchivos.MaxEditorial = max(infoArchivos.MaxEditorial, uint32(len(a.Metadata.Editorial)))
	infoArchivos.MaxUrl = max(infoArchivos.MaxUrl, uint32(len(a.Metadata.Url)))
}

func (a *Archivo) InsertarDatos(baseDeDatos *sql.DB) {
	if a.Metadata == nil {
		return
	}

	// Agregar informacion general
	var err error
	if a.IdArchivo, err = Insertar(func() (sql.Result, error) { return baseDeDatos.Exec(insertar_archivo, a.Path) }); err != nil {
		fmt.Printf("Error al obtener insertar el archivo: %s, con error: %v\n", a.Nombre(), err)
		return
	}

	for _, tag := range a.Metadata.Tags {
		if _, err := baseDeDatos.Exec(insertar_tag, tag, a.IdArchivo); err != nil {
			fmt.Printf("Error al insertar tag: %s en el archivo: %s\n", tag, a.Nombre())
		}
	}

	// Agregar informacion especifica

	if a.Metadata.TipoCita == "Libro" {
		if err = CargarDatosDelLibro(baseDeDatos, a.IdArchivo, a.Metadata); err != nil {
			fmt.Printf("Error al insertar libro en el archivo: %s, con error: %v\n", a.Nombre(), err)
		}
	}
}

func CargarDatosDelLibro(bdd *sql.DB, idArchivo int64, meta *Frontmatter) error {
	var idLibro int64
	var err error

	if idEditorial, err := ObtenerOInsertar(
		func() (*sql.Rows, error) { return bdd.Query(query_editorial, meta.Editorial) },
		func() (sql.Result, error) { return bdd.Exec(insertar_editorial, meta.Editorial) },
	); err != nil {
		return fmt.Errorf("error al hacer una querry de la editorial %s con error: %v", meta.Editorial, err)

	} else {
		libro := NewLibro(
			meta.TituloObra,
			meta.SubtituloObra,
			meta.Anio,
			idEditorial,
			meta.Edicion,
			meta.Volumen,
			meta.Url,
			idArchivo,
		)

		if idLibro, err = Insertar(
			func() (sql.Result, error) { return bdd.Exec(insertar_libro, libro.Valores()...) },
		); err != nil {
			return fmt.Errorf("error al insertar un libro, con error: %v", err)
		}
	}

	for _, autor := range meta.NombreAutores {
		if idAutor, err := ObtenerOInsertar(
			func() (*sql.Rows, error) { return bdd.Query(query_personas, autor.Nombre, autor.Apellido) },
			func() (sql.Result, error) { return bdd.Exec(insertar_persona, autor.Nombre, autor.Apellido) },
		); err != nil {
			return fmt.Errorf("error al hacer una querry del autor: %s %s con error: %v", autor.Nombre, autor.Apellido, err)

		} else if _, err := bdd.Exec(insertar_autor_libro, idLibro, idAutor); err != nil {
			return fmt.Errorf("error al insertar par idLibro-idAutor, con error: %v", err)
		}
	}

	for _, capitulo := range meta.Capitulos {
		ids := []any{idLibro, idArchivo}
		var idCapitulo int64

		if idCapitulo, err = Insertar(
			func() (sql.Result, error) { return bdd.Exec(insertar_capitulo, append(capitulo.Valores(), ids...)...) },
		); err != nil {
			return fmt.Errorf("error al insertar un capitulo, con error: %v", err)
		}

		for _, autor := range capitulo.Editores {
			if idAutor, err := ObtenerOInsertar(
				func() (*sql.Rows, error) { return bdd.Query(query_personas, autor.Nombre, autor.Apellido) },
				func() (sql.Result, error) { return bdd.Exec(insertar_persona, autor.Nombre, autor.Apellido) },
			); err != nil {
				return fmt.Errorf("error al hacer una querry del autor: %s %s con error: %v", autor.Nombre, autor.Apellido, err)

			} else if _, err := bdd.Exec(insertar_editor_capitulo, idCapitulo, idAutor); err != nil {
				return fmt.Errorf("error al insertar par idLibro-idAutor, con error: %v", err)
			}
		}
	}

	return nil
}

func (a *Archivo) Nombre() string {
	separacion := strings.Split(a.Path, "/")
	return separacion[len(separacion)-1]
}

func ObtenerOInsertar(query func() (*sql.Rows, error), insert func() (sql.Result, error)) (int64, error) {
	if rows, err := query(); err != nil {
		return 0, fmt.Errorf("error al hacer una querry con error: %v", err)

	} else if rows.Next() {
		var id int64
		rows.Scan(&id)
		return id, nil
	}

	return Insertar(insert)
}

func Insertar(insert func() (sql.Result, error)) (int64, error) {
	if filaAfectada, err := insert(); err != nil {
		return 0, fmt.Errorf("error al insertar con query, con error: %v", err)

	} else if id, err := filaAfectada.LastInsertId(); err != nil {
		return 0, fmt.Errorf("error al obtener id from query, con error: %v", err)

	} else {
		return id, nil
	}
}
