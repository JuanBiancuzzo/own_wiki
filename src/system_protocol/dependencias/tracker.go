package dependencias

import (
	"fmt"
	"slices"
	"text/template"

	b "own_wiki/system_protocol/bass_de_datos"
	u "own_wiki/system_protocol/utilidades"
)

const QUERY_TABLA_DEPENDIENTES = "SELECT idDatos FROM aux_dependibles WHERE nombreTabla = ? AND hashDatos = ?"
const INSERTAR_TABLA_INCOMPLETOS = "INSERT INTO aux_incompletos (tablaDependiente, idDependiente, key, tablaDestino, hashDatoDestino) VALUES (?, ?, ?, ?, ?)"
const QUERY_TABLA_INCOMPLETOS = "SELECT tablaDependiente, idDependiente, key FROM aux_incompletos WHERE AND tablaDestino = ? AND hashDatoDestino = ?"

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
			tablaDestino 		ENUM({{ range . }} "{{ .Nombre() }}", {{ end }}),
			hashDatosDestino   	INT
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

func (td *TrackerDependencias) IniciarProcesoInsertarDatos(infoArchivos *b.InfoArchivos, canalMensajes chan string) error {
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

	// TODO: Eliminar las tablas existentes

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

	for _, tabla := range slices.Backward(tablasOrdenadas) {
		if err := td.BasesDeDatos.EliminarTabla(tabla.Nombre()); err != nil {
			canalMensajes <- fmt.Sprintf("error al eliminar tabla %s con error: %v", tabla.Nombre(), err)
			continue
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

	if err := td.BasesDeDatos.EliminarTabla("aux_dependibles"); err != nil {
		return fmt.Errorf("error al eliminar tabla auxiliar dependibles, con error: %v", err)

	} else if err = td.BasesDeDatos.EliminarTabla("aux_incompletos"); err != nil {
		return fmt.Errorf("error al eliminar tabla auxiliar incompletos, con error: %v", err)
	}

	return nil
}

func (td *TrackerDependencias) InsertarIndependiente(tabla Tabla, hashDatos IntFK, datos ...any) error {
	if tipo, ok := td.RegistrarTablas[tabla]; ok {
		return fmt.Errorf("ya existe una tabla con esos parametros registrados")

	} else if EsTipoDependiente(tipo) {
		return fmt.Errorf("esta tabla (%s) no es independiente", tabla.Nombre())

	} else if id, err := tabla.Query(td.BasesDeDatos, datos...); err != nil {
		return fmt.Errorf("error al insertar elemento en la tabla %s", tabla.Nombre())

	} else if EsTipoDependible(tipo) {
		return td.procesoDependible(tabla, id, hashDatos)

	} else {
		return nil
	}
}

func (td *TrackerDependencias) InsertarDependiente(tabla Tabla, hashDatos IntFK, fKeys []ForeignKey, datos ...any) error {
	if tipo, ok := td.RegistrarTablas[tabla]; ok {
		return fmt.Errorf("ya existe una tabla con esos parametros registrados")

	} else if !EsTipoDependiente(tipo) {
		return fmt.Errorf("esta tabla (%s) no es dependiente", tabla.Nombre())

	} else if id, err := tabla.Query(td.BasesDeDatos, datos...); err != nil {
		return fmt.Errorf("error al insertar elemento en la tabla %s", tabla.Nombre())

	} else if err := td.procesoDependiente(tabla, id, fKeys); err != nil {
		return fmt.Errorf("error al verificar o actualizar el elemnto en la tabla tabla %s, con id: %d", tabla.Nombre(), id)

	} else if EsTipoDependible(tipo) {
		return td.procesoDependible(tabla, id, hashDatos)

	} else {
		return nil
	}
}

func (td *TrackerDependencias) procesoDependiente(tabla Tabla, idInsertado int64, fKeys []ForeignKey) error {
	for _, fKey := range fKeys {
		// Vemos si ya fue insertado la dependencia
		if id, err := td.BasesDeDatos.Obtener(QUERY_TABLA_DEPENDIENTES, fKey.TablaDestino, fKey.HashDatosDestino); err == nil {
			// Si fueron insertados, por lo que actualizamos la tabla
			query := fmt.Sprintf("UPDATE %s SET %s = %d WHERE id = %d", tabla.Nombre(), fKey.Key, id, idInsertado)
			if _, err = td.BasesDeDatos.MySQL.Exec(query); err != nil {
				return fmt.Errorf("error al actualizar %d en tabla %s, con error %v", idInsertado, tabla.Nombre(), err)
			}

		} else {
			// Como no fue insertada, tenemos que guardar la información para que se carge correctamente la dependencia
			datos := []any{tabla.Nombre(), idInsertado, fKey.Key, fKey.TablaDestino, fKey.HashDatosDestino}
			if _, err := td.BasesDeDatos.Insertar(INSERTAR_TABLA_INCOMPLETOS, datos...); err != nil {
				return fmt.Errorf("error al insertar en la tabla auxiliar de incompletos")
			}
		}
	}

	return nil
}

func (td *TrackerDependencias) procesoDependible(tabla Tabla, idInsertado int64, hashDatos IntFK) error {
	if filas, err := td.BasesDeDatos.MySQL.Query(QUERY_TABLA_INCOMPLETOS, tabla.Nombre(), hashDatos); err != nil {
		return fmt.Errorf("error al query cuales son los elementos incompletos con tabla: %s, con error: %v", tabla.Nombre(), err)

	} else {
		defer filas.Close()

		for filas.Next() {
			var tablaDependiente string
			var idDependiente int64
			var key string

			if err = filas.Scan(&tablaDependiente, &idDependiente, &key); err != nil {
				return fmt.Errorf("error al obtener datos de una query de incompletos, con error: %v", err)
			}

			query := fmt.Sprintf("UPDATE %s SET %s = %d WHERE id = %d", tablaDependiente, key, idInsertado, idDependiente)
			if _, err = td.BasesDeDatos.MySQL.Exec(query); err != nil {
				return fmt.Errorf("error al actualizar %d en tabla %s, con error %v", idInsertado, tabla.Nombre(), err)
			}
		}

		return nil
	}
}
