package fs

import (
	"fmt"
	"strconv"
)

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
	Planes              []string   `yaml:"planes,omitempty"`
	TieneCodigo         string     `yaml:"tieneCodigo,omitempty"`
	NombreMateria       string     `yaml:"nombreMateria,omitempty"`
	NombreReducido      string     `yaml:"nombreReducido,omitempty"`
	Plan                string     `yaml:"plan,omitempty"`
	Codigo              string     `yaml:"codigo,omitempty"`
	Correlativas        []string   `yaml:"correlativas,omitempty"`
	NombrePagina        string     `yaml:"nombrePagina,omitempty"`
	FechaPublicacion    string     `yaml:"fechaPublicacion,omitempty"`
	TituloArticulo      string     `yaml:"tituloArticulo,omitempty"`
	Cuatri              string     `yaml:"cuatri,omitempty"`
	NombreDistribuucion string     `yaml:"nombreDistribucion,omitempty"`
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

func NewLibro(titulo string, subtitulo string, anio string, idEditorial int64, edicion string, volumen string, url string, idArchivo int64) Libro {
	return Libro{
		Titulo:      titulo,
		Subtitulo:   subtitulo,
		Anio:        NumeroODefault(anio, 0),
		IdEditorial: idEditorial,
		Edicion:     NumeroODefault(edicion, 1),
		Volumen:     NumeroODefault(volumen, 0),
		Url:         url,
		IdArchivo:   idArchivo,
	}
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

func (c Capitulo) Valores() []any {
	return []any{
		NumeroODefault(c.NumeroCapitulo, 1),
		c.NombreCapitulo,
		NumeroODefault(c.Paginas.Inicio, 0),
		NumeroODefault(c.Paginas.Final, 0),
	}
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

type TipoDistribucion string

const (
	DISTRIBUCION_DISCRETA     = "Discreta"
	DISTRIBUCION_CONTINUA     = "Continua"
	DISTRIBUCION_MULTIVARIADA = "Multivariada"
)

type Distribucion struct {
	Tipo   TipoDistribucion
	Nombre string
}

func NewDistribucion(tipo TipoDistribucion, nombre string) Distribucion {
	return Distribucion{
		Tipo:   tipo,
		Nombre: nombre,
	}
}

type Etapa string

const (
	ETAPA_SIN_EMPEZAR = "SinEmpezar"
	ETAPA_EMPEZADO    = "Empezado"
	ETAPA_AMPLIAR     = "Ampliar"
	ETAPA_TERMINADO   = "Terminado"
)

type Carrera struct {
	Nombre      string
	Etapa       Etapa
	TieneCodigo bool
}

func NewCarrera(nombre string, etapa Etapa, tieneCodigo string) Carrera {
	return Carrera{
		Nombre:      nombre,
		Etapa:       etapa,
		TieneCodigo: BooleanoODefault(tieneCodigo, false),
	}
}

func (c Carrera) Valores() []any {
	return []any{
		c.Nombre,
		c.Etapa,
		c.TieneCodigo,
	}
}

type Materia struct {
}

func NumeroODefault(representacion string, valorDefault int) int {
	if nuevoValor, err := strconv.Atoi(representacion); err == nil {
		return nuevoValor
	} else {
		return valorDefault
	}
}

func BooleanoODefault(representacion string, valorDefault bool) bool {
	switch representacion {
	case "true":
		return true
	case "false":
		return false
	default:
		return valorDefault
	}
}

func ObtenerEtapa(representacionEtapa string) (Etapa, error) {
	var etapa Etapa
	switch representacionEtapa {
	case "sin-empezar":
		etapa = ETAPA_SIN_EMPEZAR
	case "empezado":
		etapa = ETAPA_EMPEZADO
	case "ampliar":
		etapa = ETAPA_AMPLIAR
	case "terminado":
		etapa = ETAPA_TERMINADO
	default:
		return ETAPA_SIN_EMPEZAR, fmt.Errorf("el tipo de etapa (%s) no es uno de los esperados", representacionEtapa)
	}

	return etapa, nil
}
