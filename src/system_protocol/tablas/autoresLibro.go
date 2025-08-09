package tablas

import (
	"fmt"
	b "own_wiki/system_protocol/bass_de_datos"
	d "own_wiki/system_protocol/dependencias"
)

const AUTORES_LIBRO = "autoresLibro"
const TABLA_AUTORES_LIBRO = `CREATE TABLE autoresLibro (
  id INT AUTO_INCREMENT PRIMARY KEY,
  idLibro INT,
  idPersona INT,

  FOREIGN KEY (idLibro) REFERENCES libros(id),
  FOREIGN KEY (idPersona) REFERENCES personas(id)
)`
const INSERTAR_AUTORES_LIBRO = "INSERT INTO autoresLibro (idLibro, idPersona) VALUES (0, 0);"

type AutoresLibro struct {
	RefLibro    *Libros
	RefPersonas *Personas
	Tracker     *d.TrackerDependencias
}

func NewAutoresLibro(refLibro *Libros, refPersona *Personas, tracker *d.TrackerDependencias, canalMensajes chan string) (AutoresLibro, error) {
	autoresPersona := AutoresLibro{
		RefLibro:    refLibro,
		RefPersonas: refPersona,
		Tracker:     tracker,
	}

	if err := tracker.RegistrarTabla(autoresPersona, d.DEPENDIENTE_NO_DEPENDIBLE, canalMensajes); err != nil {
		return autoresPersona, fmt.Errorf("error al registrar AutoresLibro con error: %v", err)
	} else {
		return autoresPersona, nil
	}
}

func (al AutoresLibro) Nombre() string {
	return AUTORES_LIBRO
}

func (al AutoresLibro) CargarAutorLibro(titulo, subtitulo string, anio int, edicion int, volumen int, nombreAutor, apellidoAutor string) error {
	hashLibro := al.RefLibro.GenerarHash(titulo, subtitulo, anio, edicion, volumen)
	hashPersona := al.RefPersonas.GenerarHash(nombreAutor, apellidoAutor)

	return al.Tracker.InsertarDependiente(
		al, al.GenerarHash(),
		[]d.ForeignKey{
			d.NewForeignKey("idLibro", al.RefLibro.Nombre(), hashLibro),
			d.NewForeignKey("idPersona", al.RefPersonas.Nombre(), hashPersona),
		},
	)
}

func (al AutoresLibro) GenerarHash() d.IntFK {
	return al.Tracker.Hash.HasearDatos([]byte{0})
}

func (al AutoresLibro) Query(bdd *b.Bdd, datos ...any) (int64, error) {
	return bdd.Insertar(INSERTAR_AUTORES_LIBRO, datos...)
}

func (al AutoresLibro) ObjetoExistente(bdd *b.Bdd, datos ...any) (bool, error) {
	return false, nil
}

func (al AutoresLibro) CrearTablaRelajada(bdd *b.Bdd, info *b.InfoArchivos) error {
	if err := bdd.CrearTabla(TABLA_AUTORES_LIBRO); err != nil {
		return fmt.Errorf("no se pudo crear la tabla de autoresLibro, con error: %v", err)
	}
	return nil
}

func (al AutoresLibro) RestringirTabla(bdd *b.Bdd) error {
	return nil
}

func (al AutoresLibro) ObtenerDependencias() []d.Tabla {
	return []d.Tabla{*al.RefLibro, *al.RefPersonas}
}
