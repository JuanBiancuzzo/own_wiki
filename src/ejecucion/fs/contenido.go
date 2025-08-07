package fs

type TextoVinculo struct {
	Titulo string
	Path   string
}

func NewTextoVinculo(titulo string, path string) TextoVinculo {
	return TextoVinculo{
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
	Minimo   TextoVinculo
	Path     []TextoVinculo
	Opciones []Opcion
}

func NewData(minimo TextoVinculo, path []TextoVinculo, opciones []Opcion) Data {
	return Data{
		Minimo:   minimo,
		Path:     path,
		Opciones: opciones,
	}
}
