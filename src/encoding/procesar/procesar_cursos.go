package procesar

import (
	"fmt"
	d "own_wiki/system_protocol/dependencias"
	"strconv"
	"strings"
)

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
		EtapaODefault(meta.Etapa, ETAPA_SIN_EMPEZAR),
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
		EtapaODefault(meta.Etapa, ETAPA_SIN_EMPEZAR),
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
