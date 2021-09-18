package atom

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (atom *Atom) InitMongoDB(address string) error {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(address))
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Connected to MongoDB: %s\n", address)
	atom.MongoClient = client
	return nil
}
