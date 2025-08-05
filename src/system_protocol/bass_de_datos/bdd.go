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

func (bdd *Bdd) Insertar(query string, datos ...any) (int64, error) {
	if filaAfectada, err := bdd.MySQL.Exec(query, datos...); err != nil {
		return 0, fmt.Errorf("error al insertar con query, con error: %v", err)

	} else if id, err := filaAfectada.LastInsertId(); err != nil {
		return 0, fmt.Errorf("error al obtener id from query, con error: %v", err)

	} else {
		return id, nil
	}
}
