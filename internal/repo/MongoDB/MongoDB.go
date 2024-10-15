package MongoDB

import (
	"MusicApp/internal/server_database"
	"context"
	"fmt"
	logger "github.com/skwizi4/lib/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func InitMongo(uri, databaseName, collectionName string) (*MongoDB, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	collection := client.Database(databaseName).Collection(collectionName)
	MongoDb := &MongoDB{
		Client:         client,
		DatabaseName:   databaseName,
		CollectionName: collectionName,
		Logger:         logger.InitLogger(),
		Collection:     collection,
	}
	if MongoDb.Health()["message"] != "It's healthy" {
		log.Fatalf("expected message to be 'It's healthy', got %s", MongoDb.Health()["message"])
	}

	return MongoDb, nil
}

func (m *MongoDB) Add(UserProcess, token, refreshToken string) error {
	AuthParam := server_database.AuthParams{
		Token:        token,
		UserProcess:  UserProcess,
		RefreshToken: refreshToken,
	}
	_, err := m.Collection.InsertOne(context.Background(), AuthParam)
	if err != nil {
		return err
	}
	return nil
}
func (m *MongoDB) Update(UserProcess, token, refreshToken string) error {
	filter := bson.D{{Key: "user_process", Value: UserProcess}}
	update := bson.D{
		{"$set", bson.D{
			{"token", token},
			{"refresh_Token", refreshToken},
		}},
	}
	res, err := m.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}
func (m *MongoDB) Delete(UserProcess string) error {
	_, err := m.Collection.DeleteOne(context.Background(), bson.D{{Key: "user_process", Value: UserProcess}})
	if err != nil {
		return err
	}
	return nil

}
func (m *MongoDB) Get(UserProcess string) (*server_database.AuthParams, error) {
	var user server_database.AuthParams
	filter := bson.D{{Key: "user_process", Value: UserProcess}}
	if err := m.Collection.FindOne(context.Background(), filter).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
func (m *MongoDB) Health() map[string]string {

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := m.Client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}
