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
const TAG_CURSO = "cursos/curso"
const TAG_CURSO_PRESENCIA = "cursos/curso-presencial"
const TAG_RESUMEN_CURSO = "cursos/resumen"
const TAG_DISTRIBUCION = "colección/distribuciones/distribución"
const TAG_LIBRO = "colección/biblioteca/libro"
const TAG_PAPER = "colección/biblioteca/paper"

type PathTipo struct {
	Path string
	Tipo e.TipoDependible
}

func NewPathTipo(path string, tipo e.TipoDependible) PathTipo {
	return PathTipo{
		Path: path,
		Tipo: tipo,
	}
}

type Archivo struct {
	Root           *Root
	Path           string
	Cargables      *l.Lista[e.Cargable]
	FnDependencias map[PathTipo]*l.Lista[e.FnVincular]
	Dependibles    map[e.TipoDependible]*l.Lista[e.Dependible]
}

func NewArchivo(root *Root, path string, info *db.InfoArchivos, canal chan string) (*Archivo, error) {
	archivo := Archivo{
		Root:           root,
		Path:           path,
		Cargables:      l.NewLista[e.Cargable](),
		FnDependencias: make(map[PathTipo]*l.Lista[e.FnVincular]),
		Dependibles:    make(map[e.TipoDependible]*l.Lista[e.Dependible]),
	}

	if !strings.Contains(path, ".md") {
		return &archivo, nil
	}

	var contenido string
	if bytes, err := os.ReadFile(fmt.Sprintf("%s/%s", root.Path, path)); err != nil {
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
	archivoCargable := e.NewArchivo(path, meta.Tags)

	archivo.Cargables.Push(archivoCargable)
	archivo.CargarDependible(e.DEP_ARCHIVO, archivoCargable)

	for _, tag := range meta.Tags {
		switch tag {
		case TAG_CARRERA:
			nombreCarrera := archivo.Nombre()
			if constructor, err := e.NewConstructorCarrera(nombreCarrera, meta.Etapa, meta.TieneCodigo); err == nil {
				archivo.CargarDependencia(path, e.DEP_ARCHIVO, constructor.CrearDependenciaArchivo)

				archivo.CargarDependible(e.DEP_CARRERA, constructor)
			} else {
				canal <- fmt.Sprintf("Error: %v\n", err)
			}

		case TAG_MATERIA:
			if constructor, err := e.NewConstructorMateria(meta.NombreMateria, meta.Codigo, meta.Plan, meta.Cuatri, meta.Etapa); err == nil {
				pathCarrera := ObtenerWikiLink(meta.PathCarrera)[0]
				archivo.CargarDependencia(path, e.DEP_ARCHIVO, constructor.CrearDependenciaArchivo)
				archivo.CargarDependencia(pathCarrera, e.DEP_CARRERA, constructor.CrearDependenciaCarrera)

				archivo.CargarDependible(e.DEP_MATERIA, constructor)
			} else {
				canal <- fmt.Sprintf("Error: %v\n", err)
			}

			for _, correlativa := range meta.Correlativas {
				constructor := e.NewConstructorMateriasCorrelativas(e.MATERIA_REAL, correlativa.Tipo)

				archivo.CargarDependencia(path, e.DEP_MATERIA, constructor.CrearDependenciaMateria)
				switch correlativa.Tipo {
				case e.MATERIA_REAL:
					archivo.CargarDependencia(correlativa.Path, e.DEP_MATERIA, constructor.CrearDependenciaCorrelativa)
				case e.MATERIA_EQUIVALENTE:
					archivo.CargarDependencia(correlativa.Path, e.DEP_MATERIA_EQUIVALENTE, constructor.CrearDependenciaCorrelativa)
				}
			}

		case TAG_MATERIA_EQUIVALENTE:
			constructor := e.NewConstructorMateriaEquivalente(meta.NombreMateria, meta.Codigo)
			archivo.CargarDependencia(path, e.DEP_ARCHIVO, constructor.CrearDependenciaArchivo)
			archivo.CargarDependencia(ObtenerWikiLink(meta.Equivalencia)[0], e.DEP_MATERIA, constructor.CrearDependenciaMateria)

			archivo.CargarDependible(e.DEP_MATERIA_EQUIVALENTE, constructor)

			for _, correlativa := range meta.Correlativas {
				constructor := e.NewConstructorMateriasCorrelativas(e.MATERIA_EQUIVALENTE, correlativa.Tipo)

				archivo.CargarDependencia(path, e.DEP_MATERIA_EQUIVALENTE, constructor.CrearDependenciaMateria)
				switch correlativa.Tipo {
				case e.MATERIA_REAL:
					archivo.CargarDependencia(correlativa.Path, e.DEP_MATERIA, constructor.CrearDependenciaCorrelativa)
				case e.MATERIA_EQUIVALENTE:
					archivo.CargarDependencia(correlativa.Path, e.DEP_MATERIA_EQUIVALENTE, constructor.CrearDependenciaCorrelativa)
				}
			}

		case TAG_RESUMEN_MATERIA:
			constructor := e.NewConstructorTemaMateria(meta.NombreResumen, meta.Capitulo, meta.Parte)
			archivo.CargarDependencia(path, e.DEP_ARCHIVO, constructor.CrearDependenciaArchivo)
			archivo.CargarDependencia(meta.MateriaResumen, e.DEP_MATERIA, constructor.CrearDependenciaMateria)

			archivo.CargarDependible(e.DEP_TEMA_MATERIA, constructor)

		case TAG_CURSO:
			if constructor, err := meta.CrearConstructorCurso(); err == nil {
				archivo.CargarDependencia(path, e.DEP_ARCHIVO, constructor.CrearDependenciaArchivo)

				archivo.CargarDependible(e.DEP_CURSO, constructor)
			} else {
				canal <- fmt.Sprintf("Error: %v\n", err)
			}

		case TAG_CURSO_PRESENCIA:
			if constructor, err := meta.CrearConstructorCursoPresencial(); err == nil {
				archivo.CargarDependencia(path, e.DEP_ARCHIVO, constructor.CrearDependenciaArchivo)

				archivo.CargarDependible(e.DEP_CURSO_PRESENCIAL, constructor)
			} else {
				canal <- fmt.Sprintf("Error: %v\n", err)
			}

		case TAG_RESUMEN_CURSO:
			constructor := e.NewConstructorTemaCurso(meta.NombreResumen, meta.Capitulo, meta.Parte, meta.TipoCurso)
			archivo.CargarDependencia(path, e.DEP_ARCHIVO, constructor.CrearDependenciaArchivo)

			pathCurso := ObtenerWikiLink(meta.Curso)[0]
			switch meta.TipoCurso {
			case e.CURSO_ONLINE:
				archivo.CargarDependencia(pathCurso, e.DEP_CURSO, constructor.CrearDependenciaCurso)
			case e.CURSO_PRESENCIAL:
				archivo.CargarDependencia(pathCurso, e.DEP_CURSO_PRESENCIAL, constructor.CrearDependenciaCurso)
			}

			archivo.CargarDependible(e.DEP_TEMA_CURSO, constructor)

		case TAG_LIBRO:
			constructor := meta.CrearConstructorLibro()
			archivo.CargarDependencia(path, e.DEP_ARCHIVO, constructor.CrearDependenciaArchivo)

		case TAG_PAPER:
			if constructor, err := meta.CrearConstructorPaper(); err == nil {
				archivo.CargarDependencia(path, e.DEP_ARCHIVO, constructor.CrearDependenciaArchivo)

			} else {
				canal <- fmt.Sprintf("Error: %v\n", err)
			}

		case TAG_DISTRIBUCION:
			if constructor, err := e.NewConstructorDistribucion(meta.NombreDistribuucion, meta.TipoDistribucion); err == nil {
				archivo.CargarDependencia(path, e.DEP_ARCHIVO, constructor.CrearDependenciaArchivo)

			} else {
				canal <- fmt.Sprintf("Error: %v\n", err)
			}
		}
	}

	info.MaxPath = max(info.MaxPath, uint32(len(path)))
	CargarInfo(info, &meta)

	return &archivo, nil
}

func (a *Archivo) CargarDependible(tipo e.TipoDependible, dependible e.Dependible) {
	lista, ok := a.Dependibles[tipo]
	if !ok {
		lista = l.NewLista[e.Dependible]()
	}
	lista.Push(dependible)
	a.Dependibles[tipo] = lista
}

func (a *Archivo) CargarDependencia(path string, tipo e.TipoDependible, relacion e.FnVincular) {
	pathTipo := NewPathTipo(path, tipo)
	lista, ok := a.FnDependencias[pathTipo]
	if !ok {
		lista = l.NewLista[e.FnVincular]()
	}
	lista.Push(relacion)
	a.FnDependencias[pathTipo] = lista
}

// Cambiar a establecer conexiones
func (a *Archivo) EstablecerDependencias(canal chan e.Cargable, canalMensajes chan string) {
	for pathTipo, listaRelacion := range a.FnDependencias {
		if archivo, err := a.Root.EncontrarArchivo(pathTipo.Path); err != nil {
			canalMensajes <- fmt.Sprintf("No se encontró el archivo '%s' para el archivo: '%s'", pathTipo.Path, a.Nombre())

		} else if listaDependible, ok := archivo.Dependibles[pathTipo.Tipo]; !ok {
			canalMensajes <- fmt.Sprintf("No se encontró el tipo '%d' para el archivo: '%s'", pathTipo.Tipo, a.Nombre())

		} else {
			for dependible := range listaDependible.Iterar {
				for relacion := range listaRelacion.Iterar {
					relacion(dependible)
				}
			}
		}
	}

	for cargable := range a.Cargables.Iterar {
		canal <- cargable
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
