package contenido

// Path del link, el titulo
type Link struct {
	NodoHeader
	Path  string `bson:"path"`
	Texto string `bson:"texto"`
}

func NewLink(path, texto string) Link {
	return Link{
		NodoHeader: NewTipo(MK_Link),
		Path:       path,
		Texto:      texto,
	}
}
