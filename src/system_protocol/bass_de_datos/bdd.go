package bass_de_datos

import (
	"database/sql"

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
