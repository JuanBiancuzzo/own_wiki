package tablas

import (
	"fmt"
	b "own_wiki/system_protocol/bass_de_datos"
	d "own_wiki/system_protocol/dependencias"
)

const PERSONAS = "personas"
const TABLA_PERSONAS = `CREATE TABLE personas (
  id INT AUTO_INCREMENT PRIMARY KEY,
  nombre VARCHAR(%d) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  apellido VARCHAR(%d) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci
)`
const INSERTAR_PERSONA = "INSERT INTO personas (nombre, apellido) VALUES (?, ?)"
const QUERY_PERSONA = "SELECT nombre, apellido FROM tags WHERE nombre = ? AND apellido = ?"

type Personas struct {
	Tracker *d.TrackerDependencias
}

func NewPersonas(tracker *d.TrackerDependencias, canalMensajes chan string) (Personas, error) {
	personas := Personas{
		Tracker: tracker,
	}

	if err := tracker.RegistrarTabla(personas, d.INDEPENDIENTE_DEPENDIBLE, canalMensajes); err != nil {
		return personas, fmt.Errorf("error al registrar Personas con error: %v", err)
	} else {
		return personas, nil
	}
}

func (p Personas) Nombre() string {
	return PERSONAS
}

func (p Personas) CargarPersona(nombre, apellido string) error {
	return p.Tracker.InsertarIndependiente(p, p.GenerarHash(nombre, apellido), nombre, apellido)
}

func (p Personas) GenerarHash(nombre, apellido string) d.IntFK {
	return p.Tracker.Hash.HasearDatos(append([]byte(nombre), []byte(apellido)...))
}

func (p Personas) Query(bdd *b.Bdd, datos ...any) (int64, error) {
	return bdd.Insertar(INSERTAR_PERSONA, datos...)
}

func (p Personas) ObjetoExistente(bdd *b.Bdd, datos ...any) (bool, error) {
	return bdd.Existe(QUERY_PERSONA, datos...)
}

func (p Personas) CrearTablaRelajada(bdd *b.Bdd, info *b.InfoArchivos) error {
	if err := bdd.CrearTabla(fmt.Sprintf(TABLA_PERSONAS, info.MaxNombre, info.MaxNombre)); err != nil {
		return fmt.Errorf("no se pudo crear la tabla de Personas, con error: %v", err)
	}
	return nil
}

func (p Personas) RestringirTabla(bdd *b.Bdd) error {
	return nil
}

func (p Personas) ObtenerDependencias() []d.Tabla {
	return []d.Tabla{}
}
