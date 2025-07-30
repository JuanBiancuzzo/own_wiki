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
	Padre          *Directorio
	Path           string
	Meta           *Frontmatter
	TiposDeArchivo *l.Lista[e.Cargable]
}

func NewArchivo(padre *Directorio, path string, info *db.InfoArchivos) (*Archivo, error) {
	archivo := Archivo{
		Padre:          padre,
		Path:           path,
		Meta:           nil,
		TiposDeArchivo: l.NewLista[e.Cargable](),
	}

	if !strings.Contains(path, ".md") {
		return &archivo, nil
	}

	var contenido string
	if bytes, err := os.ReadFile(path); err != nil {
		return nil, fmt.Errorf("error al leer %s obteniendo el error: %v", path, err)
	} else {
		contenido = strings.TrimSpace(string(bytes))
	}

	if strings.Index(contenido, "---") != 0 {
		return &archivo, nil
	}

	blob := contenido[3 : 3+strings.Index(contenido[3:], "---")]
	decodificador := yaml.NewDecoder(strings.NewReader(blob))

	var meta Frontmatter
	if err := decodificador.Decode(&meta); err != nil {
		return nil, fmt.Errorf("error al decodificar en %s la metadata, con el error: %v", path, err)
	}

	// a.Contenido = contenido[3+strings.Index(contenido[3:], "---")+len("---"):]

	info.MaxPath = max(info.MaxPath, uint32(len(path)))

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

	archivo.Meta = &meta
	return &archivo, nil
}

func (a *Archivo) Interprestarse(canal chan string) {
	if a.Meta == nil {
		return
	}

	a.TiposDeArchivo.Push(e.NewArchivo(a.Path, a.Meta.Tags))
	for _, tag := range a.Meta.Tags {
		switch tag {
		case "facultad/carrera":
			nombreCarrera := a.Nombre()
			if carrera, err := e.NewCarrera(a.Path, nombreCarrera, a.Meta.Etapa, a.Meta.TieneCodigo); err == nil {
				a.TiposDeArchivo.Push(carrera)

			} else {
				canal <- fmt.Sprintf("Error: %v\n", err)
			}
		case "facultad/materia":
			if carrera, err := a.Padre.ArchivoMasCercanoConTag("facultad/carrera"); err != nil {
				canal <- fmt.Sprintf("Error al buscar carrera, con error: %v\n", err)

			} else if materia, err := e.NewMateria(a.Path, carrera.Path, a.Meta.NombreMateria, a.Meta.Codigo, a.Meta.Plan, a.Meta.Cuatri, a.Meta.Etapa); err != nil {
				canal <- fmt.Sprintf("Error: %v\n", err)

			} else {
				a.TiposDeArchivo.Push(materia)
			}

			/*
				for _, correlativa := range a.Meta.Correlativas {
					pathMateria := ObtenerWikiLink(correlativa)[0]
					a.TiposDeArchivo.Push(e.NewMateriaEquivalente(a.Path, pathMateria, a.Meta.NombreMateria, a.Meta.Codigo))
				}
			*/
		case "facultad/materia-equivalente":
			pathMateria := ObtenerWikiLink(a.Meta.Equivalencia)[0]
			a.TiposDeArchivo.Push(e.NewMateriaEquivalente(a.Path, pathMateria, a.Meta.NombreMateria, a.Meta.Codigo))
		case "facultad/resumen":
			// a.TiposDeArchivo.Push(ES_RESUMEN_MATERIA)
		case "colección/distribuciones/distribución":
			if distribucion, err := e.NewDistribucion(a.Path, a.Meta.NombreDistribuucion, a.Meta.TipoDistribucion); err == nil {
				a.TiposDeArchivo.Push(distribucion)

			} else {
				canal <- fmt.Sprintf("Error: %v\n", err)
			}
		case "colección/biblioteca/libro":
			a.TiposDeArchivo.Push(a.Meta.CrearLibro(a.Path))
		}
	}
}

func (a *Archivo) EstablecerInfo(info *db.InfoArchivos, meta *Frontmatter) {
	// General

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

func ObtenerWikiLink(link string) []string {
	link = strings.TrimPrefix(link, "[[")
	link = strings.TrimSuffix(link, "]]")
	return strings.Split(link, "|")
}

/*
	case ES_MATERIA:
		for _, correlativa := range meta.Correlativas {
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
