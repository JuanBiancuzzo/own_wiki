package dependencias

import (
	"fmt"
	"text/template"

	b "own_wiki/system_protocol/bass_de_datos"
	u "own_wiki/system_protocol/utilidades"
)

type TrackerDependencias struct {
	BasesDeDatos    *b.Bdd
	RegistrarTablas map[Tabla]TipoTabla
	Hash            *Hash

	templateTablaDependibles *template.Template
	templateTablaIncompletos *template.Template

	permitirRegistro bool
	tablasProcesar   *u.Cola[Tabla]
}

func NewTrackerDependencias(bdd *b.Bdd) (*TrackerDependencias, error) {
	templateTablaDependibles := template.New("dependibles")
	templateTablaDependibles, err := templateTablaDependibles.Parse(`
		CREATE TABLE IF NOT EXISTE aux_dependibles (
			nombreTabla ENUM({{ range . }} "{{ .Nombre() }}", {{ end }}),
			hashDatos   INT,
			idDatos     INT
		);`)
	if err != nil {
		return nil, fmt.Errorf("error al crear el template para la tabla auxiliar de los dependibles")
	}

	templateTablaIncompletos := template.New("incompletos")
	templateTablaIncompletos, err = templateTablaIncompletos.Parse(`
		CREATE TABLE IF NOT EXISTE aux_incompletos (
			tablaDependiente 	ENUM({{ range . }} "{{ .Nombre() }}", {{ end }}),
			idDependiente   	INT,
			key 				VARCHAR(255),
			tablaResultado 		ENUM({{ range . }} "{{ .Nombre() }}", {{ end }}),
			hashDato   			INT
		);`)
	if err != nil {
		return nil, fmt.Errorf("error al crear el template para la tabla auxiliar de los incompletos")
	}

	return &TrackerDependencias{
		BasesDeDatos:    bdd,
		RegistrarTablas: make(map[Tabla]TipoTabla),
		Hash:            NewHash(),

		templateTablaDependibles: templateTablaDependibles,
		templateTablaIncompletos: templateTablaIncompletos,

		permitirRegistro: true,
		tablasProcesar:   u.NewCola[Tabla](),
	}, nil
}

func (td *TrackerDependencias) RegistrarTabla(tabla Tabla, tipo TipoTabla) error {
	if !td.permitirRegistro {
		return fmt.Errorf("ya no se permiten los registros, esto es porque ya se ejecuto la configuracion auxiliar")
	}

	if _, ok := td.RegistrarTablas[tabla]; ok {
		return fmt.Errorf("ya existe una tabla con nombre: '%s'", tabla.Nombre())
	}

	td.RegistrarTablas[tabla] = tipo
	if !EsTipoDependible(tipo) {
		td.tablasProcesar.Encolar(tabla)
	}

	return nil
}

func (td *TrackerDependencias) IniciarProcesoInsertarDatos(infoArchivos *b.InfoArchivos) error {
	if !td.permitirRegistro {
		return fmt.Errorf("ya se inicio el proceso de insertar datos")
	}
	td.permitirRegistro = false

	writer := u.NewInfiniteWriter()

	// Creamos las tablas de sql para guardar esa informacion
	if err := td.templateTablaDependibles.Execute(writer, td.RegistrarTablas); err != nil {
		return fmt.Errorf("no se pudo ejecutar el template para crear la tabla de dependibles")

	} else if _, err := td.BasesDeDatos.MySQL.Exec(string(writer.Items())); err != nil {
		return fmt.Errorf("no se pudo crear la tabla de dependibles")
	}

	writer.Reset()

	if err := td.templateTablaIncompletos.Execute(writer, td.RegistrarTablas); err != nil {
		return fmt.Errorf("no se pudo ejecutar el template para crear la tabla de incompletos")

	} else if _, err := td.BasesDeDatos.MySQL.Exec(string(writer.Items())); err != nil {
		return fmt.Errorf("no se pudo crear la tabla de incompletos")
	}

	// Eliminar las tablas existentes
	// TODO:

	// Creando las tablas relajadas
	var tablasOrdenadas []Tabla = []Tabla{}
	for tabla := range td.tablasProcesar.DesencolarIterativamente {
		nombreTabla := tabla.Nombre()

		for i, tablaExistente := range tablasOrdenadas {
			if tablaExistente.Nombre() == nombreTabla {
				tablasOrdenadas = append(tablasOrdenadas[:i], tablasOrdenadas[i+1:]...)
				break
			}
		}

		tablasOrdenadas = append([]Tabla{tabla}, tablasOrdenadas...)
		for _, tablaDependible := range tabla.ObtenerDependencias() {
			td.tablasProcesar.Encolar(tablaDependible)
		}
	}

	for _, tabla := range tablasOrdenadas {
		if err := tabla.CrearTablaRelajada(td.BasesDeDatos, infoArchivos); err != nil {
			return fmt.Errorf("error al crear tablas relajadas, especificamente en %s, con error: %v", tabla.Nombre(), err)
		}
	}

	return nil
}

func (td *TrackerDependencias) TerminarProcesoInsertarDatos() error {
	for tabla := range td.RegistrarTablas {
		if err := tabla.RestringirTabla(td.BasesDeDatos); err != nil {
			return err
		}
	}
	return nil
}

func (td *TrackerDependencias) InsertarIndependiente(tabla Tabla, datos ...any) error {
	if tipo, ok := td.RegistrarTablas[tabla]; ok {
		return fmt.Errorf("ya existe una tabla con esos parametros registrados")

	} else if EsTipoDependiente(tipo) {
		return fmt.Errorf("esta tabla (%s) no es independiente", tabla.Nombre())

	} else if id, err := td.BasesDeDatos.Insertar(tabla.Query(), datos...); err != nil {
		return fmt.Errorf("error al insertar elemento en la tabla %s", tabla.Nombre())

	} else if EsTipoDependible(tipo) {
		return td.procesoDependible(id)

	} else {
		return nil
	}
}

func (td *TrackerDependencias) InsertarDependiente(tabla Tabla, fkeys []ForeignKey, datos ...any) error {
	if tipo, ok := td.RegistrarTablas[tabla]; ok {
		return fmt.Errorf("ya existe una tabla con esos parametros registrados")

	} else if !EsTipoDependiente(tipo) {
		return fmt.Errorf("esta tabla (%s) no es dependiente", tabla.Nombre())

	} else if id, err := td.BasesDeDatos.Insertar(tabla.Query(), datos...); err != nil {
		return fmt.Errorf("error al insertar elemento en la tabla %s", tabla.Nombre())

	} else if err := td.procesoDependiente(); err != nil {
		return fmt.Errorf("error al verificar o actualizar el elemnto en la tabla tabla %s, con id: %d", tabla.Nombre(), id)

	} else if EsTipoDependible(tipo) {
		return td.procesoDependible(id)

	} else {
		return nil
	}
}

// TODO: estas 2 funciones
func (td *TrackerDependencias) procesoDependiente() error {
	return nil
}

func (td *TrackerDependencias) procesoDependible(id int64) error {
	return nil
}

/*
Necesitamos cargar los datos incluso si estan incompletos, esto nos va a permitir
tener un id para guardar, y despues podemos (si es necesario) updatear el id que
sea necesarip => tambien necesitamos tener el nombre de la key que se tiene que
actualizar
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

*/
