package fs

import (
	"fmt"
	"os"
	e "own_wiki/system_protocol/datos"
	d "own_wiki/system_protocol/dependencias"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
)

const (
	TABLA_ARCHIVOS    = "Archivos"
	TABLA_TAGS        = "Tags"
	TABLA_PERSONAS    = "Personas"
	TABLA_EDITORIALES = "Editoriales"

	TABLA_PLANES       = "PlanesCarrera"
	TABLA_CUATRI       = "CuatrimestresCarrera"
	TABLA_CARRERAS     = "Carreras"
	TABLA_MATERIAS     = "Materias"
	TABLA_MATERIAS_EQ  = "MateriasEquivalentes"
	TABLA_TEMA_MATERIA = "TemasMateria"

	TABLA_COLECCIONES               = "Colecciones"
	TABLA_LIBROS                    = "Libros"
	TABLA_AUTORES_LIBRO             = "AutoresLibro"
	TABLA_CAPITULOS                 = "Capitulos"
	TABLA_EDITORES_CAPITULO         = "EditoresCapitulo"
	TABLA_DISTRIBUCIONES            = "Distribuciones"
	TABLA_RELACIONES_DISTRIBUCIONES = "RelacionesDistribuciones"
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
	TAG_REPRESENTANTE = "colección/representante"
	TAG_DISTRIBUCION  = "colección/distribuciones/distribución"
	TAG_LIBRO         = "colección/biblioteca/libro"
	TAG_PAPER         = "colección/biblioteca/paper"
)

const (
	TAG_NOTA_FACULTAD      = "nota/facultad"
	TAG_NOTA_CURSO         = "nota/curso"
	TAG_NOTA_INVESTIGACION = "nota/investigacion"
	TAG_NOTA_COLECCION     = "nota/colección"
	TAG_NOTA_PROYECTO      = "nota/proyecto"
)

const HABILITAR_ERROR = true

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

func CargarArchivo(dirInicio string, path string, tracker *d.TrackerDependencias, canalMensajes chan string) error {
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

	err := tracker.Cargar(TABLA_ARCHIVOS, []d.ForeignKey{}, path)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar archivo con error: %v", err)
	}

	funcionesProcesar := make(map[string]func(path string, meta *Frontmatter, tracker *d.TrackerDependencias) error)

	// Carrera
	funcionesProcesar[TAG_CARRERA] = ProcesarCarrera
	funcionesProcesar[TAG_MATERIA] = ProcesarMateria
	funcionesProcesar[TAG_MATERIA_EQUIVALENTE] = ProcesarMateriaEquivalente
	funcionesProcesar[TAG_RESUMEN_MATERIA] = ProcesarTemaMateria

	// Cursos:
	// funcionesProcesar[TAG_CURSO] =
	// funcionesProcesar[TAG_CURSO_PRESENCIA] =
	// funcionesProcesar[TAG_RESUMEN_CURSO] =

	// Colecciones:
	funcionesProcesar[TAG_REPRESENTANTE] = ProcesarColeccion
	funcionesProcesar[TAG_DISTRIBUCION] = ProcesarDistribucion
	funcionesProcesar[TAG_LIBRO] = ProcesarLibro
	// funcionesProcesar[TAG_PAPER] =

	// Notas:
	// funcionesProcesar[TAG_NOTA_FACULTAD] =
	// funcionesProcesar[TAG_NOTA_CURSO] =
	// funcionesProcesar[TAG_NOTA_INVESTIGACION] =
	// funcionesProcesar[TAG_NOTA_COLECCION] =
	// funcionesProcesar[TAG_NOTA_PROYECTO] =

	for _, tag := range meta.Tags {
		funcionProcesar, ok := funcionesProcesar[tag]
		if !ok {
			continue
		}

		if err = funcionProcesar(path, &meta, tracker); err != nil {
			return err
		}
	}

	return nil
}

func ProcesarCarrera(path string, meta *Frontmatter, tracker *d.TrackerDependencias) error {
	nombreCarrera := Nombre(path)
	err := tracker.Cargar(TABLA_CARRERAS,
		[]d.ForeignKey{
			tracker.CrearReferencia(TABLA_ARCHIVOS, "refArchivo", []d.ForeignKey{}, path),
		},
		nombreCarrera,
		EtapaODefault(meta.Etapa, e.ETAPA_SIN_EMPEZAR),
		strings.ToLower(strings.TrimSpace(meta.TieneCodigo)) == "true",
	)

	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar carrera con error: %v", err)
	}
	return nil
}

func ProcesarMateria(path string, meta *Frontmatter, tracker *d.TrackerDependencias) error {
	err := tracker.Cargar(TABLA_PLANES, []d.ForeignKey{}, meta.Plan)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar plan con error: %v", err)
	}

	anio, cuatrimestre, err := ObtenerCuatrimestreParte(meta.Cuatri)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("obtener cuatrimestre con error: %v", err)
	}

	err = tracker.Cargar(TABLA_CUATRI, []d.ForeignKey{}, anio, cuatrimestre)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar cuatri con error: %v", err)
	}

	err = tracker.Cargar(TABLA_MATERIAS,
		[]d.ForeignKey{
			tracker.CrearReferencia(TABLA_ARCHIVOS, "refArchivo", []d.ForeignKey{}, path),
			tracker.CrearReferencia(TABLA_CARRERAS, "refCarrera", []d.ForeignKey{}, meta.NombreCarrera),
			tracker.CrearReferencia(TABLA_PLANES, "refPlan", []d.ForeignKey{}, meta.Plan),
			tracker.CrearReferencia(TABLA_CUATRI, "refCuatrimestre", []d.ForeignKey{}, anio, cuatrimestre),
		},
		meta.NombreMateria,
		EtapaODefault(meta.Etapa, e.ETAPA_SIN_EMPEZAR),
		meta.Codigo,
	)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar materia con error: %v", err)
	}
	return nil
}

func ProcesarMateriaEquivalente(path string, meta *Frontmatter, tracker *d.TrackerDependencias) error {
	infoMateria := meta.MateriaEquivalente
	err := tracker.Cargar(TABLA_MATERIAS_EQ,
		[]d.ForeignKey{
			tracker.CrearReferencia(TABLA_ARCHIVOS, "refArchivo", []d.ForeignKey{}, path),
			tracker.CrearReferencia(TABLA_CARRERAS, "refCarrera", []d.ForeignKey{}, meta.NombreCarrera),
			tracker.CrearReferencia(TABLA_MATERIAS, "refMateria", []d.ForeignKey{
				tracker.CrearReferencia(TABLA_CARRERAS, "refCarrera", []d.ForeignKey{}, infoMateria.Carrera),
			}, infoMateria.NombreMateria),
		},
		meta.NombreMateria,
		meta.Codigo,
	)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar materia equivalente con error: %v", err)
	}
	return nil
}

func ProcesarTemaMateria(path string, meta *Frontmatter, tracker *d.TrackerDependencias) error {
	infoTema := meta.InfoTemaMateria
	err := tracker.Cargar(TABLA_TEMA_MATERIA,
		[]d.ForeignKey{
			tracker.CrearReferencia(TABLA_ARCHIVOS, "refArchivo", []d.ForeignKey{}, path),
			tracker.CrearReferencia(TABLA_MATERIAS, "refMateria", []d.ForeignKey{
				tracker.CrearReferencia(TABLA_CARRERAS, "refCarrera", []d.ForeignKey{}, infoTema.Carrera),
			}, infoTema.Materia),
		},
		meta.NombreResumen,
		NumeroODefault(meta.Capitulo, 1),
		NumeroODefault(meta.Parte, 0),
	)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar tema materia con error: %v", err)
	}
	return nil
}

func ProcesarColeccion(path string, meta *Frontmatter, tracker *d.TrackerDependencias) error {
	err := tracker.Cargar(TABLA_COLECCIONES,
		[]d.ForeignKey{tracker.CrearReferencia(TABLA_ARCHIVOS, "refArchivo", []d.ForeignKey{}, path)},
		Nombre(path),
	)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar colecciones con error: %v", err)
	}
	return nil
}

func ProcesarDistribucion(path string, meta *Frontmatter, tracker *d.TrackerDependencias) error {
	tipoDistribucion, err := ObtenerTipoDistribucion(meta.TipoDistribucion)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("obtener tipo distribucion con error: %v", err)
	}

	err = tracker.Cargar(TABLA_DISTRIBUCIONES,
		[]d.ForeignKey{
			tracker.CrearReferencia(TABLA_COLECCIONES, "refColeccion", []d.ForeignKey{}, "Distribuciones"),
			tracker.CrearReferencia(TABLA_ARCHIVOS, "refArchivo", []d.ForeignKey{}, path),
		},
		meta.NombreDistribuucion,
		tipoDistribucion,
	)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar distribuciones con error: %v", err)
	}
	return nil
}

func ProcesarLibro(path string, meta *Frontmatter, tracker *d.TrackerDependencias) error {
	err := tracker.Cargar(TABLA_EDITORIALES, []d.ForeignKey{}, meta.Editorial)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar editoriales con error: %v", err)
	}

	anio := NumeroODefault(meta.Anio, 0)
	edicion := NumeroODefault(meta.Edicion, 1)
	volumen := NumeroODefault(meta.Volumen, 0)

	err = tracker.Cargar(
		TABLA_LIBROS, []d.ForeignKey{
			tracker.CrearReferencia(TABLA_ARCHIVOS, "refArchivo", []d.ForeignKey{}, path),
			tracker.CrearReferencia(TABLA_EDITORIALES, "refEditorial", []d.ForeignKey{}, meta.Editorial),
			tracker.CrearReferencia(TABLA_COLECCIONES, "refColeccion", []d.ForeignKey{}, "Biblioteca"),
		},
		meta.TituloObra,
		meta.SubtituloObra,
		anio,
		edicion,
		volumen,
		meta.Url,
	)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar libro con error: %v", err)
	}

	for _, autor := range meta.NombreAutores {
		nombre := strings.TrimSpace(autor.Nombre)
		apellido := strings.TrimSpace(autor.Apellido)

		err = tracker.Cargar(TABLA_PERSONAS, []d.ForeignKey{}, nombre, apellido)
		if HABILITAR_ERROR && err != nil {
			return fmt.Errorf("cargar persona con error: %v", err)
		}

		err = tracker.Cargar(TABLA_AUTORES_LIBRO, []d.ForeignKey{
			tracker.CrearReferencia(TABLA_LIBROS, "refLibro", []d.ForeignKey{},
				meta.TituloObra,
				anio,
				edicion,
				volumen,
			),
			tracker.CrearReferencia(TABLA_PERSONAS, "refPersona", []d.ForeignKey{},
				nombre,
				apellido,
			),
		})
		if HABILITAR_ERROR && err != nil {
			return fmt.Errorf("cargar autor libro con error: %v", err)
		}
	}

	for _, capitulo := range meta.Capitulos {
		numero := NumeroODefault(capitulo.NumeroCapitulo, 0)
		paginaInicio := NumeroODefault(capitulo.Paginas.Inicio, 0)
		paginaFinal := NumeroODefault(capitulo.Paginas.Final, 1)

		err = tracker.Cargar(TABLA_CAPITULOS,
			[]d.ForeignKey{
				tracker.CrearReferencia(TABLA_ARCHIVOS, "refArchivo", []d.ForeignKey{}, path),
				tracker.CrearReferencia(TABLA_LIBROS, "refLibro", []d.ForeignKey{},
					meta.TituloObra,
					anio,
					edicion,
					volumen,
				),
			},
			numero,
			capitulo.NombreCapitulo,
			paginaInicio,
			paginaFinal,
		)
		if HABILITAR_ERROR && err != nil {
			return fmt.Errorf("cargar capitulo con error: %v", err)
		}

		for _, editor := range capitulo.Editores {
			nombre := strings.TrimSpace(editor.Nombre)
			apellido := strings.TrimSpace(editor.Apellido)

			err = tracker.Cargar(TABLA_PERSONAS, []d.ForeignKey{}, nombre, apellido)
			if HABILITAR_ERROR && err != nil {
				return fmt.Errorf("cargar persona con error: %v", err)
			}

			err = tracker.Cargar(TABLA_EDITORES_CAPITULO, []d.ForeignKey{
				tracker.CrearReferencia(TABLA_CAPITULOS, "refCapitulo", []d.ForeignKey{},
					numero,
					capitulo.NombreCapitulo,
					paginaInicio,
					paginaFinal,
				),
				tracker.CrearReferencia(TABLA_PERSONAS, "refPersona", []d.ForeignKey{},
					nombre,
					apellido,
				),
			})
			if HABILITAR_ERROR && err != nil {
				return fmt.Errorf("cargar editor capitulo con error: %v", err)
			}
		}
	}

	return nil
}

/*
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


	func (a *Archivo) ProcesarPaper(path string, meta *Frontmatter, canalMensajes chan string) {
		if constructor, err := meta.CrearPaper(); err == nil {
			a.CargarDependencia(path, e.DEP_ARCHIVO, constructor.CrearDependenciaArchivo)

		} else {
			canalMensajes <- fmt.Sprintf("Error: %v\n", err)
		}
	}

*/

func Nombre(path string) string {
	separacion := strings.Split(path, "/")
	return strings.ReplaceAll(separacion[len(separacion)-1], ".md", "")
}

func ObtenerWikiLink(link string) []string {
	link = strings.TrimPrefix(link, "[[")
	link = strings.TrimSuffix(link, "]]")
	return strings.Split(link, "|")
}

type Etapa string

const (
	ETAPA_SIN_EMPEZAR = "SinEmpezar"
	ETAPA_EMPEZADO    = "Empezado"
	ETAPA_AMPLIAR     = "Ampliar"
	ETAPA_TERMINADO   = "Terminado"
)

type TipoDistribucion string

const (
	DISTRIBUCION_DISCRETA     = "Discreta"
	DISTRIBUCION_CONTINUA     = "Continua"
	DISTRIBUCION_MULTIVARIADA = "Multivariada"
)

type ParteCuatrimestre string

const (
	CUATRIMESTRE_PRIMERO = "Primero"
	CUATRIMESTRE_SEGUNDO = "Segundo"
)

func NumeroODefault(representacion string, valorDefault int) int {
	if nuevoValor, err := strconv.Atoi(representacion); err == nil {
		return nuevoValor
	} else {
		return valorDefault
	}
}

func BooleanoODefault(representacion string, valorDefault bool) bool {
	switch representacion {
	case "true":
		return true
	case "false":
		return false
	default:
		return valorDefault
	}
}

func EtapaODefault(representacion string, valorDefault Etapa) Etapa {
	var etapa Etapa
	switch representacion {
	case "sin-empezar":
		etapa = ETAPA_SIN_EMPEZAR
	case "empezado":
		etapa = ETAPA_EMPEZADO
	case "ampliar":
		etapa = ETAPA_AMPLIAR
	case "terminado":
		etapa = ETAPA_TERMINADO
	default:
		etapa = valorDefault
	}
	return etapa
}

func ObtenerEtapa(representacionEtapa string) (Etapa, error) {
	var etapa Etapa
	switch representacionEtapa {
	case "sin-empezar":
		etapa = ETAPA_SIN_EMPEZAR
	case "empezado":
		etapa = ETAPA_EMPEZADO
	case "ampliar":
		etapa = ETAPA_AMPLIAR
	case "terminado":
		etapa = ETAPA_TERMINADO
	default:
		return ETAPA_SIN_EMPEZAR, fmt.Errorf("el tipo de etapa (%s) no es uno de los esperados", representacionEtapa)
	}

	return etapa, nil
}

func ObtenerTipoDistribucion(representacion string) (TipoDistribucion, error) {
	var tipoDistribucion TipoDistribucion
	switch representacion {
	case "discreta":
		tipoDistribucion = DISTRIBUCION_DISCRETA
	case "continua":
		tipoDistribucion = DISTRIBUCION_CONTINUA
	case "multivariada":
		tipoDistribucion = DISTRIBUCION_MULTIVARIADA
	default:
		return DISTRIBUCION_DISCRETA, fmt.Errorf("el tipo de distribucion (%s) no es uno de los esperados", representacion)
	}

	return tipoDistribucion, nil
}

func ObtenerCuatrimestreParte(representacionCuatri string) (int, ParteCuatrimestre, error) {
	var anio int
	var cuatriNum int
	var cuatri ParteCuatrimestre

	if _, err := fmt.Sscanf(representacionCuatri, "%dC%d", &anio, &cuatriNum); err != nil {
		return anio, cuatri, fmt.Errorf("el tipo de anio-cuatri (%s) no es uno de los esperados", representacionCuatri)
	}

	switch cuatriNum {
	case 1:
		cuatri = CUATRIMESTRE_PRIMERO
	case 2:
		cuatri = CUATRIMESTRE_SEGUNDO
	default:
		return anio, cuatri, fmt.Errorf("el cuatri dado por %d no es posible representar", cuatriNum)
	}

	return anio, cuatri, nil
}
