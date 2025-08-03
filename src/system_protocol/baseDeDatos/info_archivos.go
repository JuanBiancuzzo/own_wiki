package baseDeDatos

type InfoArchivos struct {
	MaxPath        uint32
	MaxTags        uint32
	MaxNombre      uint32
	MaxApellido    uint32
	MaxNombreLibro uint32
	MaxEditorial   uint32
	MaxUrl         uint32
}

func (info *InfoArchivos) Incrementar() {
	info.MaxPath += 10
	info.MaxTags += 10
	info.MaxNombre += 10
	info.MaxApellido += 10
	info.MaxNombreLibro += 10
	info.MaxEditorial += 10
	info.MaxUrl += 10
}
