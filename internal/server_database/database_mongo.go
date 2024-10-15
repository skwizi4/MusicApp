package server_database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"
)

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

func (s *Mongo) Add(UserProcess, token, refreshToken string) error {
	AuthParam := AuthParams{
		Token:        token,
		UserProcess:  UserProcess,
		RefreshToken: refreshToken,
	}
	_, err := s.Collection.InsertOne(context.Background(), AuthParam)
	if err != nil {
		return err
	}
	return nil
}
func (s *Mongo) Update(UserProcess, token, refreshToken string) error {
	filter := bson.D{{Key: "user_process", Value: UserProcess}}
	update := bson.D{
		{"$set", bson.D{
			{"token", token},
			{"refresh_Token", refreshToken},
		}},
	}
	res, err := s.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}
func (s *Mongo) Delete(UserProcess string) error {
	_, err := s.Collection.DeleteOne(context.Background(), bson.D{{Key: "user_process", Value: UserProcess}})
	if err != nil {
		return err
	}
	return nil

}
func (s *Mongo) Get(UserProcess string) (*AuthParams, error) {
	var user AuthParams
	filter := bson.D{{Key: "user_process", Value: UserProcess}}
	if err := s.Collection.FindOne(context.Background(), filter).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
