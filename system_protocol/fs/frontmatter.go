package fs

type Frontmatter struct {
	Tags                []string   `yaml:"tags,omitempty"`
	Dia                 string     `yaml:"dia,omitempty"`
	Etapa               string     `yaml:"etapa,omitempty"`
	Aliases             []string   `yaml:"aliases,omitempty"`
	Referencias         []string   `yaml:"referencias,omitempty"`
	NumReferncia        int        `yaml:"numReferncia,omitempty"`
	TipoCita            string     `yaml:"tipoCita,omitempty"`
	Previo              string     `yaml:"previo,omitempty"`
	Num                 int        `yaml:"num,omitempty"`
	Url                 string     `yaml:"url,omitempty"`
	Nombre              string     `yaml:"nombre,omitempty"`
	Articulo            []Articulo `yaml:"articulo,omitempty"`
	Capitulo            string     `yaml:"capitulo,omitempty"`
	Fecha               string     `yaml:"fecha,omitempty"`
	NombreResumen       string     `yaml:"nombreResumen,omitempty"`
	Anio                string     `yaml:"anio,omitempty"`
	Tipo                string     `yaml:"tipo,omitempty"`
	NombreAutores       []Autor    `yaml:"nombreAutores,omitempty"`
	Estado              string     `yaml:"estado,omitempty"`
	NombreCanal         string     `yaml:"nombreCanal,omitempty"`
	NombreVideo         string     `yaml:"nombreVideo,omitempty"`
	NombreArticulo      string     `yaml:"nombreArticulo,omitempty"`
	Editorial           string     `yaml:"editorial,omitempty"`
	Capitulos           []Capitulo `yaml:"capitulos,omitempty"`
	TituloObra          string     `yaml:"tituloObra,omitempty"`
	SubtituloObra       string     `yaml:"subtituloObra,omitempty"`
	Edicion             string     `yaml:"edicion,omitempty"`
	Cover               string     `yaml:"cover,omitempty"`
	Volumen             string     `yaml:"volumen,omitempty"`
	NombreTema          string     `yaml:"nombreTema,omitempty"`
	Parte               int        `yaml:"parte,omitempty"`
	Curso               string     `yaml:"curso,omitempty"`
	Profesores          []int      `yaml:"profesores,omitempty"`
	Autores             []Autor    `yaml:"autores,omitempty"`
	Editores            []string   `yaml:"editores,omitempty"`
	NumeroInforme       string     `yaml:"numeroInforme,omitempty"`
	TituloInforme       string     `yaml:"tituloInforme,omitempty"`
	NombreMateria       string     `yaml:"nombreMateria,omitempty"`
	NombreReducido      string     `yaml:"nombreReducido,omitempty"`
	Plan                string     `yaml:"plan,omitempty"`
	Codigo              string     `yaml:"codigo,omitempty"`
	Correlativas        []string   `yaml:"correlativas,omitempty"`
	NombrePagina        string     `yaml:"nombrePagina,omitempty"`
	FechaPublicacion    string     `yaml:"fechaPublicacion,omitempty"`
	TituloArticulo      string     `yaml:"tituloArticulo,omitempty"`
	Cuatri              string     `yaml:"cuatri,omitempty"`
	NombreDistribuucion string     `yaml:"nombreDistribuucion,omitempty"`
	TipoDistribucion    string     `yaml:"tipoDistribucion,omitempty"`
	Equivalencia        string     `yaml:"equivalencia,omitempty"`
	NombreSubtema       string     `yaml:"nombreSubtema,omitempty"`
}

type Libro struct {
	Titulo      string
	Subtitulo   string
	Anio        int
	IdEditorial int64
	Edicion     int
	Volumen     int
	Url         string
	IdArchivo   int64
}

func (l Libro) Valores() []any {
	return []any{
		l.Titulo,
		l.Subtitulo,
		l.Anio,
		l.IdEditorial,
		l.Edicion,
		l.Volumen,
		l.Url,
		l.IdArchivo,
	}
}

type Capitulo struct {
	NumeroCapitulo string  `yaml:"numeroCapitulo"`
	NombreCapitulo string  `yaml:"nombreCapitulo,omitempty"`
	NumReferencia  int     `yaml:"numReferencia,omitempty"`
	Editores       []Autor `yaml:"editores,omitempty"`
	Paginas        Pagina  `yaml:"paginas,omitempty"`
}

type Autor struct {
	Nombre   string `yaml:"nombre"`
	Apellido string `yaml:"apellido"`
}

type Pagina struct {
	Inicio string `yaml:"inicio"`
	Final  string `yaml:"final"`
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
