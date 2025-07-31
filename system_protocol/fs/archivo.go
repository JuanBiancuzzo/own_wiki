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

const TAG_CARRERA = "facultad/carrera"
const TAG_MATERIA = "facultad/materia"
const TAG_MATERIA_EQUIVALENTE = "facultad/materia-equivalente"
const TAG_RESUMEN_MATERIA = "facultad/resumen"
const TAG_DISTRIBUCION = "colección/distribuciones/distribución"
const TAG_LIBRO = "colección/biblioteca/libro"

type Archivo struct {
	Padre                            *Directorio
	Path                             string
	Archivo                          *e.Archivo
	Carrera                          *e.ConstructorCarrera
	Materia                          *e.ConstructorMateria
	MateriaEquivalente               *e.ConstructorMateriaEquivalente
	MateriasCorrelativas             *l.Lista[*e.ConstructorMateriasCorrelativas]
	MateriasEquivalentesCorrelativas *l.Lista[*e.ConstructorMateriasCorrelativas]
}

func NewArchivo(padre *Directorio, path string, info *db.InfoArchivos, canal chan string) (*Archivo, error) {
	archivo := Archivo{
		Padre:                            padre,
		Path:                             path,
		MateriasCorrelativas:             l.NewLista[*e.ConstructorMateriasCorrelativas](),
		MateriasEquivalentesCorrelativas: l.NewLista[*e.ConstructorMateriasCorrelativas](),
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
			if constructor, err := e.NewConstructorMateria(meta.PathCarrera, meta.NombreMateria, meta.Codigo, meta.Plan, meta.Cuatri, meta.Etapa); err == nil {
				archivo.Archivo.CargarDependencia(constructor.CumpleDependenciaArchivo)

			} else {
				canal <- fmt.Sprintf("Error: %v\n", err)
			}

			for _, correlativa := range meta.Correlativas {
				pathCorrelativa := ObtenerWikiLink(correlativa)[0]
				constructor := e.NewConstructorMateriasCorrelativas(e.MATERIA_REAL, pathCorrelativa)
				archivo.MateriasCorrelativas.Push(constructor)
			}

		case TAG_MATERIA_EQUIVALENTE:
			pathMateria := ObtenerWikiLink(meta.Equivalencia)[0]
			constructor := e.NewConstructorMateriaEquivalente(pathMateria, meta.NombreMateria, meta.Codigo)
			archivo.Archivo.CargarDependencia(constructor.CumpleDependenciaArchivo)

			for _, correlativa := range meta.Correlativas {
				pathCorrelativa := ObtenerWikiLink(correlativa)[0]
				constructor := e.NewConstructorMateriasCorrelativas(e.MATERIA_EQUIVALENTE, pathCorrelativa)
				archivo.MateriasCorrelativas.Push(constructor)
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

func (a *Archivo) RelativizarPath(path string) {
	a.Path = strings.Replace(a.Path, path, "", 1)
}

// Cambiar a establecer conexiones
func (a *Archivo) InsertarDatos(canal chan e.Cargable, canalMensajes chan string) {
	if a.Materia != nil {
		if archivo, err := a.Padre.EncontrarArchivo(a.Materia.PathCarrera); err != nil {
			canalMensajes <- fmt.Sprintf("Error al buscar carrera en '%s' en la materia '%s', con error %v", a.Materia.PathCarrera, a.Materia.Nombre, err)

		} else if archivo.Carrera == nil {
			canalMensajes <- fmt.Sprintf("Error el archivo de carrera '%s' no tiene la estructura de carrera, con error %v", a.Materia.PathCarrera, err)

		} else {
			archivo.Carrera.CargarDependencia(a.Materia.CumpleDependenciaCarrera)
		}

		for correlativa := range a.MateriasCorrelativas.Iterar {
			if archivo, err := a.Padre.EncontrarArchivo(correlativa.PathCorrelativa); err != nil {
				canalMensajes <- fmt.Sprintf("Error al buscar correlativa en '%s' en la materia '%s', con error %v", correlativa.PathCorrelativa, a.Materia.Nombre, err)

			} else if archivo.Materia != nil {
				correlativa.TipoCorrelativa = e.MATERIA_REAL
				a.Materia.CargarDependencia(correlativa.CumpleDependenciaMateria)
				archivo.Materia.CargarDependencia(correlativa.CumpleDependenciaMateria)

			} else if archivo.MateriaEquivalente != nil {
				correlativa.TipoCorrelativa = e.MATERIA_EQUIVALENTE
				a.Materia.CargarDependencia(correlativa.CumpleDependenciaMateria)
				archivo.MateriaEquivalente.CargarDependencia(correlativa.CumpleDependenciaMateria)

			} else {
				canalMensajes <- fmt.Sprintf("Error el archivo de materia '%s' no tiene la estructura de materi, con error %v", a.MateriaEquivalente.PathMateria, err)
			}
		}
	}

	if a.MateriaEquivalente != nil {
		if archivo, err := a.Padre.EncontrarArchivo(a.MateriaEquivalente.PathMateria); err != nil {
			canalMensajes <- fmt.Sprintf("Error al buscar materia en '%s' en la materia equivalente '%s', con error %v", a.MateriaEquivalente.PathMateria, a.MateriaEquivalente.Nombre, err)

		} else if archivo.Materia == nil {
			canalMensajes <- fmt.Sprintf("Error el archivo de materia '%s' no tiene la estructura de materi, con error %v", a.MateriaEquivalente.PathMateria, err)

		} else {
			archivo.Materia.CargarDependencia(a.MateriaEquivalente.CumpleDependenciaMateria)
		}

		for correlativa := range a.MateriasEquivalentesCorrelativas.Iterar {
			if archivo, err := a.Padre.EncontrarArchivo(correlativa.PathCorrelativa); err != nil {
				canalMensajes <- fmt.Sprintf("Error al buscar correlativa en '%s' en la materia '%s', con error %v", correlativa.PathCorrelativa, a.Materia.Nombre, err)

			} else if archivo.Materia != nil {
				correlativa.TipoCorrelativa = e.MATERIA_REAL
				a.Materia.CargarDependencia(correlativa.CumpleDependenciaMateria)
				archivo.Materia.CargarDependencia(correlativa.CumpleDependenciaMateria)

			} else if archivo.MateriaEquivalente != nil {
				correlativa.TipoCorrelativa = e.MATERIA_EQUIVALENTE
				a.Materia.CargarDependencia(correlativa.CumpleDependenciaMateria)
				archivo.MateriaEquivalente.CargarDependencia(correlativa.CumpleDependenciaMateria)

			} else {
				canalMensajes <- fmt.Sprintf("Error el archivo de materia '%s' no tiene la estructura de materi, con error %v", a.MateriaEquivalente.PathMateria, err)
			}
		}
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
