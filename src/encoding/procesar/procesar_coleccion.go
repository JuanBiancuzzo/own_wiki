package procesar

import (
	"fmt"
	d "own_wiki/system_protocol/dependencias"
	"strconv"
	"strings"
)

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
