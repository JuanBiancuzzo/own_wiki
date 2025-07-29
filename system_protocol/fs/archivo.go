package fs

import (
	"fmt"
	"os"
	"own_wiki/system_protocol/db"
	e "own_wiki/system_protocol/estructura"
	l "own_wiki/system_protocol/listas"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
)

type Archivo struct {
	Path           string
	TiposDeArchivo *l.Lista[e.Cargable]
}

func NewArchivo(path string) *Archivo {
	return &Archivo{
		Path:           path,
		TiposDeArchivo: l.NewLista[e.Cargable](),
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

	meta := new(Frontmatter)
	err = decodificador.Decode(meta)
	if err != nil {
		canal <- fmt.Sprintf("Error al decodificar en %s la metadata, con el error: %v\n", a.Path, err)
		return
	}

	// a.Contenido = contenido[3+strings.Index(contenido[3:], "---")+len("---"):]

	a.EstablecerInfo(infoArchivos, meta)
	if len(meta.Tags) == 0 {
		return
	}

	a.TiposDeArchivo.Push(e.NewArchivo(a.Path, meta.Tags))
	for _, tag := range meta.Tags {
		switch tag {
		case "facultad/carrera":
			nombreCarrera := a.Nombre()
			if carrera, err := e.NewCarrera(a.Path, nombreCarrera, meta.Etapa, meta.TieneCodigo); err == nil {
				a.TiposDeArchivo.Push(carrera)

			} else {
				canal <- fmt.Sprintf("Error: %v\n", err)
			}
		case "facultad/materia":
			seccionesPath := strings.Split(a.Path, "/")
			carpetaPrevia := strings.Join(seccionesPath[:len(seccionesPath)-2], "/")
			carpetaPrevia = fmt.Sprintf("%s/%s", carpetaPrevia, "%")

			if materia, err := e.NewMateria(a.Path, carpetaPrevia, meta.NombreMateria, meta.Codigo, meta.Plan, meta.Cuatri, meta.Etapa); err == nil {
				a.TiposDeArchivo.Push(materia)

			} else {
				canal <- fmt.Sprintf("Error: %v\n", err)
			}
		case "facultad/materia-equivalente":
			// a.TiposDeArchivo.Push(ES_MATERIA_EQUIVALENTE)
		case "facultad/resumen":
			// a.TiposDeArchivo.Push(ES_RESUMEN_MATERIA)
		case "colección/distribuciones/distribución":
			if distribucion, err := e.NewDistribucion(a.Path, meta.NombreDistribuucion, meta.TipoDistribucion); err == nil {
				a.TiposDeArchivo.Push(distribucion)

			} else {
				canal <- fmt.Sprintf("Error: %v\n", err)
			}
		case "colección/biblioteca/libro":
			a.TiposDeArchivo.Push(meta.CrearLibro(a.Path))
		}
	}
}

func (a *Archivo) EstablecerInfo(info *db.InfoArchivos, meta *Frontmatter) {
	// General
	info.MaxPath = max(info.MaxPath, uint32(len(a.Path)))

	if meta == nil {
		return
	}

	for _, tag := range meta.Tags {
		info.MaxTags = max(info.MaxTags, uint32(len(tag)))
	}

	// Libros
	for _, autor := range meta.Autores {
		info.MaxNombre = max(info.MaxNombre, uint32(len(autor.Nombre)))
		info.MaxApellido = max(info.MaxApellido, uint32(len(autor.Apellido)))
	}

	info.MaxNombreLibro = max(info.MaxNombreLibro, uint32(len(meta.TituloObra)))
	info.MaxNombreLibro = max(info.MaxNombreLibro, uint32(len(meta.SubtituloObra)))
	for _, capitulo := range meta.Capitulos {
		info.MaxNombreLibro = max(info.MaxNombreLibro, uint32(len(capitulo.NombreCapitulo)))
	}

	info.MaxEditorial = max(info.MaxEditorial, uint32(len(meta.Editorial)))
	info.MaxUrl = max(info.MaxUrl, uint32(len(meta.Url)))

	// Distribuciones

}

func (a *Archivo) RelativizarPath(path string) {
	a.Path = strings.Replace(a.Path, path, "", 1)
}

func (a *Archivo) InsertarDatos(canal chan e.Cargable) {
	for tipoArchivo := range a.TiposDeArchivo.Iterar {
		canal <- tipoArchivo
	}
}

func (a *Archivo) Nombre() string {
	return e.Nombre(a.Path)
}

/*
	case ES_MATERIA:
		for _, correlativa := range meta.Correlativas {
			correlativa = strings.TrimPrefix(correlativa, "[[")
			correlativa = strings.TrimSuffix(correlativa, "]]")
			pathCorrelativa := strings.Split(correlativa, "|")[0]
			canal <- func(bdd *sql.DB) bool { return CargarCorrelativas(bdd, pathArchivo, pathCorrelativa) }
		}
	case ES_MATERIA_EQUIVALENTE:
		canal <- func(bdd *sql.DB) bool {
			fmt.Println("Insertando materia: ", nombreArchivo)
			if idCarrera, existe := ExisteArchivoCarpetaPrevia(bdd, "carreras", pathArchivo); !existe {
				return false
			} else if err = CargarDatosDeLaMateria(bdd, idArchivo, idCarrera, meta); err != nil {
				fmt.Printf("Error al insertar una materia en el archivo: %s, con error: %v\n", nombreArchivo, err)
			}
			return true
		}
		for _, correlativa := range meta.Correlativas {
			correlativa = strings.TrimPrefix(correlativa, "[[")
			correlativa = strings.TrimSuffix(correlativa, "]]")
			pathCorrelativa := strings.Split(correlativa, "|")[0]
			canal <- func(bdd *sql.DB) bool { return CargarCorrelativas(bdd, pathArchivo, pathCorrelativa) }
		}
	}
*/
