package contenido

// Path a la imagen, y el titulo
type Imagen struct {
	NodoHeader
	Path   string `bson:"path"`
	Titulo string `bson:"titulo"` // Lo que aparece en las tooltips
}

func NewImagen(path, titulo string) Imagen {
	return Imagen{
		NodoHeader: NewTipo(MK_Image),
		Path:       path,
		Titulo:     titulo,
	}
}
