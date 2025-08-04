package fs

import (
	"fmt"
	"os"
	bdd "own_wiki/system_protocol/bass_de_datos"
	e "own_wiki/system_protocol/datos"
	l "own_wiki/system_protocol/utilidades"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
)

const (
	TAG_CARRERA             = "facultad/carrera"
	TAG_MATERIA             = "facultad/materia"
	TAG_MATERIA_EQUIVALENTE = "facultad/materia-equivalente"
	TAG_RESUMEN_MATERIA     = "facultad/resumen"
)

const (
	TAG_CURSO           = "cursos/curso"
	TAG_CURSO_PRESENCIA = "cursos/curso-presencial"
	TAG_RESUMEN_CURSO   = "cursos/resumen"
)

const (
	TAG_DISTRIBUCION = "colección/distribuciones/distribución"
	TAG_LIBRO        = "colección/biblioteca/libro"
	TAG_PAPER        = "colección/biblioteca/paper"
)

const (
	TAG_NOTA_FACULTAD      = "nota/facultad"
	TAG_NOTA_CURSO         = "nota/curso"
	TAG_NOTA_INVESTIGACION = "nota/investigacion"
	TAG_NOTA_COLECCION     = "nota/colección"
	TAG_NOTA_PROYECTO      = "nota/proyecto"
)

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

func NewArchivo(root *Root, path string, info *bdd.InfoArchivos, canalMensajes chan string) (*Archivo, error) {
	a := Archivo{
		Root:           root,
		Path:           path,
		Cargables:      l.NewLista[e.Cargable](),
		FnDependencias: make(map[PathTipo]*l.Lista[e.FnVincular]),
		Dependibles:    make(map[e.TipoDependible]*l.Lista[e.Dependible]),
	}

	if !strings.Contains(path, ".md") {
		return &a, nil
	}

	var contenido string
	if bytes, err := os.ReadFile(fmt.Sprintf("%s/%s", root.Path, path)); err != nil {
		return nil, fmt.Errorf("error al leer %s obteniendo el error: %v", path, err)
	} else {
		contenido = strings.TrimSpace(string(bytes))
	}

	if strings.Index(contenido, "---") != 0 {
		return &a, nil
	}

	blob := contenido[3 : 3+strings.Index(contenido[3:], "---")]
	decodificador := yaml.NewDecoder(strings.NewReader(blob))

	var meta Frontmatter
	if err := decodificador.Decode(&meta); err != nil {
		return nil, fmt.Errorf("error al decodificar en %s la metadata, con el error: %v", path, err)
	}

	// a.Contenido = contenido[3+strings.Index(contenido[3:], "---")+len("---"):]
	archivoCargable := e.NewArchivo(path, meta.Tags)

	a.Cargables.Push(archivoCargable)
	a.CargarDependible(e.DEP_ARCHIVO, archivoCargable)

	procesamiento := make(map[string]func(path string, meta *Frontmatter, canalMensajes chan string))

	// Facultad
	procesamiento[TAG_CARRERA] = a.ProcesarCarrera
	procesamiento[TAG_MATERIA] = a.ProcesarMateria
	procesamiento[TAG_MATERIA_EQUIVALENTE] = a.ProcesarMateriaEquivalente
	procesamiento[TAG_RESUMEN_MATERIA] = a.ProcesarTemaMateria

	// Cursos
	procesamiento[TAG_CURSO] = a.ProcesarCurso
	procesamiento[TAG_CURSO_PRESENCIA] = a.ProcesarCursoPresencial
	procesamiento[TAG_RESUMEN_CURSO] = a.ProcesarTemaCurso

	// Colecciones
	procesamiento[TAG_DISTRIBUCION] = a.ProcesarDistribucion
	procesamiento[TAG_LIBRO] = a.ProcesarLibro
	procesamiento[TAG_PAPER] = a.ProcesarPaper

	// Notas
	procesamiento[TAG_NOTA_FACULTAD] = a.ProcesarNota
	procesamiento[TAG_NOTA_CURSO] = a.ProcesarNota
	procesamiento[TAG_NOTA_INVESTIGACION] = a.ProcesarNota
	procesamiento[TAG_NOTA_COLECCION] = a.ProcesarNota
	procesamiento[TAG_NOTA_PROYECTO] = a.ProcesarNota

	for _, tag := range meta.Tags {
		if funcionProcesamiento, ok := procesamiento[tag]; ok {
			funcionProcesamiento(path, &meta, canalMensajes)
		}
	}

	info.MaxPath = max(info.MaxPath, uint32(len(path)))
	CargarInfo(info, &meta)

	return &a, nil
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
func (a *Archivo) EstablecerDependencias(canalDatos chan e.Cargable, canalDocumentos chan e.A, canalMensajes chan string) {
	for pathTipo, listaRelacion := range a.FnDependencias {
		if archivo, err := a.Root.EncontrarArchivo(pathTipo.Path); err != nil {
			canalMensajes <- fmt.Sprintf("No se encontró el archivo '%s' para el archivo: '%s'", pathTipo.Path, a.Path)

		} else if listaDependible, ok := archivo.Dependibles[pathTipo.Tipo]; !ok {
			canalMensajes <- fmt.Sprintf("No se encontró el tipo '%s' para el archivo: '%s'", e.TipoDependible2String(pathTipo.Tipo), a.Path)

		} else {
			for dependible := range listaDependible.Iterar {
				for relacion := range listaRelacion.Iterar {
					relacion(dependible)
				}
			}
		}
	}

	for cargable := range a.Cargables.Iterar {
		canalDatos <- cargable
	}
}

func (a *Archivo) ProcesarCarrera(path string, meta *Frontmatter, canalMensajes chan string) {
	nombreCarrera := e.Nombre(path)
	if constructor, err := e.NewCarrera(nombreCarrera, meta.Etapa, meta.TieneCodigo); err == nil {
		a.CargarDependencia(path, e.DEP_ARCHIVO, constructor.CrearDependenciaArchivo)

		a.CargarDependible(e.DEP_CARRERA, constructor)
	} else {
		canalMensajes <- fmt.Sprintf("Error: %v\n", err)
	}
}

func (a *Archivo) ProcesarMateria(path string, meta *Frontmatter, canalMensajes chan string) {
	if constructor, err := e.NewMateria(meta.NombreMateria, meta.Codigo, meta.Plan, meta.Cuatri, meta.Etapa); err == nil {
		pathCarrera := ObtenerWikiLink(meta.PathCarrera)[0]
		a.CargarDependencia(path, e.DEP_ARCHIVO, constructor.CrearDependenciaArchivo)
		a.CargarDependencia(pathCarrera, e.DEP_CARRERA, constructor.CrearDependenciaCarrera)

		a.CargarDependible(e.DEP_MATERIA, constructor)
	} else {
		canalMensajes <- fmt.Sprintf("Error: %v\n", err)
	}

	for _, correlativa := range meta.Correlativas {
		constructor := e.NewMateriasCorrelativas(e.MATERIA_REAL, correlativa.Tipo)

		a.CargarDependencia(path, e.DEP_MATERIA, constructor.CrearDependenciaMateria)
		switch correlativa.Tipo {
		case e.MATERIA_REAL:
			a.CargarDependencia(correlativa.Path, e.DEP_MATERIA, constructor.CrearDependenciaCorrelativa)
		case e.MATERIA_EQUIVALENTE:
			a.CargarDependencia(correlativa.Path, e.DEP_MATERIA_EQUIVALENTE, constructor.CrearDependenciaCorrelativa)
		}
	}
}

func (a *Archivo) ProcesarMateriaEquivalente(path string, meta *Frontmatter, canalMensajes chan string) {
	constructor := e.NewMateriaEquivalente(meta.NombreMateria, meta.Codigo)
	a.CargarDependencia(path, e.DEP_ARCHIVO, constructor.CrearDependenciaArchivo)
	a.CargarDependencia(ObtenerWikiLink(meta.Equivalencia)[0], e.DEP_MATERIA, constructor.CrearDependenciaMateria)
	a.CargarDependencia(ObtenerWikiLink(meta.PathCarrera)[0], e.DEP_CARRERA, constructor.CrearDependenciaCarrera)

	a.CargarDependible(e.DEP_MATERIA_EQUIVALENTE, constructor)

	for _, correlativa := range meta.Correlativas {
		constructor := e.NewMateriasCorrelativas(e.MATERIA_EQUIVALENTE, correlativa.Tipo)

		a.CargarDependencia(path, e.DEP_MATERIA_EQUIVALENTE, constructor.CrearDependenciaMateria)
		switch correlativa.Tipo {
		case e.MATERIA_REAL:
			a.CargarDependencia(correlativa.Path, e.DEP_MATERIA, constructor.CrearDependenciaCorrelativa)
		case e.MATERIA_EQUIVALENTE:
			a.CargarDependencia(correlativa.Path, e.DEP_MATERIA_EQUIVALENTE, constructor.CrearDependenciaCorrelativa)
		}
	}
}

func (a *Archivo) ProcesarTemaMateria(path string, meta *Frontmatter, canalMensajes chan string) {
	constructor := e.NewTemaMateria(meta.NombreResumen, meta.Capitulo, meta.Parte)
	a.CargarDependencia(path, e.DEP_ARCHIVO, constructor.CrearDependenciaArchivo)
	a.CargarDependencia(meta.MateriaResumen, e.DEP_MATERIA, constructor.CrearDependenciaMateria)

	a.CargarDependible(e.DEP_TEMA_MATERIA, constructor)
}

func (a *Archivo) ProcesarCurso(path string, meta *Frontmatter, canalMensajes chan string) {
	if constructor, err := meta.CrearCurso(); err == nil {
		a.CargarDependencia(path, e.DEP_ARCHIVO, constructor.CrearDependenciaArchivo)

		a.CargarDependible(e.DEP_CURSO, constructor)
	} else {
		canalMensajes <- fmt.Sprintf("Error: %v\n", err)
	}
}

func (a *Archivo) ProcesarCursoPresencial(path string, meta *Frontmatter, canalMensajes chan string) {
	if constructor, err := meta.CrearCursoPresencial(); err == nil {
		a.CargarDependencia(path, e.DEP_ARCHIVO, constructor.CrearDependenciaArchivo)

		a.CargarDependible(e.DEP_CURSO_PRESENCIAL, constructor)
	} else {
		canalMensajes <- fmt.Sprintf("Error: %v\n", err)
	}
}

func (a *Archivo) ProcesarTemaCurso(path string, meta *Frontmatter, canalMensajes chan string) {
	constructor := e.NewTemaCurso(meta.NombreResumen, meta.Capitulo, meta.Parte, meta.TipoCurso)
	a.CargarDependencia(path, e.DEP_ARCHIVO, constructor.CrearDependenciaArchivo)

	pathCurso := ObtenerWikiLink(meta.Curso)[0]
	switch meta.TipoCurso {
	case e.CURSO_ONLINE:
		a.CargarDependencia(pathCurso, e.DEP_CURSO, constructor.CrearDependenciaCurso)
	case e.CURSO_PRESENCIAL:
		a.CargarDependencia(pathCurso, e.DEP_CURSO_PRESENCIAL, constructor.CrearDependenciaCurso)
	}

	a.CargarDependible(e.DEP_TEMA_CURSO, constructor)
}

func (a *Archivo) ProcesarNota(path string, meta *Frontmatter, canalMensajes chan string) {
	constructor := meta.CrearNota(e.Nombre(path))
	a.CargarDependencia(path, e.DEP_ARCHIVO, constructor.CrearDependenciaArchivo)
	a.CargarDependible(e.DEP_NOTA, constructor)

	vinculosNota := [][]string{meta.VinculoFacultad, meta.VinculoCurso}
	tipoNota := []e.TipoNota{e.TN_FACULTAD, e.TN_CURSO}
	tipoDependencia := []e.TipoDependible{e.DEP_TEMA_MATERIA, e.DEP_TEMA_CURSO}

	for i, vinculos := range vinculosNota {
		for _, vinculo := range vinculos {
			pathVinculo := ObtenerWikiLink(vinculo)[0]

			notaVinculo := e.NewNotaVinculo(tipoNota[i])
			a.CargarDependencia(path, e.DEP_NOTA, notaVinculo.CrearDependenciaNota)
			a.CargarDependencia(pathVinculo, tipoDependencia[i], notaVinculo.CrearDependenciaVinculo)
		}
	}
}

func (a *Archivo) ProcesarLibro(path string, meta *Frontmatter, canalMensajes chan string) {
	constructor := meta.CrearLibro()
	a.CargarDependencia(path, e.DEP_ARCHIVO, constructor.CrearDependenciaArchivo)
}

func (a *Archivo) ProcesarPaper(path string, meta *Frontmatter, canalMensajes chan string) {
	if constructor, err := meta.CrearPaper(); err == nil {
		a.CargarDependencia(path, e.DEP_ARCHIVO, constructor.CrearDependenciaArchivo)

	} else {
		canalMensajes <- fmt.Sprintf("Error: %v\n", err)
	}
}

func (a *Archivo) ProcesarDistribucion(path string, meta *Frontmatter, canalMensajes chan string) {
	if constructor, err := e.NewDistribucion(meta.NombreDistribuucion, meta.TipoDistribucion); err == nil {
		a.CargarDependencia(path, e.DEP_ARCHIVO, constructor.CrearDependenciaArchivo)

	} else {
		canalMensajes <- fmt.Sprintf("Error: %v\n", err)
	}
}

func ObtenerWikiLink(link string) []string {
	link = strings.TrimPrefix(link, "[[")
	link = strings.TrimSuffix(link, "]]")
	return strings.Split(link, "|")
}
