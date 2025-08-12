package procesar

import (
	"fmt"
	d "own_wiki/system_protocol/dependencias"
	"strings"
)

func ProcesarCarrera(path string, meta *Frontmatter, tracker *d.TrackerDependencias) error {
	nombreCarrera := Nombre(path)
	err := tracker.Cargar(TABLA_CARRERAS,
		[]d.RelacionTabla{
			d.NewRelacionSimple(TABLA_ARCHIVOS, "refArchivo", path),
		},
		nombreCarrera,
		EtapaODefault(meta.Etapa, ETAPA_SIN_EMPEZAR),
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
		EtapaODefault(meta.Etapa, ETAPA_SIN_EMPEZAR),
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
