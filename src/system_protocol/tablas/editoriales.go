package tablas

import (
	"fmt"
	b "own_wiki/system_protocol/bass_de_datos"
	d "own_wiki/system_protocol/dependencias"
)

const EDITORIALES = "editoriales"
const TABLA_EDITORIALES = `CREATE TABLE editoriales (
  id INT AUTO_INCREMENT PRIMARY KEY,
  editorial VARCHAR(%d) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci
)`
const INSERTAR_EDITORIAL = "INSERT INTO editoriales (editorial) VALUES (?)"
const QUERY_EDITORIAL = "SELECT editorial FROM tags WHERE editorial = ?"

type Editoriales struct {
	Tracker *d.TrackerDependencias
}

func NewEditoriales(tracker *d.TrackerDependencias, canalMensajes chan string) (Editoriales, error) {
	editoriales := Editoriales{
		Tracker: tracker,
	}

	if err := tracker.RegistrarTabla(editoriales, d.INDEPENDIENTE_DEPENDIBLE, canalMensajes); err != nil {
		return editoriales, fmt.Errorf("error al registrar Editoriales con error: %v", err)
	} else {
		return editoriales, nil
	}
}

func (e Editoriales) Nombre() string {
	return EDITORIALES
}

func (e Editoriales) CargarEditorial(editorial string) error {
	return e.Tracker.InsertarIndependiente(e, e.GenerarHash(editorial), editorial)
}

func (e Editoriales) GenerarHash(editorial string) d.IntFK {
	return e.Tracker.Hash.HasearDatos([]byte(editorial))
}

func (e Editoriales) Query(bdd *b.Bdd, datos ...any) (int64, error) {
	return bdd.Insertar(INSERTAR_EDITORIAL, datos...)
}

func (e Editoriales) ObjetoExistente(bdd *b.Bdd, datos ...any) (bool, error) {
	return bdd.Existe(QUERY_EDITORIAL, datos...)
}

func (e Editoriales) CrearTablaRelajada(bdd *b.Bdd, info *b.InfoArchivos) error {
	if err := bdd.CrearTabla(fmt.Sprintf(TABLA_EDITORIALES, info.MaxNombre)); err != nil {
		return fmt.Errorf("no se pudo crear la tabla de Editoriales, con error: %v", err)
	}
	return nil
}

func (e Editoriales) RestringirTabla(bdd *b.Bdd) error {
	return nil
}

func (e Editoriales) ObtenerDependencias() []d.Tabla {
	return []d.Tabla{}
}
