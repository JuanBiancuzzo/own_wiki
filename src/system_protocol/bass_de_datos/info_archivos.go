package bass_de_datos

type InfoArchivos struct {
	MaxPath        uint32
	MaxTags        uint32
	MaxNombre      uint32
	MaxApellido    uint32
	MaxNombreLibro uint32
	MaxEditorial   uint32
	MaxUrl         uint32
}

func NewInfoArchivos() *InfoArchivos {
	return &InfoArchivos{
		MaxPath:        255,
		MaxTags:        255,
		MaxNombre:      255,
		MaxApellido:    255,
		MaxNombreLibro: 255,
		MaxEditorial:   255,
		MaxUrl:         255,
	}
}
