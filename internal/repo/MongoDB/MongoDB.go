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
		Client:         client,
		DatabaseName:   databaseName,
		CollectionName: collectionName,
		Logger:         logger.InitLogger(),
		Collection:     collection,
	}, nil
}
