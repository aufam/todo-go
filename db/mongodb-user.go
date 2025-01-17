package db

import (
	"context"
	"fmt"
	"todo-go/core"
	"todo-go/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (d MongoDB) GetUserID(user models.UserForm) (string, error) {
	col := d.DB.Collection("users")
	ctx := context.Background()

	hashedPassword, err := core.HashPassword(user.Password)
	if err != nil {
		return "", err
	}

	var userDB models.User
	err = col.FindOne(ctx, bson.M{"username": user.Username, "hashedPassword": hashedPassword}).Decode(&userDB)
	if err != nil {
		return "", err
	}

	return userDB.ID, nil
}

func (d MongoDB) GetUser(userID string) (models.UserResponse, error) {
	col := d.DB.Collection("users")
	ctx := context.Background()

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return models.UserResponse{}, err
	}

	var user models.User
	err = col.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return models.UserResponse{}, err
	}

	return user.AsResponse(), nil
}

func (d MongoDB) GetUsers() ([]models.UserResponse, error) {
	col := d.DB.Collection("users")
	ctx := context.Background()

	users := make([]models.UserResponse, 0)
	cursor, err := col.Find(ctx, bson.M{})
	if err != nil {
		return users, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return users, err
		}

		users = append(users, user.AsResponse())
	}

	return users, nil
}

func (d MongoDB) AddUser(user models.UserForm) (string, error) {
	col := d.DB.Collection("users")
	ctx := context.Background()

	newUser, err := user.CreateModel()
	if err != nil {
		return "", err
	}

	result, err := col.InsertOne(ctx, newUser)
	if err != nil {
		return "", err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("Failed to parse InsertedID")
	}

	return insertedID.Hex(), nil
}

func (d MongoDB) DelUser(userID string) error {
	col := d.DB.Collection("users")
	ctx := context.Background()

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	result, err := col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("No User with id `%s` found", userID)
	}
	return nil
}
