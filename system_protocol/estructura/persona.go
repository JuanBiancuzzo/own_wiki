package estructura

const QUERY_PERSONAS = "SELECT id FROM personas WHERE nombre = ? AND apellido = ?"
const INSERTAR_PERSONA = "INSERT INTO personas (nombre, apellido) VALUES (?, ?)"

type Persona struct {
	Nombre   string
	Apellido string
}

func NewPersona(nombre string, apellido string) *Persona {
	return &Persona{
		Nombre:   nombre,
		Apellido: apellido,
	}
}

func (p *Persona) Insertar() []any {
	return []any{
		p.Nombre,
		p.Apellido,
	}
}
