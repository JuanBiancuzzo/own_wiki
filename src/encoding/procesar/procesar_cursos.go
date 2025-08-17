package procesar

import (
	"fmt"
	d "own_wiki/system_protocol/dependencias"
	"strconv"
	"strings"
)

func ProcesarCursoOnline(path string, meta *Frontmatter, tracker *d.TrackerDependencias) error {
	err := tracker.Cargar(TABLA_PAGINAS_CURSOS, d.ConjuntoDato{"nombre": meta.NombrePagina})
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar paginas cursos con error: %v", err)
	}

	anio, err := strconv.Atoi(meta.FechaCurso)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("obtener anio del curso online con error: %v", err)
	}

	err = tracker.Cargar(TABLA_CURSOS_ONLINE, d.ConjuntoDato{
		"nombre":     meta.NombreCurso,
		"etapa":      EtapaODefault(meta.Etapa, ETAPA_SIN_EMPEZAR),
		"anio":       anio,
		"url":        meta.Url,
		"refArchivo": d.NewRelacion(TABLA_ARCHIVOS, d.ConjuntoDato{"path": path}),
		"refPagina":  d.NewRelacion(TABLA_PAGINAS_CURSOS, d.ConjuntoDato{"nombre": meta.NombrePagina}),
	})
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar curso online con error: %v", err)
	}

	for _, profesor := range meta.NombreAutores {
		nombre := strings.TrimSpace(profesor.Nombre)
		apellido := strings.TrimSpace(profesor.Apellido)

		err = tracker.Cargar(TABLA_PERSONAS, d.ConjuntoDato{
			"nombre":   nombre,
			"apellido": apellido,
		})
		if HABILITAR_ERROR && err != nil {
			return fmt.Errorf("cargar persona en curso online con error: %v", err)
		}

		err = tracker.Cargar(TABLA_PROFESORES_CURSO, d.ConjuntoDato{
			"refCurso": d.NewRelacion(TABLA_CURSOS_ONLINE, d.ConjuntoDato{
				"nombre": meta.NombreCurso,
				"anio":   anio,
			}),
			"refPersona": d.NewRelacion(TABLA_PERSONAS, d.ConjuntoDato{
				"nombre":   nombre,
				"apellido": apellido,
			}),
		})
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

	err = tracker.Cargar(TABLA_CURSOS_PRESENECIAL, d.ConjuntoDato{
		"nombre":     meta.NombreCurso,
		"etapa":      EtapaODefault(meta.Etapa, ETAPA_SIN_EMPEZAR),
		"anio":       anio,
		"refArchivo": d.NewRelacion(TABLA_ARCHIVOS, d.ConjuntoDato{"path": path}),
	})
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar curso presencial con error: %v", err)
	}

	for _, profesor := range meta.NombreAutores {
		nombre := strings.TrimSpace(profesor.Nombre)
		apellido := strings.TrimSpace(profesor.Apellido)

		err = tracker.Cargar(TABLA_PERSONAS, d.ConjuntoDato{
			"nombre":   nombre,
			"apellido": apellido,
		})
		if HABILITAR_ERROR && err != nil {
			return fmt.Errorf("cargar persona en curso presencial con error: %v", err)
		}

		err = tracker.Cargar(TABLA_PROFESORES_CURSO, d.ConjuntoDato{
			"refCurso": d.NewRelacion(TABLA_CURSOS_PRESENECIAL, d.ConjuntoDato{
				"nombre": meta.NombreCurso,
				"anio":   anio,
			}),
			"refPersona": d.NewRelacion(TABLA_PERSONAS, d.ConjuntoDato{
				"nombre":   nombre,
				"apellido": apellido,
			}),
		})
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

	err = tracker.Cargar(TABLA_TEMA_CURSO, d.ConjuntoDato{
		"nombre":     meta.NombreResumen,
		"capitulo":   NumeroODefault(meta.Capitulo, 1),
		"parte":      NumeroODefault(meta.Parte, 0),
		"refArchivo": d.NewRelacion(TABLA_ARCHIVOS, d.ConjuntoDato{"path": path}),
		"refCurso": d.NewRelacion(tablaCurso, d.ConjuntoDato{
			"nombre": infoTema.Curso,
			"anio":   anio,
		}),
	})
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar tema curso con error: %v", err)
	}
	return nil
}
