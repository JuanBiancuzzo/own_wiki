package procesar

type Frontmatter struct {
	Tags                 []string          `yaml:"tags,omitempty"`
	Dia                  string            `yaml:"dia,omitempty"`
	VinculoFacultad      []VinculoFacultad `yaml:"vinculoFacultad,omitempty"`
	VinculoInvestigacion []string          `yaml:"vinculoInvestigacion,omitempty"`
	VinculoCurso         []VinculoCurso    `yaml:"vinculoCurso,omitempty"`
	VinculoProyecto      []string          `yaml:"vinculoProyecto,omitempty"`
	VinculoColeccion     []string          `yaml:"vinculoColeccion,omitempty"`
	Etapa                string            `yaml:"etapa,omitempty"`
	Aliases              []string          `yaml:"aliases,omitempty"`
	Referencias          []string          `yaml:"referencias,omitempty"`
	NumReferencia        int               `yaml:"numReferncia,omitempty"`
	TipoCita             string            `yaml:"tipoCita,omitempty"`
	Previo               string            `yaml:"previo,omitempty"`
	Num                  int               `yaml:"num,omitempty"`
	Url                  string            `yaml:"url,omitempty"`
	Nombre               string            `yaml:"nombre,omitempty"`
	Articulo             []Articulo        `yaml:"articulo,omitempty"`
	Capitulo             string            `yaml:"capitulo,omitempty"`
	Fecha                string            `yaml:"fecha,omitempty"`
	NombreResumen        string            `yaml:"nombreResumen,omitempty"`
	MateriaResumen       string            `yaml:"materiaResumen,omitempty"`
	Anio                 string            `yaml:"anio,omitempty"`
	Tipo                 string            `yaml:"tipo,omitempty"`
	NombreAutores        []Persona         `yaml:"nombreAutores,omitempty"`
	Estado               string            `yaml:"estado,omitempty"`
	NombreCanal          string            `yaml:"nombreCanal,omitempty"`
	NombreVideo          string            `yaml:"nombreVideo,omitempty"`
	NombreArticulo       string            `yaml:"nombreArticulo,omitempty"`
	Editorial            string            `yaml:"editorial,omitempty"`
	Capitulos            []Capitulo        `yaml:"capitulos,omitempty"`
	TituloObra           string            `yaml:"tituloObra,omitempty"`
	SubtituloObra        string            `yaml:"subtituloObra,omitempty"`
	Edicion              string            `yaml:"edicion,omitempty"`
	Cover                string            `yaml:"cover,omitempty"`
	Volumen              string            `yaml:"volumen,omitempty"`
	NombreTema           string            `yaml:"nombreTema,omitempty"`
	Parte                string            `yaml:"parte,omitempty"`
	Curso                string            `yaml:"curso,omitempty"`
	InfoCurso            InfoTemaCurso     `yaml:"infoCurso,omitempty"`
	NombreCurso          string            `yaml:"nombreCurso,omitempty"`
	FechaCurso           string            `yaml:"fechaCurso,omitempty"`
	TipoCurso            TipoCurso         `yaml:"tipoCurso,omitempty"`
	Profesores           []int             `yaml:"profesores,omitempty"`
	Autores              []Persona         `yaml:"autores,omitempty"`
	Editores             []Persona         `yaml:"editores,omitempty"`
	NumeroInforme        string            `yaml:"numeroInforme,omitempty"`
	TituloInforme        string            `yaml:"tituloInforme,omitempty"`
	SubtituloInforme     string            `yaml:"subtituloInforme,omitempty"`
	NombreRevista        string            `yaml:"nombreRevista,omitempty"`
	VolumenInforme       string            `yaml:"volumenRevista,omitempty"`
	Paginas              Pagina            `yaml:"paginas,omitempty"`
	Planes               []string          `yaml:"planes,omitempty"`
	TieneCodigo          string            `yaml:"tieneCodigo,omitempty"`
	NombreMateria        string            `yaml:"nombreMateria,omitempty"`
	NombreReducido       string            `yaml:"nombreReducido,omitempty"`
	PathCarrera          string            `yaml:"pathCarrera,omitempty"`
	NombreCarrera        string            `yaml:"nombreCarrera,omitempty"`
	MateriaEquivalente   InfoMateria       `yaml:"materiaEquivalente,omitempty"`
	InfoTemaMateria      InfoTemaMateria   `yaml:"infoTemaMateria,omitempty"`
	Plan                 string            `yaml:"plan,omitempty"`
	Codigo               string            `yaml:"codigo,omitempty"`
	Correlativas         []Correlativa     `yaml:"correlativas,omitempty"`
	NombrePagina         string            `yaml:"nombrePagina,omitempty"`
	FechaPublicacion     string            `yaml:"fechaPublicacion,omitempty"`
	TituloArticulo       string            `yaml:"tituloArticulo,omitempty"`
	Cuatri               string            `yaml:"cuatri,omitempty"`
	NombreDistribuucion  string            `yaml:"nombreDistribucion,omitempty"`
	TipoDistribucion     string            `yaml:"tipoDistribucion,omitempty"`
	Equivalencia         string            `yaml:"equivalencia,omitempty"`
	NombreSubtema        string            `yaml:"nombreSubtema,omitempty"`
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
	Materia string      `yaml:"materia"`
	Tipo    TipoMateria `yaml:"tipo"`
}

type InfoMateria struct {
	NombreMateria string `yaml:"nombre"`
	Carrera       string `yaml:"carrera"`
}

type InfoTemaMateria struct {
	Materia string `yaml:"materia"`
	Carrera string `yaml:"carrera"`
}

type InfoTemaCurso struct {
	Curso string    `yaml:"curso"`
	Tipo  TipoCurso `yaml:"tipo"`
	Anio  string    `yaml:"anio"`
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

type VinculoFacultad struct {
	NombreTema    string `yaml:"tema,omitempty"`
	CapituloTema  string `yaml:"capitulo,omitempty"`
	NombreMateria string `yaml:"materia,omitempty"`
	NombreCarrera string `yaml:"carrera,omitempty"`
}

type VinculoCurso struct {
	NombreTema   string    `yaml:"tema,omitempty"`
	CapituloTema string    `yaml:"capitulo,omitempty"`
	TipoCurso    TipoCurso `yaml:"tipo,omitempty"`
	NombreCurso  string    `yaml:"curso,omitempty"`
	Anio         string    `yaml:"anio,omitempty"`
}
