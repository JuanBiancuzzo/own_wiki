package estructura

import (
	"database/sql"
	"fmt"
)

const INSERTAR_MATERIA = "INSERT INTO materias (nombre, codigo, etapa, idCarrera, idPlan, idCuatrimestre, idArchivo) VALUES (?, ?, ?, ?, ?, ?, ?)"
const QUERY_MATERIA_PATH = `SELECT res.id FROM (
	SELECT materias.id, archivos.path FROM archivos INNER JOIN materias ON archivos.id = materias.idArchivo
) AS res WHERE res.path = ?`

const QUERY_PLANES = "SELECT id FROM planesCarrera WHERE nombre = ?"
const INSERTAR_PLAN = "INSERT INTO planesCarrera (nombre) VALUES (?)"

const QUERY_CUATRIMESTRES = "SELECT id FROM cuatrimestreCarrera WHERE anio = ? AND cuatrimestre = ?"
const INSERTAR_CUATRIMESTRE = "INSERT INTO cuatrimestreCarrera (anio, cuatrimestre) VALUES (?, ?)"

const INSERTAR_CORRELATIVAS = "INSERT INTO materiasCorrelativas (idMateria, idCorrelativa) VALUES (?, ?)"

type ParteCuatrimestre string

const (
	CUATRIMESTRE_PRIMERO = "Primero"
	CUATRIMESTRE_SEGUNDO = "Segundo"
)

type Materia struct {
	PathArchivo string
	PathCarrera string
	Nombre      string
	Codigo      string
	Plan        string
	Anio        int
	Cuatri      ParteCuatrimestre
	Etapa       Etapa
}

func NewMateria(pathArchivo string, pathCarrera string, nombre string, codigo string, plan string, repCuatri string, repEtapa string) (*Materia, error) {
	if etapa, err := ObtenerEtapa(repEtapa); err != nil {
		return nil, fmt.Errorf("error al crear materia con error: %v", err)

	} else if anio, cuatri, err := ObtenerCuatrimestreParte(repCuatri); err != nil {
		return nil, fmt.Errorf("error al crear materia con error: %v", err)

	} else {
		return &Materia{
			PathArchivo: pathArchivo,
			PathCarrera: pathCarrera,
			Nombre:      nombre,
			Codigo:      codigo,
			Anio:        anio,
			Cuatri:      cuatri,
			Etapa:       etapa,
		}, nil
	}
}

func (m *Materia) Insertar(idCarrea int64, idPlan int64, idCuatrimestre int64, idArchivo int64) []any {
	return []any{
		m.Nombre,
		m.Codigo,
		m.Etapa,
		idCarrea,
		idPlan,
		idCuatrimestre,
		idArchivo,
	}
}

func ObtenerCuatrimestreParte(representacionCuatri string) (int, ParteCuatrimestre, error) {
	var anio int
	var cuatriNum int
	var cuatri ParteCuatrimestre

	if _, err := fmt.Sscanf(representacionCuatri, "%dC%d", &anio, &cuatriNum); err != nil {
		return anio, cuatri, fmt.Errorf("el tipo de anio-cuatri (%s) no es uno de los esperados", representacionCuatri)
	}

	switch cuatriNum {
	case 1:
		cuatri = CUATRIMESTRE_PRIMERO
	case 2:
		cuatri = CUATRIMESTRE_SEGUNDO
	default:
		return anio, cuatri, fmt.Errorf("el cuatri dado por %d no es posible representar", cuatriNum)
	}

	return anio, cuatri, nil
}

func (m *Materia) CargarDatos(bdd *sql.DB, canal chan string) bool {
	canal <- fmt.Sprintf("Insertar Materia: %s", m.Nombre)

	if idArchivo, existe := Obtener(
		func() *sql.Row { return bdd.QueryRow(QUERY_ARCHIVO, m.PathArchivo) },
	); !existe {
		return false

	} else if idCarrera, existe := Obtener(
		func() *sql.Row { return bdd.QueryRow(QUERY_CARRERA_PATH, m.PathCarrera) },
	); !existe {
		return false

	} else if idPlan, err := ObtenerOInsertar(
		func() *sql.Row { return bdd.QueryRow(QUERY_PLANES, m.Plan) },
		func() (sql.Result, error) { return bdd.Exec(INSERTAR_PLAN, m.Plan) },
	); err != nil {
		canal <- fmt.Sprintf("error al hacer una querry del plan %s con error: %v", m.Plan, err)

	} else if idCuatrimestre, err := ObtenerOInsertar(
		func() *sql.Row { return bdd.QueryRow(QUERY_CUATRIMESTRES, m.Anio, m.Cuatri) },
		func() (sql.Result, error) { return bdd.Exec(INSERTAR_CUATRIMESTRE, m.Anio, m.Cuatri) },
	); err != nil {
		canal <- fmt.Sprintf("error al hacer una querry del cuatri %s parte de %d con error: %v", m.Cuatri, m.Anio, err)

	} else if _, err := bdd.Exec(INSERTAR_MATERIA, m.Insertar(idCarrera, idPlan, idCuatrimestre, idArchivo)...); err != nil {
		canal <- fmt.Sprintf("error al insertar una materia, con error: %v", err)

	}

	return true
}

/*
func ExisteArchivoCarpetaPrevia(bdd *sql.DB, tabla string, pathArchivo string) (int64, bool) {

	query := fmt.Sprintf(`
		SELECT res.id FROM (
			SELECT %s.id, archivos.path FROM archivos INNER JOIN %s ON archivos.id = %s.idArchivo
		) AS res WHERE res.path LIKE "%s/%s"
	`, tabla, tabla, tabla, carpetaPrevia, "%")

	return Obtener(func() *sql.Row { return bdd.QueryRow(query) })
}

func CargarCorrelativas(bdd *sql.DB, pathMateria string, pathCorrelativa string) bool {
	query := `SELECT res.id FROM (
			SELECT materias.id, archivos.path FROM archivos INNER JOIN materias ON archivos.id = materias.idArchivo
		) AS res WHERE res.path LIKE "%s%s"`

	if idMateria, existe := Obtener(
		func() *sql.Row { return bdd.QueryRow(fmt.Sprintf(query, "%", pathMateria)) },
	); !existe {
		fmt.Printf("No se pudo encontrar materia: %s\n", pathMateria)
		return false
	} else if idCorrelativa, existe := Obtener(
		func() *sql.Row { return bdd.QueryRow(fmt.Sprintf(query, "%", pathCorrelativa)) },
	); !existe {
		fmt.Printf("No se pudo encontrar correlativa: %s\n", pathCorrelativa)
		return false
	} else if _, err := Insertar(
		func() (sql.Result, error) { return bdd.Exec(INSERTAR_CORRELATIVAS, idMateria, idCorrelativa) },
	); err != nil {
		fmt.Printf("Error al insertar materias correlativas, con error: %v\n", err)
	}

	return true
}
*/
