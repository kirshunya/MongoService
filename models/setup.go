package models

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	connectionString = "mongodb+srv://admin:23022005ru@users.62pdg13.mongodb.net/?retryWrites=true&w=majority&appName=users"
	dbColumn         = "users"
	mongoClient      *mongo.Client
	db               = "users"
)

func ConnectMongo() {
	clientOptions := options.Client().ApplyURI(connectionString)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		panic(err)
	}

	mongoClient = client
}
