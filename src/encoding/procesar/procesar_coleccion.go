package procesar

import (
	"fmt"
	d "own_wiki/system_protocol/dependencias"
	"strconv"
	"strings"
)

func ProcesarColeccion(path string, meta *Frontmatter, tracker *d.TrackerDependencias) error {
	err := tracker.Cargar(TABLA_COLECCIONES, d.ConjuntoDato{
		"nombre":     Nombre(path),
		"refArchivo": d.NewRelacion(TABLA_ARCHIVOS, d.ConjuntoDato{"path": path}),
	})
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

	err = tracker.Cargar(TABLA_DISTRIBUCIONES, d.ConjuntoDato{
		"nombre":       meta.NombreDistribuucion,
		"tipo":         tipoDistribucion,
		"refArchivo":   d.NewRelacion(TABLA_ARCHIVOS, d.ConjuntoDato{"path": path}),
		"refColeccion": d.NewRelacion(TABLA_COLECCIONES, d.ConjuntoDato{"nombre": "Distribuciones"}),
	})
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar distribuciones con error: %v", err)
	}
	return nil
}

func ProcesarLibro(path string, meta *Frontmatter, tracker *d.TrackerDependencias) error {
	err := tracker.Cargar(TABLA_EDITORIALES, d.ConjuntoDato{"editorial": meta.Editorial})
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar editoriales con error: %v", err)
	}

	anio, err := strconv.Atoi(meta.Anio)
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("obtener anio del libro con error: %v", err)
	}

	edicion := NumeroODefault(meta.Edicion, 1)
	volumen := NumeroODefault(meta.Volumen, 0)

	autores := make([]d.ConjuntoDato, len(meta.NombreAutores))
	for i, autor := range meta.NombreAutores {
		datosAutor := d.ConjuntoDato{
			"nombre":   strings.TrimSpace(autor.Nombre),
			"apellido": strings.TrimSpace(autor.Apellido),
		}

		err = tracker.Cargar(TABLA_PERSONAS, datosAutor)
		if HABILITAR_ERROR && err != nil {
			return fmt.Errorf("cargar persona con error: %v", err)
		}

		autores[i] = d.ConjuntoDato{"refAutor": d.NewRelacion(TABLA_PERSONAS, datosAutor)}
	}

	err = tracker.Cargar(TABLA_LIBROS, d.ConjuntoDato{
		"titulo":       meta.TituloObra,
		"subtitulo":    meta.SubtituloObra,
		"anio":         anio,
		"edicion":      edicion,
		"volumen":      volumen,
		"url":          meta.Url,
		"refArchivo":   d.NewRelacion(TABLA_ARCHIVOS, d.ConjuntoDato{"path": path}),
		"refEditorial": d.NewRelacion(TABLA_EDITORIALES, d.ConjuntoDato{"editorial": meta.Editorial}),
		"refColeccion": d.NewRelacion(TABLA_COLECCIONES, d.ConjuntoDato{"nombre": "Biblioteca"}),
		"refAutores":   autores,
	})
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar libro con error: %v", err)
	}

	for _, capitulo := range meta.Capitulos {
		numero := NumeroODefault(capitulo.NumeroCapitulo, 0)
		paginaInicio := NumeroODefault(capitulo.Paginas.Inicio, 0)
		paginaFinal := NumeroODefault(capitulo.Paginas.Final, 1)

		editores := make([]d.ConjuntoDato, len(capitulo.Editores))
		for i, editor := range capitulo.Editores {
			datosEditor := d.ConjuntoDato{
				"nombre":   strings.TrimSpace(editor.Nombre),
				"apellido": strings.TrimSpace(editor.Apellido),
			}

			err = tracker.Cargar(TABLA_PERSONAS, datosEditor)
			if HABILITAR_ERROR && err != nil {
				return fmt.Errorf("cargar persona con error: %v", err)
			}

			editores[i] = d.ConjuntoDato{"refEditor": d.NewRelacion(TABLA_PERSONAS, datosEditor)}
		}

		err = tracker.Cargar(TABLA_CAPITULOS, d.ConjuntoDato{
			"numero":       numero,
			"nombre":       capitulo.NombreCapitulo,
			"paginaInicio": paginaInicio,
			"paginaFinal":  paginaFinal,
			"refArchivo":   d.NewRelacion(TABLA_ARCHIVOS, d.ConjuntoDato{"path": path}),
			"refLibro": d.NewRelacion(TABLA_LIBROS, d.ConjuntoDato{
				"titulo":  meta.TituloObra,
				"anio":    anio,
				"edicion": edicion,
				"volumen": volumen,
			}),
			"refEditores": editores,
		})
		if HABILITAR_ERROR && err != nil {
			return fmt.Errorf("cargar capitulo con error: %v", err)
		}
	}

	return nil
}

func ProcesarPaper(path string, meta *Frontmatter, tracker *d.TrackerDependencias) error {
	nombreRevista := strings.TrimSpace(meta.NombreRevista)
	if nombreRevista == "" {
		nombreRevista = "No fue ingresado - TODO"
	}
	err := tracker.Cargar(TABLA_REVISTAS_PAPER, d.ConjuntoDato{"nombre": nombreRevista})
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

	autores := make([]d.ConjuntoDato, len(meta.Autores))
	for i, autor := range meta.Autores {
		datosAutor := d.ConjuntoDato{
			"nombre":   strings.TrimSpace(autor.Nombre),
			"apellido": strings.TrimSpace(autor.Apellido),
		}

		err = tracker.Cargar(TABLA_PERSONAS, datosAutor)
		if HABILITAR_ERROR && err != nil {
			return fmt.Errorf("cargar persona con error: %v", err)
		}

		autores[i] = d.ConjuntoDato{"refAutor": d.NewRelacion(TABLA_PERSONAS, datosAutor)}
	}

	editores := make([]d.ConjuntoDato, len(meta.Editores))
	for i, editor := range meta.Editores {
		datosEditor := d.ConjuntoDato{
			"nombre":   strings.TrimSpace(editor.Nombre),
			"apellido": strings.TrimSpace(editor.Apellido),
		}

		err = tracker.Cargar(TABLA_PERSONAS, datosEditor)
		if HABILITAR_ERROR && err != nil {
			return fmt.Errorf("cargar persona con error: %v", err)
		}

		editores[i] = d.ConjuntoDato{"refEditor": d.NewRelacion(TABLA_PERSONAS, datosEditor)}
	}

	err = tracker.Cargar(TABLA_PAPERS, d.ConjuntoDato{
		"titulo":       meta.TituloInforme,
		"subtitulo":    meta.SubtituloInforme,
		"anio":         anio,
		"volumen":      volumen,
		"numero":       numero,
		"paginaInicio": paginaInicio,
		"paginaFinal":  paginaFinal,
		"doi":          meta.Url,
		"refArchivo":   d.NewRelacion(TABLA_ARCHIVOS, d.ConjuntoDato{"path": path}),
		"refRevista":   d.NewRelacion(TABLA_REVISTAS_PAPER, d.ConjuntoDato{"nombre": nombreRevista}),
		"refColeccion": d.NewRelacion(TABLA_COLECCIONES, d.ConjuntoDato{"nombre": "Papers"}),
		"refAutores":   autores,
		"refEditores":  editores,
	})
	if HABILITAR_ERROR && err != nil {
		return fmt.Errorf("cargar paper con error: %v", err)
	}

	return nil
}
