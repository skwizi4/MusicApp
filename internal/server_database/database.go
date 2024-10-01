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
	Create(telegramId, token string) error
}

type Mongo struct {
	Client         *mongo.Client
	Logger         logger.GoLogger
	DatabaseName   string
	CollectionName string
	Collection     *mongo.Collection
}
type AuthParams struct {
	Token      string `json:"token"`
	TelegramID string `json:"telegram_id"`
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

func (s *Mongo) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.Client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func (s *Mongo) Create(telegramId, token string) error {
	AuthParam := AuthParams{
		Token:      token,
		TelegramID: telegramId,
	}
	_, err := s.Collection.InsertOne(context.Background(), AuthParam)
	if err != nil {
		return err
	}
	return nil
}
