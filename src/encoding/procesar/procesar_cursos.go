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

	profesores := make([]d.ConjuntoDato, len(meta.NombreAutores))
	for i, profesor := range meta.NombreAutores {
		datosProfesor := d.ConjuntoDato{
			"nombre":   strings.TrimSpace(profesor.Nombre),
			"apellido": strings.TrimSpace(profesor.Apellido),
		}

		err = tracker.Cargar(TABLA_PERSONAS, datosProfesor)
		if HABILITAR_ERROR && err != nil {
			return fmt.Errorf("cargar persona en curso online con error: %v", err)
		}

		profesores[i] = d.ConjuntoDato{"refProfesor": d.NewRelacion(TABLA_PERSONAS, datosProfesor)}
	}

	bibliografia := make([]d.ConjuntoDato, len(meta.Referencias))
	for i, referencia := range meta.Referencias {
		datosReferencia := make(d.ConjuntoDato)
		if datosReferencia["numero"], err = strconv.Atoi(referencia); HABILITAR_ERROR && err != nil {
			return fmt.Errorf("al convertir la referencia '%s' tuve el error: %v", referencia, err)
		}

		err = tracker.Cargar(TABLA_REFERENCIAS, datosReferencia)
		if HABILITAR_ERROR && err != nil {
			return fmt.Errorf("cargar referencia en curso online con error: %v", err)
		}

		bibliografia[i] = d.ConjuntoDato{
			"refReferencia": d.NewRelacion(TABLA_REFERENCIAS, datosReferencia),
		}
	}

	err = tracker.Cargar(TABLA_CURSOS_ONLINE, d.ConjuntoDato{
		"nombre":          meta.NombreCurso,
		"etapa":           EtapaODefault(meta.Etapa, ETAPA_SIN_EMPEZAR),
		"anio":            anio,
		"url":             meta.Url,
		"refArchivo":      d.NewRelacion(TABLA_ARCHIVOS, d.ConjuntoDato{"path": path}),
		"refPagina":       d.NewRelacion(TABLA_PAGINAS_CURSOS, d.ConjuntoDato{"nombre": meta.NombrePagina}),
		"refProfesores":   profesores,
		"refBibliografia": bibliografia,
	})
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar curso online con error: %v", err)
	}

	return nil
}

func ProcesarCursoPresencial(path string, meta *Frontmatter, tracker *d.TrackerDependencias) error {
	anio, err := strconv.Atoi(meta.FechaCurso)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("obtener anio del curso presencial con error: %v", err)
	}

	profesores := make([]d.ConjuntoDato, len(meta.NombreAutores))
	for i, profesor := range meta.NombreAutores {
		datosProfesor := d.ConjuntoDato{
			"nombre":   strings.TrimSpace(profesor.Nombre),
			"apellido": strings.TrimSpace(profesor.Apellido),
		}

		err = tracker.Cargar(TABLA_PERSONAS, datosProfesor)
		if HABILITAR_ERROR && err != nil {
			return fmt.Errorf("cargar persona en curso online con error: %v", err)
		}

		profesores[i] = d.ConjuntoDato{"refProfesor": d.NewRelacion(TABLA_PERSONAS, datosProfesor)}
	}

	bibliografia := make([]d.ConjuntoDato, len(meta.Referencias))
	for i, referencia := range meta.Referencias {
		datosReferencia := make(d.ConjuntoDato)
		if datosReferencia["numero"], err = strconv.Atoi(referencia); HABILITAR_ERROR && err != nil {
			return fmt.Errorf("al convertir la referencia '%s' tuve el error: %v", referencia, err)
		}

		err = tracker.Cargar(TABLA_REFERENCIAS, datosReferencia)
		if HABILITAR_ERROR && err != nil {
			return fmt.Errorf("cargar referencia en curso presencial con error: %v", err)
		}

		bibliografia[i] = d.ConjuntoDato{
			"refReferencia": d.NewRelacion(TABLA_REFERENCIAS, datosReferencia),
		}
	}

	err = tracker.Cargar(TABLA_CURSOS_PRESENECIAL, d.ConjuntoDato{
		"nombre":          meta.NombreCurso,
		"etapa":           EtapaODefault(meta.Etapa, ETAPA_SIN_EMPEZAR),
		"anio":            anio,
		"refArchivo":      d.NewRelacion(TABLA_ARCHIVOS, d.ConjuntoDato{"path": path}),
		"refProfesores":   profesores,
		"refBibliografia": bibliografia,
	})
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar curso presencial con error: %v", err)
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

	bibliografia := make([]d.ConjuntoDato, len(meta.Referencias))
	for i, referencia := range meta.Referencias {
		datosReferencia := make(d.ConjuntoDato)
		if datosReferencia["numero"], err = strconv.Atoi(referencia); HABILITAR_ERROR && err != nil {
			return fmt.Errorf("al convertir la referencia '%s' tuve el error: %v", referencia, err)
		}

		err = tracker.Cargar(TABLA_REFERENCIAS, datosReferencia)
		if HABILITAR_ERROR && err != nil {
			return fmt.Errorf("cargar referencia en el tema del curso con error: %v", err)
		}

		bibliografia[i] = d.ConjuntoDato{
			"refReferencia": d.NewRelacion(TABLA_REFERENCIAS, datosReferencia),
		}
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
		"refBibliografia": bibliografia,
	})
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar tema curso con error: %v", err)
	}
	return nil
}
