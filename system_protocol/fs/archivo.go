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

const INSERTAR_ARCHIVO = "INSERT INTO archivos (path) VALUES (?)"

const INSERTAR_TAG = "INSERT INTO tags (tag, idArchivo) VALUES (?, ?)"

const QUERY_EDITORIAL = "SELECT id FROM editoriales WHERE editorial = ?"
const INSERTAR_EDITORIAL = "INSERT INTO editoriales (editorial) VALUES (?)"

const QUERY_PERSONAS = "SELECT id FROM personas WHERE nombre = ? AND apellido = ?"
const INSERTAR_PERSONA = "INSERT INTO personas (nombre, apellido) VALUES (?, ?)"

const INSERTAR_LIBRO = "INSERT INTO libros (titulo, subtitulo, anio, idEditorial, edicion, volumen, url, idArchivo) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
const INSERTAR_CAPITULO = "INSERT INTO capitulos (capitulo, nombre, paginaInicial, paginaFinal, idLibro, idArchivo) VALUES (?, ?, ?, ?, ?, ?)"
const INSERTAR_AUTOR_LIBRO = "INSERT INTO autoresLibro (idLibro, idPersona) VALUES (?, ?)"
const INSERTAR_EDITOR_CAPITULO = "INSERT INTO editoresCapitulo (idCapitulo, idPersona) VALUES (?, ?)"

const INSERTAR_DISTRIBUCION = "INSERT INTO distribuciones (nombre, tipo, idArchivo) VALUES (?, ?, ?)"

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

	// Libros
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

	// Distribuciones
}

func (a *Archivo) InsertarDatos(bdd *sql.DB) {
	if a.Metadata == nil {
		return
	}

	// Agregar informacion general
	var err error
	if a.IdArchivo, err = Insertar(func() (sql.Result, error) { return bdd.Exec(INSERTAR_ARCHIVO, a.Path) }); err != nil {
		fmt.Printf("Error al obtener insertar el archivo: %s, con error: %v\n", a.Nombre(), err)
		return
	}

	for _, tag := range a.Metadata.Tags {
		if _, err := bdd.Exec(INSERTAR_TAG, tag, a.IdArchivo); err != nil {
			fmt.Printf("Error al insertar tag: %s en el archivo: %s\n", tag, a.Nombre())
		}
	}

	// Agregar informacion especifica
	if err = CargarDatosDeLasDistribuciones(bdd, a.IdArchivo, a.Metadata); err != nil {
		fmt.Printf("Error al insertar distribuciones en el archivo: %s, con error: %v\n", a.Nombre(), err)
	}

	if err = CargarDatosDelLibro(bdd, a.IdArchivo, a.Metadata); err != nil {
		fmt.Printf("Error al insertar libro en el archivo: %s, con error: %v\n", a.Nombre(), err)
	}

	if err = CargarDatosDelPaper(bdd, a.IdArchivo, a.Metadata); err != nil {
		fmt.Printf("Error al insertar paper en el archivo: %s, con error: %v\n", a.Nombre(), err)
	}
}

func CargarDatosDeLasDistribuciones(bdd *sql.DB, idArchivo int64, meta *Frontmatter) error {
	var tipoDistribucion TipoDistribucion
	switch meta.TipoDistribucion {
	case "discreta":
		tipoDistribucion = DISTRIBUCION_DISCRETA
	case "continua":
		tipoDistribucion = DISTRIBUCION_CONTINUA
	case "multivariada":
		tipoDistribucion = DISTRIBUCION_MULTIVARIADA
	default:
		return nil
	}

	distribucion := NewDistribucion(tipoDistribucion, meta.NombreDistribuucion)
	if _, err := bdd.Exec(INSERTAR_DISTRIBUCION, distribucion.Nombre, distribucion.Tipo, idArchivo); err != nil {
		return fmt.Errorf("error al insertar una distribucion, con error: %v", err)
	}
	// Cuando parsee el texto intentar ver si puedo obtener las relaciones que hay entre las distribuciones

	return nil
}

func CargarDatosDelPaper(bdd *sql.DB, idArchivo int64, meta *Frontmatter) error {
	if meta.TipoCita != "Paper" {
		return nil
	}

	return nil
}

func CargarDatosDelLibro(bdd *sql.DB, idArchivo int64, meta *Frontmatter) error {
	if meta.TipoCita != "Libro" {
		return nil
	}

	var idLibro int64
	var err error

	if idEditorial, err := ObtenerOInsertar(
		func() (*sql.Rows, error) { return bdd.Query(QUERY_EDITORIAL, meta.Editorial) },
		func() (sql.Result, error) { return bdd.Exec(INSERTAR_EDITORIAL, meta.Editorial) },
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
			func() (sql.Result, error) { return bdd.Exec(INSERTAR_LIBRO, libro.Valores()...) },
		); err != nil {
			return fmt.Errorf("error al insertar un libro, con error: %v", err)
		}
	}

	for _, autor := range meta.NombreAutores {
		if idAutor, err := ObtenerOInsertar(
			func() (*sql.Rows, error) { return bdd.Query(QUERY_PERSONAS, autor.Nombre, autor.Apellido) },
			func() (sql.Result, error) { return bdd.Exec(INSERTAR_PERSONA, autor.Nombre, autor.Apellido) },
		); err != nil {
			return fmt.Errorf("error al hacer una querry del autor: %s %s con error: %v", autor.Nombre, autor.Apellido, err)

		} else if _, err := bdd.Exec(INSERTAR_AUTOR_LIBRO, idLibro, idAutor); err != nil {
			return fmt.Errorf("error al insertar par idLibro-idAutor, con error: %v", err)
		}
	}

	for _, capitulo := range meta.Capitulos {
		ids := []any{idLibro, idArchivo}
		var idCapitulo int64

		if idCapitulo, err = Insertar(
			func() (sql.Result, error) { return bdd.Exec(INSERTAR_CAPITULO, append(capitulo.Valores(), ids...)...) },
		); err != nil {
			return fmt.Errorf("error al insertar un capitulo, con error: %v", err)
		}

		for _, autor := range capitulo.Editores {
			if idAutor, err := ObtenerOInsertar(
				func() (*sql.Rows, error) { return bdd.Query(QUERY_PERSONAS, autor.Nombre, autor.Apellido) },
				func() (sql.Result, error) { return bdd.Exec(INSERTAR_PERSONA, autor.Nombre, autor.Apellido) },
			); err != nil {
				return fmt.Errorf("error al hacer una querry del autor: %s %s con error: %v", autor.Nombre, autor.Apellido, err)

			} else if _, err := bdd.Exec(INSERTAR_EDITOR_CAPITULO, idCapitulo, idAutor); err != nil {
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
