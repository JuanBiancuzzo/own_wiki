package fs

import (
	"fmt"
	"os"
	"own_wiki/system_protocol/db"
	e "own_wiki/system_protocol/estructura"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
)

const TAG_CARRERA = "facultad/carrera"
const TAG_MATERIA = "facultad/materia"
const TAG_MATERIA_EQUIVALENTE = "facultad/materia-equivalente"
const TAG_RESUMEN_MATERIA = "facultad/resumen"
const TAG_DISTRIBUCION = "colección/distribuciones/distribución"
const TAG_LIBRO = "colección/biblioteca/libro"

type Archivo struct {
	Padre   *Directorio
	Path    string
	Archivo *e.Archivo
	Carrera *e.ConstructorCarrera
	Materia *e.ConstructorMateria
}

func NewArchivo(padre *Directorio, path string, info *db.InfoArchivos, canal chan string) (*Archivo, error) {
	archivo := Archivo{
		Padre:   padre,
		Path:    path,
		Archivo: nil,
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

	archivo.Archivo = e.NewArchivo(path, meta.Tags)
	for _, tag := range meta.Tags {
		switch tag {
		case TAG_DISTRIBUCION:
			if constructor, err := e.NewConstructorDistribucion(meta.NombreDistribuucion, meta.TipoDistribucion); err == nil {
				archivo.Archivo.CargarDependencia(constructor.CumpleDependencia)

			} else {
				canal <- fmt.Sprintf("Error: %v\n", err)
			}
		case TAG_CARRERA:
			nombreCarrera := archivo.Nombre()
			if constructor, err := e.NewConstructorCarrera(nombreCarrera, meta.Etapa, meta.TieneCodigo); err == nil {
				archivo.Archivo.CargarDependencia(constructor.CumpleDependencia)

			} else {
				canal <- fmt.Sprintf("Error: %v\n", err)
			}

		case TAG_MATERIA:
			if constructor, err := e.NewConstructorMateria(meta.NombreMateria, meta.Codigo, meta.Plan, meta.Cuatri, meta.Etapa); err == nil {
				archivo.Archivo.CargarDependencia(constructor.CumpleDependenciaArchivo)

			} else {
				canal <- fmt.Sprintf("Error: %v\n", err)
			}

		case TAG_LIBRO:
			constructor := meta.CrearLibro()
			archivo.Archivo.CargarDependencia(constructor.CumpleDependencia)
		}
	}

	info.MaxPath = max(info.MaxPath, uint32(len(path)))
	CargarInfo(info, &meta)

	return &archivo, nil
}

/*
func (a *Archivo) Interprestarse(canal chan string) {


	a.TiposDeArchivo.Push(e.NewArchivo(a.Path, a.Meta.Tags))
	for _, tag := range a.Meta.Tags {
		switch tag {}
		case TAG_MATERIA:
			if carrera, err := a.Padre.ArchivoMasCercanoConTag(TAG_CARRERA); err != nil {
				canal <- fmt.Sprintf("Error al buscar carrera, con error: %v\n", err)

			} else if materia, err := e.NewMateria(a.Path, carrera.Path, a.Meta.NombreMateria, a.Meta.Codigo, a.Meta.Plan, a.Meta.Cuatri, a.Meta.Etapa); err != nil {
				canal <- fmt.Sprintf("Error: %v\n", err)

			} else {
				a.TiposDeArchivo.Push(materia)
			}

			for _, correlativa := range a.Meta.Correlativas {
				pathCorrelativa := ObtenerWikiLink(correlativa)[0]
				if archivo, err := a.Padre.EncontrarArchivo(pathCorrelativa); err != nil {
					canal <- fmt.Sprintf("No existe el archivo: %s, con error: %v", pathCorrelativa, err)

				} else {
					var tipoCorrelativa e.TipoMateria = e.MATERIA_REAL
					if slices.Contains(archivo.Meta.Tags, TAG_MATERIA_EQUIVALENTE) {
						tipoCorrelativa = e.MATERIA_EQUIVALENTE
					}

					a.TiposDeArchivo.Push(e.NewMateriasCorrelativas(a.Path, e.MATERIA_REAL, pathCorrelativa, tipoCorrelativa))
				}
			}
		case TAG_MATERIA_EQUIVALENTE:
			pathMateria := ObtenerWikiLink(a.Meta.Equivalencia)[0]
			a.TiposDeArchivo.Push(e.NewMateriaEquivalente(a.Path, pathMateria, a.Meta.NombreMateria, a.Meta.Codigo))

			for _, correlativa := range a.Meta.Correlativas {
				pathCorrelativa := ObtenerWikiLink(correlativa)[0]

				if archivo, err := a.Padre.EncontrarArchivo(pathCorrelativa); err != nil {
					canal <- fmt.Sprintf("No existe el archivo: %s, con error: %v", pathCorrelativa, err)

				} else {
					var tipoCorrelativa e.TipoMateria = e.MATERIA_REAL
					if slices.Contains(archivo.Meta.Tags, TAG_MATERIA_EQUIVALENTE) {
						tipoCorrelativa = e.MATERIA_EQUIVALENTE
					}

					a.TiposDeArchivo.Push(e.NewMateriasCorrelativas(a.Path, e.MATERIA_EQUIVALENTE, pathCorrelativa, tipoCorrelativa))
				}
			}
		}
	}
}
*/

func (a *Archivo) RelativizarPath(path string) {
	a.Path = strings.Replace(a.Path, path, "", 1)
}

// Cambiar a establecer conexiones
func (a *Archivo) InsertarDatos(canal chan e.Cargable) {
	if a.Materia != nil {
		// Buscar carrera
		// (esto es un ejemplo)
		a.Carrera.CargarDependencia(a.Materia.CumpleDependenciaCarrera)

	}

	canal <- a.Archivo
}

func (a *Archivo) Nombre() string {
	return e.Nombre(a.Path)
}

func ObtenerWikiLink(link string) []string {
	link = strings.TrimPrefix(link, "[[")
	link = strings.TrimSuffix(link, "]]")
	return strings.Split(link, "|")
}
