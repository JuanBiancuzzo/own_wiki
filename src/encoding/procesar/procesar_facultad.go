package procesar

import (
	"fmt"
	d "own_wiki/system_protocol/dependencias"
	"strings"
)

func ProcesarCarrera(path string, meta *Frontmatter, tracker *d.TrackerDependencias) error {
	nombreCarrera := Nombre(path)
	err := tracker.Cargar(TABLA_CARRERAS, d.ConjuntoDato{
		"nombre":             nombreCarrera,
		"etapa":              EtapaODefault(meta.Etapa, ETAPA_SIN_EMPEZAR),
		"tieneCodigoMateria": strings.ToLower(strings.TrimSpace(meta.TieneCodigo)) == "true",
		"refArchivo":         d.NewRelacion(TABLA_ARCHIVOS, d.ConjuntoDato{"path": path}),
	})

	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar carrera con error: %v", err)
	}
	return nil
}

func ProcesarMateria(path string, meta *Frontmatter, tracker *d.TrackerDependencias) error {
	err := tracker.Cargar(TABLA_PLANES, d.ConjuntoDato{"nombre": meta.Plan})
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar plan con error: %v", err)
	}

	anio, cuatrimestre, err := ObtenerCuatrimestreParte(meta.Cuatri)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("obtener cuatrimestre con error: %v", err)
	}

	err = tracker.Cargar(TABLA_CUATRI, d.ConjuntoDato{
		"anio":         anio,
		"cuatrimestre": cuatrimestre,
	})
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar cuatri con error: %v", err)
	}

	materiasCorrelativas := make([]d.ConjuntoDato, len(meta.Correlativas))
	for i, infoCorrelativa := range meta.Correlativas {
		tablaCorrelativa := TABLA_MATERIAS
		if infoCorrelativa.Tipo == MATERIA_EQUIVALENTE {
			tablaCorrelativa = TABLA_MATERIAS_EQ
		}

		materiasCorrelativas[i] = d.ConjuntoDato{
			"refCorrelativa": d.NewRelacion(tablaCorrelativa, d.ConjuntoDato{
				"nombre":     infoCorrelativa.Materia,
				"refCarrera": d.NewRelacion(TABLA_CARRERAS, d.ConjuntoDato{"nombre": meta.NombreCarrera}),
			}),
		}
	}

	err = tracker.Cargar(TABLA_MATERIAS, d.ConjuntoDato{
		"nombre":     meta.NombreMateria,
		"etapa":      EtapaODefault(meta.Etapa, ETAPA_SIN_EMPEZAR),
		"codigo":     meta.Codigo,
		"refArchivo": d.NewRelacion(TABLA_ARCHIVOS, d.ConjuntoDato{"path": path}),
		"refCarrera": d.NewRelacion(TABLA_CARRERAS, d.ConjuntoDato{"nombre": meta.NombreCarrera}),
		"refPlan":    d.NewRelacion(TABLA_PLANES, d.ConjuntoDato{"nombre": meta.Plan}),
		"refCuatrimestre": d.NewRelacion(TABLA_CUATRI, d.ConjuntoDato{
			"anio":         anio,
			"cuatrimestre": cuatrimestre,
		}),
		"refCorrelativas": materiasCorrelativas,
	})
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar materia con error: %v", err)
	}

	return nil
}

func ProcesarMateriaEquivalente(path string, meta *Frontmatter, tracker *d.TrackerDependencias) error {
	infoMateria := meta.MateriaEquivalente

	materiasCorrelativas := make([]d.ConjuntoDato, len(meta.Correlativas))
	for i, infoCorrelativa := range meta.Correlativas {
		tablaCorrelativa := TABLA_MATERIAS
		if infoCorrelativa.Tipo == MATERIA_EQUIVALENTE {
			tablaCorrelativa = TABLA_MATERIAS_EQ
		}

		materiasCorrelativas[i] = d.ConjuntoDato{
			"refCorrelativa": d.NewRelacion(tablaCorrelativa, d.ConjuntoDato{
				"nombre":     infoCorrelativa.Materia,
				"refCarrera": d.NewRelacion(TABLA_CARRERAS, d.ConjuntoDato{"nombre": meta.NombreCarrera}),
			}),
		}
	}

	err := tracker.Cargar(TABLA_MATERIAS_EQ, d.ConjuntoDato{
		"nombre":     meta.NombreMateria,
		"codigo":     meta.Codigo,
		"refArchivo": d.NewRelacion(TABLA_ARCHIVOS, d.ConjuntoDato{"path": path}),
		"refCarrera": d.NewRelacion(TABLA_CARRERAS, d.ConjuntoDato{"nombre": meta.NombreCarrera}),
		"refMateria": d.NewRelacion(TABLA_MATERIAS, d.ConjuntoDato{
			"nombre":     infoMateria.NombreMateria,
			"refCarrera": d.NewRelacion(TABLA_CARRERAS, d.ConjuntoDato{"nombre": infoMateria.Carrera}),
		}),
		"refCorrelativas": materiasCorrelativas,
	})
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar materia equivalente con error: %v", err)
	}

	return nil
}

func ProcesarTemaMateria(path string, meta *Frontmatter, tracker *d.TrackerDependencias) error {
	infoTema := meta.InfoTemaMateria
	err := tracker.Cargar(TABLA_TEMA_MATERIA, d.ConjuntoDato{
		"nombre":     meta.NombreResumen,
		"capitulo":   NumeroODefault(meta.Capitulo, 1),
		"parte":      NumeroODefault(meta.Parte, 0),
		"refArchivo": d.NewRelacion(TABLA_ARCHIVOS, d.ConjuntoDato{"path": path}),
		"refMateria": d.NewRelacion(TABLA_MATERIAS, d.ConjuntoDato{
			"nombre":     infoTema.Materia,
			"refCarrera": d.NewRelacion(TABLA_CARRERAS, d.ConjuntoDato{"nombre": infoTema.Carrera}),
		}),
	})
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar tema materia con error: %v", err)
	}
	return nil
}
