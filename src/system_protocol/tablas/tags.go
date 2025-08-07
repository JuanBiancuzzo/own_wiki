package tablas

import (
	"fmt"
	b "own_wiki/system_protocol/bass_de_datos"
	d "own_wiki/system_protocol/dependencias"
)

const TAGS = "tags"
const TABLA_TAGS = `CREATE TABLE tags (
  tag VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  idArchivo INT NOT NULL,

  FOREIGN KEY (idArchivo) REFERENCES archivos(id)
)`
const INSERTAR_TAG = "INSERT INTO tags (tag, idArchivo) VALUES (?, ?)"

type Tags struct {
	RefArchivos *Archivos
	Tracker     *d.TrackerDependencias
}

func NewTags(refArchivos *Archivos, tracker *d.TrackerDependencias) (Tags, error) {
	tags := Tags{
		RefArchivos: refArchivos,
		Tracker:     tracker,
	}

	if err := tracker.RegistrarTabla(tags, d.DEPENDIENTE_NO_DEPENDIBLE); err != nil {
		return tags, fmt.Errorf("error al registrar Tags con error: %v", err)
	} else {
		return tags, nil
	}
}

func (t Tags) Nombre() string {
	return TAGS
}

func (t Tags) CargarTag(pathArchivo, tag string) error {
	fKey := t.RefArchivos.GenerarForeignKey(pathArchivo)
	return t.Tracker.InsertarDependiente(t, []d.ForeignKey{fKey}, tag)
}

func (t Tags) Query() string {
	return INSERTAR_TAG
}

func (t Tags) CrearTablaRelajada(bdd *b.Bdd, info *b.InfoArchivos) error {
	if _, err := bdd.MySQL.Exec(TABLA_TAGS, info.MaxTags); err != nil {
		return fmt.Errorf("no se pudo crear la tabla de tags, con error: %v", err)
	}
	return nil
}

func (t Tags) RestringirTabla(bdd *b.Bdd) error {
	return nil
}

func (t Tags) ObtenerDependencias() []d.Tabla {
	return []d.Tabla{}
}
