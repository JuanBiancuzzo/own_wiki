package fs

import (
	"database/sql"
	"fmt"
	"os"
	"own_wiki/system_protocol/db"
	"strconv"
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
var insertar_autor_libro = "INSERT INTO autoresLibro (idLibro, idPersona) VALUES (?, ?)"

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

func (a *Archivo) Interprestarse(infoArchivos *db.InfoArchivos) {
	if !strings.Contains(a.Path, ".md") {
		return
	}

	bytes, err := os.ReadFile(a.Path)
	if err != nil {
		fmt.Printf("Error al leer %s obteniendo el error: %v", a.Path, err)
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
		fmt.Printf("Error al decodificar en %s la metadata, con el error: %v\n", a.Path, err)
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

	if a.Metadata.TipoCita == "Libro" {
		a.CargarDatosDelLibro(baseDeDatos)
	}
}

func (a *Archivo) CargarDatosDelLibro(baseDeDatos *sql.DB) {
	var libro Libro
	if idEditorial, err := ObtenerOInsertar(
		func() (*sql.Rows, error) { return baseDeDatos.Query(query_editorial, a.Metadata.Editorial) },
		func() (sql.Result, error) { return baseDeDatos.Exec(insertar_editorial, a.Metadata.Editorial) },
	); err != nil {
		fmt.Printf("Error al hacer una querry de la editorial %s con error: %v\n", a.Metadata.Editorial, err)
		return

	} else {
		libro = Libro{
			Titulo:      a.Metadata.TituloObra,
			Subtitulo:   a.Metadata.SubtituloObra,
			Anio:        NumeroODefault(a.Metadata.Anio, -1),
			IdEditorial: idEditorial,
			Edicion:     NumeroODefault(a.Metadata.Edicion, 1),
			Volumen:     NumeroODefault(a.Metadata.Volumen, 1),
			Url:         a.Metadata.Url,
			IdArchivo:   a.IdArchivo,
		}
	}

	for _, autor := range a.Metadata.NombreAutores {
		if idAutor, err := ObtenerOInsertar(
			func() (*sql.Rows, error) { return baseDeDatos.Query(query_personas, autor.Nombre, autor.Apellido) },
			func() (sql.Result, error) { return baseDeDatos.Exec(insertar_persona, autor.Nombre, autor.Apellido) },
		); err != nil {
			fmt.Printf("Error al hacer una querry del autor: %s %s con error: %v\n", autor.Nombre, autor.Apellido, err)

		} else if idLibro, err := Insertar(
			func() (sql.Result, error) { return baseDeDatos.Exec(insertar_libro, libro.Valores()...) },
		); err != nil {
			fmt.Printf("Error al insertar un libro en el archivo: %s\n, con error: %v", a.Nombre(), err)

		} else if _, err := baseDeDatos.Exec(insertar_autor_libro, idLibro, idAutor); err != nil {
			fmt.Printf("Error al insertar par idLibro-idAutor en el archivo: %s\n, con error: %v", a.Nombre(), err)
		}
	}
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

	} else {
		return Insertar(insert)
	}
}

func Insertar(insert func() (sql.Result, error)) (int64, error) {
	if filaAfectada, err := insert(); err != nil {
		return 0, fmt.Errorf("error al insertar con query")

	} else if id, err := filaAfectada.LastInsertId(); err != nil {
		return 0, fmt.Errorf("error al obtener id from query")

	} else {
		return id, nil
	}
}

func NumeroODefault(representacion string, valorDefault int) int {
	if nuevoValor, err := strconv.Atoi(representacion); err == nil {
		return nuevoValor
	} else {
		return valorDefault
	}
}
