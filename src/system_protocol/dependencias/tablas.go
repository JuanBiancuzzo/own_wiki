package dependencias

import "fmt"

const TABLA_ARCHIVOS = "archivos"

type Archivos struct {
	Tracker *TrackerDependencias
}

func NewArchivos(tracker *TrackerDependencias) (*Archivos, error) {
	archivos := &Archivos{
		Tracker: tracker,
	}

	if err := tracker.RegistrarTabla(archivos, INDEPENDIENTE_DEPENDIBLE); err != nil {
		return nil, fmt.Errorf("error al registrar Archivos con error: %v", err)
	} else {
		return archivos, nil
	}
}

func (a *Archivos) Nombre() string {
	return TABLA_ARCHIVOS
}

func (a *Archivos) HashearDatos(datos ...any) IntFK {
	// Implementar
	return 0
}

func (a *Archivos) CrearCargable(path string) error {

	if _, err := a.Tracker.RegistrarCargable(TABLA_ARCHIVOS, path); err != nil {
		return err
	} else if err := a.Tracker.RegistrarDependible(TABLA_ARCHIVOS, path); err != nil {
		return err
	}

	return nil
}

const TABLA_TAGS = "tags"

type Tags struct {
	Tracker    *TrackerDependencias
	DepArchivo *Archivos
}

func NewTags(tracker *TrackerDependencias, dependenciaArchivos *Archivos) (*Tags, error) {
	tags := &Tags{
		Tracker:    tracker,
		DepArchivo: dependenciaArchivos,
	}

	if err := tracker.RegistrarTabla(tags, DEPENDIENTE_NO_DEPENDIBLE); err != nil {
		return nil, fmt.Errorf("error al registrar tags con error: %v", err)
	} else {
		return tags, nil
	}
}

func (t *Tags) Nombre() string {
	return TABLA_TAGS
}

func (t *Tags) CrearCargable(path string, tag string) error {

	if id, err := t.Tracker.RegistrarCargable(tag); err != nil {
		return err

	} else {
		t.DepArchivo.EstablecerDependencia(path, id)

		return nil
	}
}
