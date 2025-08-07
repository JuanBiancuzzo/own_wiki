package tablas

import (
	d "own_wiki/system_protocol/dependencias"
)

type Tablas struct {
	Archivos Archivos
	Tags     Tags
}

func NewTablas(tracker *d.TrackerDependencias) (*Tablas, error) {
	if archivos, err := NewArchivos(tracker); err != nil {
		return nil, err

	} else if tags, err := NewTags(&archivos, tracker); err != nil {
		return nil, err

	} else {
		return &Tablas{
			Archivos: archivos,
			Tags:     tags,
		}, nil
	}
}
