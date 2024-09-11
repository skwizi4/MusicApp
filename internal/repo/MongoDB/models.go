package MongoDB

import (
	logger "github.com/skwizi4/lib/logs"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDB struct {
	client         *mongo.Client
	logger         logger.GoLogger
	databaseName   string
	collectionName string
	collection     *mongo.Collection
}
