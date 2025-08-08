package fs

import (
	"fmt"
	"os"
	e "own_wiki/system_protocol/datos"
	t "own_wiki/system_protocol/tablas"
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

func CargarArchivo(dirInicio string, path string, tablas *t.Tablas, canalMensajes chan string) error {
	// Tal vez cambiarlo despues por la existencia de imagenes
	if !strings.Contains(path, ".md") {
		return nil
	}

	var contenido string
	if bytes, err := os.ReadFile(fmt.Sprintf("%s/%s", dirInicio, path)); err != nil {
		return fmt.Errorf("error al leer %s obteniendo el error: %v", path, err)
	} else {
		contenido = strings.TrimSpace(string(bytes))
	}

	if strings.Index(contenido, "---") != 0 {
		return nil
	}

	blob := contenido[3 : 3+strings.Index(contenido[3:], "---")]
	decodificador := yaml.NewDecoder(strings.NewReader(blob))

	var meta Frontmatter
	if err := decodificador.Decode(&meta); err != nil {
		return fmt.Errorf("error al decodificar en %s la metadata, con el error: %v", path, err)
	}

	// El resto del contenido del archivo
	// a.Contenido = contenido[3+strings.Index(contenido[3:], "---")+len("---"):]

	if err := tablas.Archivos.CargarArchivo(path); err != nil {
		return err
	}

	for _, tag := range meta.Tags {
		if err := tablas.Tags.CargarTag(path, tag); err != nil {
			return err
		}

		switch tag {
		case TAG_CARRERA:

		case TAG_MATERIA:

		case TAG_MATERIA_EQUIVALENTE:

		case TAG_RESUMEN_MATERIA:

		case TAG_CURSO:

		case TAG_CURSO_PRESENCIA:

		case TAG_RESUMEN_CURSO:

		case TAG_DISTRIBUCION:

		case TAG_LIBRO:

		case TAG_PAPER:

		case TAG_NOTA_FACULTAD:
			fallthrough
		case TAG_NOTA_CURSO:
			fallthrough
		case TAG_NOTA_INVESTIGACION:
			fallthrough
		case TAG_NOTA_COLECCION:
			fallthrough
		case TAG_NOTA_PROYECTO:

		}
	}

	return nil
}

/*

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
*/

func ObtenerWikiLink(link string) []string {
	link = strings.TrimPrefix(link, "[[")
	link = strings.TrimSuffix(link, "]]")
	return strings.Split(link, "|")
}
