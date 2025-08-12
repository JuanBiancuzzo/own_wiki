package bass_de_datos

import (
	"database/sql"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Bdd struct {
	MySQL   *sql.DB
	MongoDB *mongo.Database
}

func NewBdd(mySQL *sql.DB, mongoDB *mongo.Database) *Bdd {
	return &Bdd{
		MySQL:   mySQL,
		MongoDB: mongoDB,
	}
}

func (bdd *Bdd) CrearTabla(query string, datos ...any) error {
	_, err := bdd.MySQL.Exec(query, datos...)
	return err
}

func (bdd *Bdd) EliminarTabla(nombreTabla string) error {
	_, err := bdd.MySQL.Exec(fmt.Sprintf("DROP TABLE %s", nombreTabla))
	return err
}

func (bdd *Bdd) Existe(query string, datos ...any) (bool, error) {
	lectura := make([]any, len(datos))
	fila := bdd.MySQL.QueryRow(query, datos...)

	if err := fila.Scan(lectura...); err != nil {
		return false, nil
	}

	return true, nil
}

func (bdd *Bdd) Obtener(query string, datos ...any) (int64, error) {
	var id int64
	fila := bdd.MySQL.QueryRow(query, datos...)

	if err := fila.Scan(&id); err != nil {
		return id, fmt.Errorf("error al intentar query la bdd, con error: %v", err)
	}

	return id, nil
}

func (bdd *Bdd) Insertar(query string, datos ...any) (int64, error) {
	if filaAfectada, err := bdd.MySQL.Exec(query, datos...); err != nil {
		return 0, fmt.Errorf("error al insertar con query, con error: %v", err)

	} else if id, err := filaAfectada.LastInsertId(); err != nil {
		return 0, fmt.Errorf("error al obtener id from query, con error: %v", err)

	} else {
		return id, nil
	}
}

func (bdd *Bdd) ObtenerOInsertar(queryObtener, queryInsertar string, datos ...any) (int64, error) {
	if id, err := bdd.Obtener(queryObtener, datos...); err == nil {
		return id, nil
	}
	return bdd.Insertar(queryInsertar, datos...)
}
