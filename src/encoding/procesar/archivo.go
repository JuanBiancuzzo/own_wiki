package procesar

import (
	"fmt"
	"os"
	d "own_wiki/system_protocol/dependencias"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
)

const HABILITAR_ERROR = true

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
	funcionesProcesar[TAG_NOTA_FACULTAD] = ProcesarNota
	funcionesProcesar[TAG_NOTA_CURSO] = ProcesarNota
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

func ProcesarNota(path string, meta *Frontmatter, tracker *d.TrackerDependencias) error {
	fecha, err := d.NewDate(meta.Dia)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("obtener fecha de la nota con error: %v", err)
	}

	nombreNota := Nombre(path)

	err = tracker.Cargar(TABLA_NOTAS,
		[]d.RelacionTabla{
			d.NewRelacionSimple(TABLA_ARCHIVOS, "refArchivo", path),
		},
		nombreNota,
		EtapaODefault(meta.Etapa, ETAPA_SIN_EMPEZAR),
		fecha.Representacion(),
	)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar nota con error: %v", err)
	}

	for _, vFacultad := range meta.VinculoFacultad {
		err = tracker.Cargar(TABLA_NOTAS_VINCULO,
			[]d.RelacionTabla{
				d.NewRelacionSimple(TABLA_NOTAS, "refNota", nombreNota),
				d.NewRelacionCompleja(TABLA_TEMA_MATERIA, "refVinculo",
					[]d.RelacionTabla{
						d.NewRelacionCompleja(TABLA_MATERIAS, "refMateria", []d.RelacionTabla{
							d.NewRelacionSimple(TABLA_CARRERAS, "refCarrera", vFacultad.NombreCarrera),
						}, vFacultad.NombreMateria),
					},
					vFacultad.NombreTema,
					NumeroODefault(vFacultad.CapituloTema, 1),
				),
			},
			TN_FACULTAD,
		)
		if HABILITAR_ERROR && err != nil {
			return fmt.Errorf("cargar vinculo con nota de faultad con error: %v", err)
		}
	}

	for _, vCurso := range meta.VinculoCurso {
		tablaCurso := TABLA_CURSOS_ONLINE
		if vCurso.TipoCurso == CURSO_PRESENCIAL {
			tablaCurso = TABLA_CURSOS_PRESENECIAL
		}

		anio, err := strconv.Atoi(vCurso.Anio)
		if HABILITAR_ERROR && err != nil {
			return fmt.Errorf("obtener anio del curso %s con error: %v", vCurso.TipoCurso, err)
		}

		err = tracker.Cargar(TABLA_NOTAS_VINCULO,
			[]d.RelacionTabla{
				d.NewRelacionSimple(TABLA_NOTAS, "refNota", nombreNota),
				d.NewRelacionCompleja(TABLA_TEMA_CURSO, "refVinculo",
					[]d.RelacionTabla{
						d.NewRelacionSimple(tablaCurso, "refCurso", vCurso.NombreCurso, anio),
					},
					vCurso.NombreTema,
					vCurso.TipoCurso,
					NumeroODefault(vCurso.CapituloTema, 1),
				),
			},
			TN_CURSO,
		)
		if HABILITAR_ERROR && err != nil {
			return fmt.Errorf("cargar vinculo con nota de curso con error: %v", err)
		}
	}

	return nil
}

func Nombre(path string) string {
	separacion := strings.Split(path, "/")
	return strings.ReplaceAll(separacion[len(separacion)-1], ".md", "")
}
