package bdd

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func EstablecerConexionNoSQL(canalMensajes chan string) (*mongo.Database, error) {
	dbHost := os.Getenv("MONGO_HOST")
	dbPort := os.Getenv("MONGO_PORT")
	dbName := os.Getenv("MONGO_NAME")

	uri := fmt.Sprintf("mongodb://%s:%s/", dbHost, dbPort)
	canalMensajes <- fmt.Sprintf("Conectando a MongoDB con: %s", uri)

	cliente, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("error connecting to mongoDB, con error: %v", err)
	}

	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if err = cliente.Ping(ctx, nil); err == nil {
			canalMensajes <- "Se conecto correctamente a la MongoDB"
			break
		} else {
			// canalMensajes <- fmt.Sprintf("Error al hacer ping a MongoDB, error: %v", err)
		}
	}

	return cliente.Database(dbName), nil
}

func CrearColecciones(bdd *mongo.Database) error {
	return nil
}

func CerrarBddNoSQL(bdd *mongo.Database) {
	if bdd == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	bdd.Client().Disconnect(ctx)
	cancel()
}
