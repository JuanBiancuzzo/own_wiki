package fs

type ContenidoMinimo struct {
	Titulo string
	Path   string
}

func NewContenidoMinimo(titulo string, path string) ContenidoMinimo {
	return ContenidoMinimo{
		Titulo: titulo,
		Path:   path,
	}
}

type Opcion struct {
	Nombre string
	Path   string
}

func NewOpcion(nombre string, path string) Opcion {
	return Opcion{
		Nombre: nombre,
		Path:   path,
	}
}

type Data struct {
	Minimo   ContenidoMinimo
	Opciones []Opcion
}

func NewData(minimo ContenidoMinimo, opciones []Opcion) Data {
	return Data{
		Minimo:   minimo,
		Opciones: opciones,
	}
}
