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

	TABLA_PAGINAS_CURSOS     = "PaginasCursos"
	TABLA_CURSOS_ONLINE      = "CursosOnline"
	TABLA_CURSOS_PRESENECIAL = "CursosPresencial"
	TABLA_TEMA_CURSO         = "TemasCurso"
	TABLA_PROFESORES_CURSO   = "ProfesoresCurso"

	TABLA_PLANES       = "PlanesCarrera"
	TABLA_CUATRI       = "CuatrimestresCarrera"
	TABLA_CARRERAS     = "Carreras"
	TABLA_MATERIAS     = "Materias"
	TABLA_MATERIAS_EQ  = "MateriasEquivalentes"
	TABLA_TEMA_MATERIA = "TemasMateria"
	TABLA_CORRELATIVAS = "MateriasCorrelativas"

	TABLA_COLECCIONES               = "Colecciones"
	TABLA_LIBROS                    = "Libros"
	TABLA_AUTORES_LIBRO             = "AutoresLibro"
	TABLA_CAPITULOS                 = "Capitulos"
	TABLA_EDITORES_CAPITULO         = "EditoresCapitulo"
	TABLA_DISTRIBUCIONES            = "Distribuciones"
	TABLA_RELACIONES_DISTRIBUCIONES = "RelacionesDistribuciones"
	TABLA_REVISTAS_PAPER            = "RevistasDePaper"
	TABLA_PAPERS                    = "Papers"
	TABLA_ESCRITORES_PAPER          = "EscritoresPaper"
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

	err := tracker.Cargar(TABLA_ARCHIVOS, []d.RelacionTabla{}, path)
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
	funcionesProcesar[TAG_CURSO] = ProcesarCursoOnline
	funcionesProcesar[TAG_CURSO_PRESENCIA] = ProcesarCursoPresencial
	funcionesProcesar[TAG_RESUMEN_CURSO] = ProcesarTemaCurso

	// Colecciones:
	funcionesProcesar[TAG_REPRESENTANTE] = ProcesarColeccion
	funcionesProcesar[TAG_DISTRIBUCION] = ProcesarDistribucion
	funcionesProcesar[TAG_LIBRO] = ProcesarLibro
	funcionesProcesar[TAG_PAPER] = ProcesarPaper

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

func ProcesarCursoOnline(path string, meta *Frontmatter, tracker *d.TrackerDependencias) error {
	err := tracker.Cargar(TABLA_PAGINAS_CURSOS, []d.RelacionTabla{}, meta.NombrePagina)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar paginas cursos con error: %v", err)
	}

	anio, err := strconv.Atoi(meta.FechaCurso)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("obtener anio del curso online con error: %v", err)
	}

	err = tracker.Cargar(TABLA_CURSOS_ONLINE,
		[]d.RelacionTabla{
			d.NewRelacionSimple(TABLA_ARCHIVOS, "refArchivo", path),
			d.NewRelacionSimple(TABLA_PAGINAS_CURSOS, "refPagina", meta.NombrePagina),
		},
		meta.NombreCurso,
		EtapaODefault(meta.Etapa, e.ETAPA_SIN_EMPEZAR),
		anio,
		meta.Url,
	)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar curso online con error: %v", err)
	}

	for _, profesor := range meta.NombreAutores {
		nombre := strings.TrimSpace(profesor.Nombre)
		apellido := strings.TrimSpace(profesor.Apellido)

		err = tracker.Cargar(TABLA_PERSONAS, []d.RelacionTabla{}, nombre, apellido)
		if HABILITAR_ERROR && err != nil {
			return fmt.Errorf("cargar persona en curso online con error: %v", err)
		}

		err = tracker.Cargar(TABLA_PROFESORES_CURSO,
			[]d.RelacionTabla{
				d.NewRelacionSimple(TABLA_CURSOS_ONLINE, "refCurso",
					meta.NombreCurso,
					anio,
				),
				d.NewRelacionSimple(TABLA_PERSONAS, "refPersona",
					nombre,
					apellido,
				),
			},
			CURSO_ONLINE,
		)
		if HABILITAR_ERROR && err != nil {
			return fmt.Errorf("cargar profesor en el curso online con error: %v", err)
		}
	}

	return nil
}

func ProcesarCursoPresencial(path string, meta *Frontmatter, tracker *d.TrackerDependencias) error {
	anio, err := strconv.Atoi(meta.FechaCurso)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("obtener anio del curso presencial con error: %v", err)
	}

	err = tracker.Cargar(TABLA_CURSOS_PRESENECIAL,
		[]d.RelacionTabla{
			d.NewRelacionSimple(TABLA_ARCHIVOS, "refArchivo", path),
		},
		meta.NombreCurso,
		EtapaODefault(meta.Etapa, e.ETAPA_SIN_EMPEZAR),
		anio,
	)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar curso presencial con error: %v", err)
	}

	for _, profesor := range meta.NombreAutores {
		nombre := strings.TrimSpace(profesor.Nombre)
		apellido := strings.TrimSpace(profesor.Apellido)

		err = tracker.Cargar(TABLA_PERSONAS, []d.RelacionTabla{}, nombre, apellido)
		if HABILITAR_ERROR && err != nil {
			return fmt.Errorf("cargar persona en curso presencial con error: %v", err)
		}

		err = tracker.Cargar(TABLA_PROFESORES_CURSO,
			[]d.RelacionTabla{
				d.NewRelacionSimple(TABLA_CURSOS_PRESENECIAL, "refCurso",
					meta.NombreCurso,
					anio,
				),
				d.NewRelacionSimple(TABLA_PERSONAS, "refPersona",
					nombre,
					apellido,
				),
			},
			CURSO_PRESENCIAL,
		)
		if HABILITAR_ERROR && err != nil {
			return fmt.Errorf("cargar profesor en el curso presencial con error: %v", err)
		}
	}

	return nil
}

func ProcesarTemaCurso(path string, meta *Frontmatter, tracker *d.TrackerDependencias) error {
	infoTema := meta.InfoCurso
	tablaCurso := TABLA_CURSOS_ONLINE
	if infoTema.Tipo == CURSO_PRESENCIAL {
		tablaCurso = TABLA_CURSOS_PRESENECIAL
	}

	anio, err := strconv.Atoi(infoTema.Anio)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("obtener anio del tema del curso con error: %v", err)
	}

	err = tracker.Cargar(TABLA_TEMA_CURSO,
		[]d.RelacionTabla{
			d.NewRelacionSimple(TABLA_ARCHIVOS, "refArchivo", path),
			d.NewRelacionSimple(tablaCurso, "refCurso", infoTema.Curso, anio),
		},
		meta.NombreResumen,
		infoTema.Tipo,
		NumeroODefault(meta.Capitulo, 1),
		NumeroODefault(meta.Parte, 0),
	)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar tema curso con error: %v", err)
	}
	return nil
}

func ProcesarCarrera(path string, meta *Frontmatter, tracker *d.TrackerDependencias) error {
	nombreCarrera := Nombre(path)
	err := tracker.Cargar(TABLA_CARRERAS,
		[]d.RelacionTabla{
			d.NewRelacionSimple(TABLA_ARCHIVOS, "refArchivo", path),
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
	err := tracker.Cargar(TABLA_PLANES, []d.RelacionTabla{}, meta.Plan)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar plan con error: %v", err)
	}

	anio, cuatrimestre, err := ObtenerCuatrimestreParte(meta.Cuatri)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("obtener cuatrimestre con error: %v", err)
	}

	err = tracker.Cargar(TABLA_CUATRI, []d.RelacionTabla{}, anio, cuatrimestre)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar cuatri con error: %v", err)
	}

	err = tracker.Cargar(TABLA_MATERIAS,
		[]d.RelacionTabla{
			d.NewRelacionSimple(TABLA_ARCHIVOS, "refArchivo", path),
			d.NewRelacionSimple(TABLA_CARRERAS, "refCarrera", meta.NombreCarrera),
			d.NewRelacionSimple(TABLA_PLANES, "refPlan", meta.Plan),
			d.NewRelacionSimple(TABLA_CUATRI, "refCuatrimestre", anio, cuatrimestre),
		},
		meta.NombreMateria,
		EtapaODefault(meta.Etapa, e.ETAPA_SIN_EMPEZAR),
		meta.Codigo,
	)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar materia con error: %v", err)
	}

	for _, infoCorrelativa := range meta.Correlativas {
		tablaCorrelativa := TABLA_MATERIAS
		if infoCorrelativa.Tipo == MATERIA_EQUIVALENTE {
			tablaCorrelativa = TABLA_MATERIAS_EQ
		}

		err = tracker.Cargar(TABLA_CORRELATIVAS,
			[]d.RelacionTabla{
				d.NewRelacionCompleja(TABLA_MATERIAS, "refMateria", []d.RelacionTabla{
					d.NewRelacionSimple(TABLA_CARRERAS, "refCarrera", meta.NombreCarrera),
				}, meta.NombreMateria),
				d.NewRelacionCompleja(tablaCorrelativa, "refCorrelativa", []d.RelacionTabla{
					d.NewRelacionSimple(TABLA_CARRERAS, "refCarrera", meta.NombreCarrera),
				}, infoCorrelativa.Materia),
			},
			MATERIA_REAL,
			infoCorrelativa.Tipo,
		)
		if HABILITAR_ERROR && err != nil {
			return fmt.Errorf("cargar materias correlativas de una materia normal con error: %v", err)
		}
	}

	return nil
}

func ProcesarMateriaEquivalente(path string, meta *Frontmatter, tracker *d.TrackerDependencias) error {
	infoMateria := meta.MateriaEquivalente
	relacionCarrera := d.NewRelacionSimple(TABLA_CARRERAS, "refCarrera", meta.NombreCarrera)

	err := tracker.Cargar(TABLA_MATERIAS_EQ,
		[]d.RelacionTabla{
			d.NewRelacionSimple(TABLA_ARCHIVOS, "refArchivo", path),
			relacionCarrera,
			d.NewRelacionCompleja(TABLA_MATERIAS, "refMateria", []d.RelacionTabla{
				d.NewRelacionSimple(TABLA_CARRERAS, "refCarrera", infoMateria.Carrera),
			}, infoMateria.NombreMateria),
		},
		meta.NombreMateria,
		meta.Codigo,
	)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar materia equivalente con error: %v", err)
	}

	for _, infoCorrelativa := range meta.Correlativas {
		tablaCorrelativa := TABLA_MATERIAS
		if infoCorrelativa.Tipo == MATERIA_EQUIVALENTE {
			tablaCorrelativa = TABLA_MATERIAS_EQ
		}

		err = tracker.Cargar(TABLA_CORRELATIVAS,
			[]d.RelacionTabla{
				d.NewRelacionCompleja(TABLA_MATERIAS_EQ, "refMateria", []d.RelacionTabla{relacionCarrera}, meta.NombreMateria),
				d.NewRelacionCompleja(tablaCorrelativa, "refCorrelativa", []d.RelacionTabla{relacionCarrera}, infoCorrelativa.Materia),
			},
			MATERIA_EQUIVALENTE,
			infoCorrelativa.Tipo,
		)
		if HABILITAR_ERROR && err != nil {
			return fmt.Errorf("cargar materias correlativas de una materia equivalente con error: %v", err)
		}
	}

	return nil
}

func ProcesarTemaMateria(path string, meta *Frontmatter, tracker *d.TrackerDependencias) error {
	infoTema := meta.InfoTemaMateria
	err := tracker.Cargar(TABLA_TEMA_MATERIA,
		[]d.RelacionTabla{
			d.NewRelacionSimple(TABLA_ARCHIVOS, "refArchivo", path),
			d.NewRelacionCompleja(TABLA_MATERIAS, "refMateria", []d.RelacionTabla{
				d.NewRelacionSimple(TABLA_CARRERAS, "refCarrera", infoTema.Carrera),
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
		[]d.RelacionTabla{d.NewRelacionSimple(TABLA_ARCHIVOS, "refArchivo", path)},
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
		[]d.RelacionTabla{
			d.NewRelacionSimple(TABLA_COLECCIONES, "refColeccion", "Distribuciones"),
			d.NewRelacionSimple(TABLA_ARCHIVOS, "refArchivo", path),
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
	err := tracker.Cargar(TABLA_EDITORIALES, []d.RelacionTabla{}, meta.Editorial)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar editoriales con error: %v", err)
	}

	anio, err := strconv.Atoi(meta.Anio)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("obtener anio del libro con error: %v", err)
	}

	edicion := NumeroODefault(meta.Edicion, 1)
	volumen := NumeroODefault(meta.Volumen, 0)

	err = tracker.Cargar(
		TABLA_LIBROS, []d.RelacionTabla{
			d.NewRelacionSimple(TABLA_ARCHIVOS, "refArchivo", path),
			d.NewRelacionSimple(TABLA_EDITORIALES, "refEditorial", meta.Editorial),
			d.NewRelacionSimple(TABLA_COLECCIONES, "refColeccion", "Biblioteca"),
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

		err = tracker.Cargar(TABLA_PERSONAS, []d.RelacionTabla{}, nombre, apellido)
		if HABILITAR_ERROR && err != nil {
			return fmt.Errorf("cargar persona con error: %v", err)
		}

		err = tracker.Cargar(TABLA_AUTORES_LIBRO, []d.RelacionTabla{
			d.NewRelacionSimple(TABLA_LIBROS, "refLibro",
				meta.TituloObra,
				anio,
				edicion,
				volumen,
			),
			d.NewRelacionSimple(TABLA_PERSONAS, "refPersona",
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
			[]d.RelacionTabla{
				d.NewRelacionSimple(TABLA_ARCHIVOS, "refArchivo", path),
				d.NewRelacionSimple(TABLA_LIBROS, "refLibro",
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

			err = tracker.Cargar(TABLA_PERSONAS, []d.RelacionTabla{}, nombre, apellido)
			if HABILITAR_ERROR && err != nil {
				return fmt.Errorf("cargar persona con error: %v", err)
			}

			err = tracker.Cargar(TABLA_EDITORES_CAPITULO, []d.RelacionTabla{
				d.NewRelacionSimple(TABLA_CAPITULOS, "refCapitulo",
					numero,
					capitulo.NombreCapitulo,
					paginaInicio,
					paginaFinal,
				),
				d.NewRelacionSimple(TABLA_PERSONAS, "refPersona",
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

func ProcesarPaper(path string, meta *Frontmatter, tracker *d.TrackerDependencias) error {
	nombreRevista := strings.TrimSpace(meta.NombreRevista)
	if nombreRevista == "" {
		nombreRevista = "No fue ingresado - TODO"
	}
	err := tracker.Cargar(TABLA_REVISTAS_PAPER, []d.RelacionTabla{}, nombreRevista)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar revista con error: %v", err)
	}

	anio, err := strconv.Atoi(meta.Anio)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("obtener anio del paper con error: %v", err)
	}

	volumen := NumeroODefault(meta.VolumenInforme, 0)
	numero := NumeroODefault(meta.NumeroInforme, 0)
	paginaInicio := NumeroODefault(meta.Paginas.Inicio, 0)
	paginaFinal := NumeroODefault(meta.Paginas.Final, 1)

	err = tracker.Cargar(TABLA_PAPERS,
		[]d.RelacionTabla{
			d.NewRelacionSimple(TABLA_ARCHIVOS, "refArchivo", path),
			d.NewRelacionSimple(TABLA_REVISTAS_PAPER, "refRevista", nombreRevista),
			d.NewRelacionSimple(TABLA_COLECCIONES, "refColeccion", "Papers"),
		},
		meta.TituloInforme,
		meta.SubtituloInforme,
		anio,
		volumen,
		numero,
		paginaInicio,
		paginaFinal,
		meta.Url,
	)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar paper con error: %v", err)
	}

	datosCaracteristicosPaper := []any{meta.TituloInforme, anio, volumen, numero, paginaInicio, paginaFinal}

	for _, autor := range meta.Autores {
		nombre := strings.TrimSpace(autor.Nombre)
		apellido := strings.TrimSpace(autor.Apellido)

		err = tracker.Cargar(TABLA_PERSONAS, []d.RelacionTabla{}, nombre, apellido)
		if HABILITAR_ERROR && err != nil {
			return fmt.Errorf("cargar persona con error: %v", err)
		}

		err = tracker.Cargar(TABLA_ESCRITORES_PAPER,
			[]d.RelacionTabla{
				d.NewRelacionSimple(TABLA_PAPERS, "refPaper", datosCaracteristicosPaper...),
				d.NewRelacionSimple(TABLA_PERSONAS, "refPersona",
					nombre,
					apellido,
				),
			},
			PAPER_AUTOR,
		)
		if HABILITAR_ERROR && err != nil {
			return fmt.Errorf("cargar autor del paper con error: %v", err)
		}
	}

	for _, editor := range meta.Editores {
		nombre := strings.TrimSpace(editor.Nombre)
		apellido := strings.TrimSpace(editor.Apellido)

		err = tracker.Cargar(TABLA_PERSONAS, []d.RelacionTabla{}, nombre, apellido)
		if HABILITAR_ERROR && err != nil {
			return fmt.Errorf("cargar persona con error: %v", err)
		}

		err = tracker.Cargar(TABLA_ESCRITORES_PAPER,
			[]d.RelacionTabla{
				d.NewRelacionSimple(TABLA_PAPERS, "refPaper", datosCaracteristicosPaper...),
				d.NewRelacionSimple(TABLA_PERSONAS, "refPersona",
					nombre,
					apellido,
				),
			},
			PAPER_EDITOR,
		)
		if HABILITAR_ERROR && err != nil {
			return fmt.Errorf("cargar editor del paper con error: %v", err)
		}
	}
	return nil
}

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

type TipoMateria string

const (
	MATERIA_REAL        = "Materia"
	MATERIA_EQUIVALENTE = "Equivalente"
)

type TipoCurso string

const (
	CURSO_ONLINE     = "Online"
	CURSO_PRESENCIAL = "Presencial"
)

type TipoEscritorPaper string

const (
	PAPER_EDITOR = "Editor"
	PAPER_AUTOR  = "Autor"
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
