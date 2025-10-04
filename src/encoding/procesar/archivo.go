package procesar

import (
	"fmt"
	"os"
	d "own_wiki/system_protocol/dependencias"
	"strconv"
	"strings"

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

	err := tracker.Cargar(TABLA_ARCHIVOS, d.ConjuntoDato{"path": path})
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

	vinculosNota := make([]d.ConjuntoDato, len(meta.VinculoFacultad)+len(meta.VinculoCurso))
	for i, vFacultad := range meta.VinculoFacultad {
		vinculosNota[i] = d.ConjuntoDato{
			"refVinculo": d.NewRelacion(TABLA_TEMA_MATERIA, d.ConjuntoDato{
				"nombre":   vFacultad.NombreTema,
				"capitulo": NumeroODefault(vFacultad.CapituloTema, 1),
				"refMateria": d.NewRelacion(TABLA_MATERIAS, d.ConjuntoDato{
					"nombre":     vFacultad.NombreMateria,
					"refCarrera": d.NewRelacion(TABLA_CARRERAS, d.ConjuntoDato{"nombre": vFacultad.NombreCarrera}),
				}),
			}),
		}
	}

	desfaseVinculo := len(meta.VinculoFacultad)
	for i, vCurso := range meta.VinculoCurso {
		tablaCurso := TABLA_CURSOS_ONLINE
		if vCurso.TipoCurso == CURSO_PRESENCIAL {
			tablaCurso = TABLA_CURSOS_PRESENECIAL
		}

		anio, err := strconv.Atoi(vCurso.Anio)
		if HABILITAR_ERROR && err != nil {
			return fmt.Errorf("obtener anio del curso %s con error: %v", vCurso.TipoCurso, err)
		}

		vinculosNota[i+desfaseVinculo] = d.ConjuntoDato{
			"refVinculo": d.NewRelacion(TABLA_TEMA_CURSO, d.ConjuntoDato{
				"nombre":   vCurso.NombreTema,
				"capitulo": NumeroODefault(vCurso.CapituloTema, 1),
				"refCurso": d.NewRelacion(tablaCurso, d.ConjuntoDato{
					"nombre": vCurso.NombreCurso,
					"anio":   anio,
				}),
			}),
		}
	}

	err = tracker.Cargar(TABLA_NOTAS, d.ConjuntoDato{
		"nombre":      Nombre(path),
		"etapa":       EtapaODefault(meta.Etapa, ETAPA_SIN_EMPEZAR),
		"dia":         fecha.Representacion(),
		"refArchivo":  d.NewRelacion(TABLA_ARCHIVOS, d.ConjuntoDato{"path": path}),
		"refVinculos": vinculosNota,
	})
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar nota con error: %v", err)
	}

	return nil
}

func Nombre(path string) string {
	separacion := strings.Split(path, "/")
	return strings.ReplaceAll(separacion[len(separacion)-1], ".md", "")
}
