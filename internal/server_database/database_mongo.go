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

func (s *Mongo) Add(telegramId, token, refreshToken string) error {
	AuthParam := AuthParams{
		Token:        token,
		TelegramID:   telegramId,
		RefreshToken: refreshToken,
	}
	_, err := s.Collection.InsertOne(context.Background(), AuthParam)
	if err != nil {
		return err
	}
	return nil
}
func (s *Mongo) Update(telegramId, token, refreshToken string) error {
	filter := bson.D{{Key: "telegram_Id", Value: telegramId}}
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
func (s *Mongo) Delete(telegramId string) error {
	_, err := s.Collection.DeleteOne(context.Background(), bson.D{{Key: "telegram_Id", Value: telegramId}})
	if err != nil {
		return err
	}
	return nil

}
func (s *Mongo) Get(telegramId string) (*AuthParams, error) {
	var user AuthParams
	filter := bson.D{{Key: "telegram_Id", Value: telegramId}}
	if err := s.Collection.FindOne(context.Background(), filter).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
