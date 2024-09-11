package MongoDB

import (
	"context"
	logger "github.com/skwizi4/lib/logs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongo(uri, databaseName, collectionName string) (*MongoDB, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	collection := client.Database(databaseName).Collection(collectionName)

	return &MongoDB{
		client:         client,
		databaseName:   databaseName,
		collectionName: collectionName,
		logger:         logger.InitLogger(),
		collection:     collection,
	}, nil
}

// Create, Read, Update, Delete todo - Написать Универсальные методы

func (m MongoDB) Create() {}
func (m MongoDB) Read()   {}
func (m MongoDB) Delete() {}
func (m MongoDB) Update() {}
