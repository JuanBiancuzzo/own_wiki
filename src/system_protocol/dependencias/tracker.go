package dependencias

import (
	"fmt"
	b "own_wiki/system_protocol/bass_de_datos"
	u "own_wiki/system_protocol/utilidades"
	"strings"
)

type TipoTabla byte

const (
	DEPENDIENTE_NO_DEPENDIBLE   byte = 0b00
	DEPENDIENTE_DEPENDIBLE      byte = 0b10
	INDEPENDIENTE_NO_DEPENDIBLE byte = 0b01
	INDEPENDIENTE_DEPENDIBLE    byte = 0b11
)

func EsTipoDependiente(tipo TipoTabla) bool {
	return tipo&0b01 != 0b01
}

func EsTipoDependible(tipo TipoTabla) bool {
	return tipo&0b10 == 0b10
}

type TrackerDependencias struct {
	BasesDeDatos    *b.Bdd
	RegistrarTablas map[string]TipoTabla

	GuardarAccion map[string]map[IntFK]int64
	GuardarUpdate map[string]map[IntFK]*u.Lista[*u.Par[string, int64]]
}

func NewTrackerDependencias(bdd *b.Bdd) (*TrackerDependencias, error) {
	return &TrackerDependencias{
		BasesDeDatos:    bdd,
		RegistrarTablas: make(map[string]TipoTabla),

		GuardarAccion: make(map[string]map[IntFK]int64),
		GuardarUpdate: make(map[string]map[IntFK]*u.Lista[*u.Par[string, int64]]),
	}, nil
}

func (td *TrackerDependencias) RegistrarTabla(tabla string, tipo TipoTabla) error {
	if _, ok := td.RegistrarTablas[tabla]; ok {
		return fmt.Errorf("ya existe una tabla con esos parametros registrados")

	} else {
		td.RegistrarTablas[tabla] = tipo
	}

	if EsTipoDependiente(tipo) {
		td.GuardarAccion[tabla] = make(map[IntFK]int64)
	}

	if EsTipoDependible(tipo) {
		td.GuardarUpdate[tabla] = make(map[IntFK]*u.Lista[*u.Par[string, int64]])
	}

	return nil
}

func (td *TrackerDependencias) CargarIndepentiente(tabla string, datos ...any) error {
	return td.CargarDependientes(tabla, []ForeignKey{}, datos...)
}

/*
Necesitamos cargar los datos incluso si estan incompletos, esto nos va a permitir
tener un id para guardar, y despues podemos (si es necesario) updatear el id que
sea necesarip => tambien necesitamos tener el nombre de la key que se tiene que
actualizar
*/
func (td *TrackerDependencias) CargarDependientes(tabla string, fKeys []ForeignKey, datos ...any) error {
	if tipo, ok := td.RegistrarTablas[tabla]; ok {
		return fmt.Errorf("ya existe una tabla con esos parametros registrados")

	} else {
		// Primero cargar en tablas
		var id int64 = 0
		var hashDatos IntFK = HashDatos(datos...)

		if EsTipoDependiente(tipo) {
			td.GuardarAccion[tabla][hashDatos] = id
		}

		if EsTipoDependible(tipo) {
			for _, fKey := range fKeys {

			}
			lista, ok := td.GuardarUpdate[tabla][hashDatos]
			if !ok {
				lista := u.NewLista[*u.Par[string, int64]]()
			}
			lista.Push(u.NewPar())

			td.GuardarUpdate[tabla][hashDatos] = lista


			if 
			// td.GuardarUpdate[tabla] = make(map[IntFK]*u.Lista[*u.Par[string, int64]])
		}
	}
}

func (td *TrackerDependencias) RegistrarIndependiente(tabla Tabla, datos ...any) error {
	if tipo, ok := td.RegistrarTablas[tabla]; !ok {
		return fmt.Errorf("no se registró la tabla %s", tabla)

	} else if !EsTipoIndependiente(tipo) {
		return fmt.Errorf("la tabla %s no se registró como una tabla independiente", tabla)

	} else {
		// Crear id al insertar en sql
		var idCreado int64 = 0

		if EsTipoDependible(tipo) {
		}

		return nil
	}
}

func (td *TrackerDependencias) RegistrarDependiente(tabla Tabla, fKeys []ForeignKey, datos ...any) error {
	if _, ok := td.RegistrarTablas[tabla]; !ok {
		return fmt.Errorf("no se registró la tabla %s", tabla)

	}

	return nil
}

func (t *Tabla) InsertarQuery() string {
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", t.Nombre, strings.Join(t.Keys, ", "), strings.Repeat(", ?", len(t.Keys))[2:])
}

func (td *TrackerDependencias) RegistrarTabla(nombreTabla string) error {
	if _, ok := td.RegistroTablas[nombreTabla]; ok {
		return fmt.Errorf("ya existe una tabla con esos parametros registrados")
	}
	// Hacer una funcione que use SQL para crear una tabla para tener un query mucho mas rapido
	td.RegistroTablas[nombreTabla] = 0
	return nil
}

func (td *TrackerDependencias) RegistrarDependible(tabla *Tabla, values []string) (uint, error) {
	if _, ok := td.RegistroTablas[tabla]; !ok {
		return 0, fmt.Errorf("no se registro la tabla que se esta usando para el registro")
	}

	return 0, nil
}

func (td *TrackerDependencias) InsertarSQL(tabla *Tabla, datos ...any) (int64, error) {
	return td.BasesDeDatos.Insertar(tabla.InsertarQuery(), datos)
}
