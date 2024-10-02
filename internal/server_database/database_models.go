package server_database

import (
	"context"
	"fmt"
	logger "github.com/skwizi4/lib/logs"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	Health() map[string]string
	Add(telegramId, token, refreshToken string) error
	Update(telegramId, token, refreshToken string) error
	Delete(telegramId string) error
	Get(telegramId string) (*AuthParams, error)
}

type Mongo struct {
	Client         *mongo.Client
	Logger         logger.GoLogger
	DatabaseName   string
	CollectionName string
	Collection     *mongo.Collection
}
type AuthParams struct {
	Token        string `json:"token" bson:"token"`
	TelegramID   string `json:"telegramId" bson:"telegram_Id"`
	RefreshToken string `json:"refreshToken" bson:"refresh_Token"`
}

var (
	uri            = os.Getenv("DB_URI")
	databaseName   = os.Getenv("DB_NAME")
	collectionName = os.Getenv("DB_COLLECTION_NAME")
)

func New() (*Mongo, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	collection := client.Database(databaseName).Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}
	return &Mongo{
		Client:         client,
		DatabaseName:   databaseName,
		CollectionName: collectionName,
		Logger:         logger.InitLogger(),
		Collection:     collection,
	}, nil
}
