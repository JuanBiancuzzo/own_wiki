package fs

import (
	"own_wiki/system_protocol/db"
	e "own_wiki/system_protocol/estructura"
)

type Frontmatter struct {
	Tags                []string      `yaml:"tags,omitempty"`
	Dia                 string        `yaml:"dia,omitempty"`
	Etapa               string        `yaml:"etapa,omitempty"`
	Aliases             []string      `yaml:"aliases,omitempty"`
	Referencias         []string      `yaml:"referencias,omitempty"`
	NumReferncia        int           `yaml:"numReferncia,omitempty"`
	TipoCita            string        `yaml:"tipoCita,omitempty"`
	Previo              string        `yaml:"previo,omitempty"`
	Num                 int           `yaml:"num,omitempty"`
	Url                 string        `yaml:"url,omitempty"`
	Nombre              string        `yaml:"nombre,omitempty"`
	Articulo            []Articulo    `yaml:"articulo,omitempty"`
	Capitulo            string        `yaml:"capitulo,omitempty"`
	Fecha               string        `yaml:"fecha,omitempty"`
	NombreResumen       string        `yaml:"nombreResumen,omitempty"`
	MateriaResumen      string        `yaml:"materiaResumen,omitempty"`
	Anio                string        `yaml:"anio,omitempty"`
	Tipo                string        `yaml:"tipo,omitempty"`
	NombreAutores       []Persona     `yaml:"nombreAutores,omitempty"`
	Estado              string        `yaml:"estado,omitempty"`
	NombreCanal         string        `yaml:"nombreCanal,omitempty"`
	NombreVideo         string        `yaml:"nombreVideo,omitempty"`
	NombreArticulo      string        `yaml:"nombreArticulo,omitempty"`
	Editorial           string        `yaml:"editorial,omitempty"`
	Capitulos           []Capitulo    `yaml:"capitulos,omitempty"`
	TituloObra          string        `yaml:"tituloObra,omitempty"`
	SubtituloObra       string        `yaml:"subtituloObra,omitempty"`
	Edicion             string        `yaml:"edicion,omitempty"`
	Cover               string        `yaml:"cover,omitempty"`
	Volumen             string        `yaml:"volumen,omitempty"`
	NombreTema          string        `yaml:"nombreTema,omitempty"`
	Parte               string        `yaml:"parte,omitempty"`
	Curso               string        `yaml:"curso,omitempty"`
	NombreCurso         string        `yaml:"nombreCurso,omitempty"`
	FechaCurso          string        `yaml:"fechaCurso,omitempty"`
	TipoCurso           e.TipoCurso   `yaml:"tipoCurso,omitempty"`
	Profesores          []int         `yaml:"profesores,omitempty"`
	Autores             []Persona     `yaml:"autores,omitempty"`
	Editores            []Persona     `yaml:"editores,omitempty"`
	NumeroInforme       string        `yaml:"numeroInforme,omitempty"`
	TituloInforme       string        `yaml:"tituloInforme,omitempty"`
	SubtituloInforme    string        `yaml:"subtituloInforme,omitempty"`
	NombreRevista       string        `yaml:"nombreRevista,omitempty"`
	VolumenInforme      string        `yaml:"volumenRevista,omitempty"`
	Paginas             Pagina        `yaml:"paginas,omitempty"`
	Planes              []string      `yaml:"planes,omitempty"`
	TieneCodigo         string        `yaml:"tieneCodigo,omitempty"`
	NombreMateria       string        `yaml:"nombreMateria,omitempty"`
	NombreReducido      string        `yaml:"nombreReducido,omitempty"`
	PathCarrera         string        `yaml:"pathCarrera,omitempty"`
	Plan                string        `yaml:"plan,omitempty"`
	Codigo              string        `yaml:"codigo,omitempty"`
	Correlativas        []Correlativa `yaml:"correlativas,omitempty"`
	NombrePagina        string        `yaml:"nombrePagina,omitempty"`
	FechaPublicacion    string        `yaml:"fechaPublicacion,omitempty"`
	TituloArticulo      string        `yaml:"tituloArticulo,omitempty"`
	Cuatri              string        `yaml:"cuatri,omitempty"`
	NombreDistribuucion string        `yaml:"nombreDistribucion,omitempty"`
	TipoDistribucion    string        `yaml:"tipoDistribucion,omitempty"`
	Equivalencia        string        `yaml:"equivalencia,omitempty"`
	NombreSubtema       string        `yaml:"nombreSubtema,omitempty"`
}
type Capitulo struct {
	NumeroCapitulo string    `yaml:"numeroCapitulo"`
	NombreCapitulo string    `yaml:"nombreCapitulo,omitempty"`
	NumReferencia  int       `yaml:"numReferencia,omitempty"`
	Editores       []Persona `yaml:"editores,omitempty"`
	Paginas        Pagina    `yaml:"paginas,omitempty"`
}

type Persona struct {
	Nombre   string `yaml:"nombre"`
	Apellido string `yaml:"apellido"`
}

type Pagina struct {
	Inicio string `yaml:"inicio"`
	Final  string `yaml:"final"`
}

type Correlativa struct {
	Path string        `yaml:"path"`
	Tipo e.TipoMateria `yaml:"tipo"`
}

type Articulo struct {
	Tipo        string `yaml:"tipo,omitempty"`
	Enumeracion int    `yaml:"enumeracion,omitempty"`
	Texto       string `yaml:"texto,omitempty"`
	Textos      []struct {
		Tipo  string `yaml:"tipo,omitempty"`
		Texto string `yaml:"texto,omitempty"`
	} `yaml:"textos,omitempty"`
}

func (f *Frontmatter) CrearConstructorLibro() *e.ConstructorLibro {

	autores := make([]*e.Persona, len(f.Autores))
	for i, autor := range f.Autores {
		autores[i] = e.NewPersona(autor.Nombre, autor.Apellido)
	}
	capitulos := make([]*e.Capitulo, len(f.Capitulos))
	for i, capitulo := range f.Capitulos {
		editores := make([]*e.Persona, len(capitulo.Editores))
		for i, editor := range capitulo.Editores {
			editores[i] = e.NewPersona(editor.Nombre, editor.Apellido)
		}

		capitulos[i] = e.NewCapitulo(
			capitulo.NumeroCapitulo,
			capitulo.NombreCapitulo,
			editores,
			capitulo.Paginas.Inicio,
			capitulo.Paginas.Final,
		)
	}

	return e.NewConstructorLibro(
		f.TituloObra,
		f.SubtituloObra,
		f.Editorial,
		f.Anio,
		f.Edicion,
		f.Volumen,
		f.Url,
		autores,
		capitulos,
	)
}

func (f *Frontmatter) CrearConstructorPaper() (*e.ConstructorPaper, error) {
	if f.NombreRevista == "" {
		f.NombreRevista = "No fue ingresado"
	}

	autores := make([]*e.Persona, len(f.Autores))
	for i, autor := range f.Autores {
		autores[i] = e.NewPersona(autor.Nombre, autor.Apellido)
	}
	editores := make([]*e.Persona, len(f.Editores))
	for i, editor := range f.Editores {
		editores[i] = e.NewPersona(editor.Nombre, editor.Apellido)
	}

	return e.NewConstructorPaper(
		f.TituloInforme,
		f.SubtituloInforme,
		f.NombreRevista,
		f.VolumenInforme,
		f.NumeroInforme,
		f.Paginas.Inicio,
		f.Paginas.Final,
		f.Anio,
		f.Url,
		autores,
		editores,
	)
}

func (f *Frontmatter) CrearConstructorCurso() (*e.ConstructorCurso, error) {
	profesores := make([]*e.Persona, len(f.NombreAutores))
	for i, profesor := range f.NombreAutores {
		profesores[i] = e.NewPersona(profesor.Nombre, profesor.Apellido)
	}

	return e.NewConstructorCurso(
		f.NombreCurso,
		f.Etapa,
		f.FechaCurso,
		f.NombrePagina,
		f.Url,
		profesores,
	)
}

func (f *Frontmatter) CrearConstructorCursoPresencial() (*e.ConstructorCursoPresencial, error) {
	profesores := make([]*e.Persona, len(f.NombreAutores))
	for i, profesor := range f.NombreAutores {
		profesores[i] = e.NewPersona(profesor.Nombre, profesor.Apellido)
	}

	return e.NewConstructorCursoPresencial(
		f.NombreCurso,
		f.Etapa,
		f.FechaCurso,
		profesores,
	)
}

func CargarInfo(info *db.InfoArchivos, meta *Frontmatter) {
	for _, tag := range meta.Tags {
		info.MaxTags = max(info.MaxTags, uint32(len(tag)))
	}

	// Libros
	for _, autor := range meta.Autores {
		info.MaxNombre = max(info.MaxNombre, uint32(len(autor.Nombre)))
		info.MaxApellido = max(info.MaxApellido, uint32(len(autor.Apellido)))
	}

	info.MaxNombreLibro = max(info.MaxNombreLibro, uint32(len(meta.TituloObra)))
	info.MaxNombreLibro = max(info.MaxNombreLibro, uint32(len(meta.SubtituloObra)))
	for _, capitulo := range meta.Capitulos {
		info.MaxNombreLibro = max(info.MaxNombreLibro, uint32(len(capitulo.NombreCapitulo)))
	}

	info.MaxEditorial = max(info.MaxEditorial, uint32(len(meta.Editorial)))
	info.MaxUrl = max(info.MaxUrl, uint32(len(meta.Url)))

	// Distribuciones
}
