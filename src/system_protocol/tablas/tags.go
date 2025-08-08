package tablas

import (
	"fmt"
	b "own_wiki/system_protocol/bass_de_datos"
	d "own_wiki/system_protocol/dependencias"
)

const TAGS = "tags"
const TABLA_TAGS = `CREATE TABLE tags (
  id INT AUTO_INCREMENT PRIMARY KEY,
  tag VARCHAR(%d) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  idArchivo INT,

  FOREIGN KEY (idArchivo) REFERENCES archivos(id)
)`
const INSERTAR_TAG = "INSERT INTO tags (tag) VALUES (?)"
const QUERY_TAG = "SELECT id FROM tags WHERE tag = ?"

type Tags struct {
	RefArchivos *Archivos
	Tracker     *d.TrackerDependencias
}

func NewTags(refArchivos *Archivos, tracker *d.TrackerDependencias, canalMensajes chan string) (Tags, error) {
	tags := Tags{
		RefArchivos: refArchivos,
		Tracker:     tracker,
	}

	if err := tracker.RegistrarTabla(tags, d.DEPENDIENTE_NO_DEPENDIBLE, canalMensajes); err != nil {
		return tags, fmt.Errorf("error al registrar Tags con error: %v", err)
	} else {
		return tags, nil
	}
}

func (t Tags) Nombre() string {
	return TAGS
}

func (t Tags) CargarTag(pathArchivo, tag string) error {
	hashArchivo := t.RefArchivos.GenerarHash(pathArchivo)
	return t.Tracker.InsertarDependiente(
		t, t.GenerarHash(tag),
		[]d.ForeignKey{d.NewForeignKey("idArchivo", t.RefArchivos.Nombre(), hashArchivo)},
		tag,
	)
}

func (t Tags) GenerarHash(tag string) d.IntFK {
	return t.Tracker.Hash.HasearDatos([]byte(tag))
}

func (t Tags) Query(bdd *b.Bdd, datos ...any) (int64, error) {
	return bdd.Insertar(INSERTAR_TAG, datos...)
}

func (t Tags) ObjetoExistente(bdd *b.Bdd, datos ...any) (bool, error) {
	_, err := bdd.Obtener(QUERY_TAG, datos...)
	return err == nil, nil
}

func (t Tags) CrearTablaRelajada(bdd *b.Bdd, info *b.InfoArchivos) error {
	if err := bdd.CrearTabla(fmt.Sprintf(TABLA_TAGS, info.MaxTags)); err != nil {
		return fmt.Errorf("no se pudo crear la tabla de tags, con error: %v", err)
	}
	return nil
}

func (t Tags) RestringirTabla(bdd *b.Bdd) error {
	return nil
}

func (t Tags) ObtenerDependencias() []d.Tabla {
	return []d.Tabla{*t.RefArchivos}
}
