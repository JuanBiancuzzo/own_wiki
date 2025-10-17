package contenido

type WikiLink struct {
	NodoHeader
	Path    string `bson:"path"`
	Texto   string `bson:"texto"`
	Mostrar bool   `bson:"mostrar"`
}

func NewWikiLink(path, texto string, mostrar bool) WikiLink {
	return WikiLink{
		NodoHeader: NewTipo(MK_Link),
		Path:       path,
		Texto:      texto,
		Mostrar:    mostrar,
	}
}
