package tablas

import (
	"fmt"
	b "own_wiki/system_protocol/bass_de_datos"
	d "own_wiki/system_protocol/dependencias"
)

const ARCHIVOS = "archivos"
const TABLA_ARCHIVOS = `CREATE TABLE archivos (
  id INT AUTO_INCREMENT PRIMARY KEY,
  path VARCHAR(?) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci 
)`
const INSERTAR_ARCHIVO = "INSERT INTO archivos (path) VALUES (?)"

type Archivos struct {
	Tracker *d.TrackerDependencias
}

func NewArchivos(tracker *d.TrackerDependencias) (Archivos, error) {
	archivos := Archivos{
		Tracker: tracker,
	}

	if err := tracker.RegistrarTabla(archivos, d.INDEPENDIENTE_DEPENDIBLE); err != nil {
		return archivos, fmt.Errorf("error al registrar Archivos con error: %v", err)
	} else {
		return archivos, nil
	}
}

func (a Archivos) Nombre() string {
	return ARCHIVOS
}

func (a Archivos) CargarArchivo(path string) error {
	return a.Tracker.InsertarIndependiente(a, path)
}

func (a Archivos) GenerarForeignKey(path string) d.ForeignKey {
	return d.NewForeignKey(
		ARCHIVOS,
		a.Tracker.Hash.HasearDatos([]byte(path)),
	)
}

func (a Archivos) Query() string {
	return INSERTAR_ARCHIVO
}

func (a Archivos) CrearTablaRelajada(bdd *b.Bdd, info *b.InfoArchivos) error {
	if _, err := bdd.MySQL.Exec(TABLA_ARCHIVOS, info.MaxPath); err != nil {
		return fmt.Errorf("no se pudo crear la tabla de archivos, con error: %v", err)
	}
	return nil
}

func (a Archivos) RestringirTabla(bdd *b.Bdd) error {
	return nil
}

func (a Archivos) ObtenerDependencias() []d.Tabla {
	return []d.Tabla{}
}
