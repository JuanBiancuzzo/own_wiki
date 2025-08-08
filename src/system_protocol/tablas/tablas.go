package tablas

import (
	d "own_wiki/system_protocol/dependencias"
)

type Tablas struct {
	Archivos Archivos
	Tags     Tags
	Personas Personas
}

func NewTablas(tracker *d.TrackerDependencias, canalMensajes chan string) (*Tablas, error) {
	if archivos, err := NewArchivos(tracker, canalMensajes); err != nil {
		return nil, err

	} else if personas, err := NewPersonas(tracker, canalMensajes); err != nil {
		return nil, err

	} else if tags, err := NewTags(&archivos, tracker, canalMensajes); err != nil {
		return nil, err

	} else {
		return &Tablas{
			Archivos: archivos,
			Tags:     tags,
			Personas: personas,
		}, nil
	}
}
