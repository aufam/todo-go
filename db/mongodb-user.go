package db

import (
	"context"
	"fmt"
	"todo-go/core"
	"todo-go/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (d *MongoDB) GetUserID(ctx context.Context, user models.UserForm) (string, error) {
	col := d.DB.Collection("users")

	var userDB models.User
	err := col.FindOne(ctx, bson.M{"username": user.Username}).Decode(&userDB)
	if err != nil {
		return "", err
	}

	if !core.CheckPassword(userDB.HashedPassword, user.Password) {
		return "", fmt.Errorf("Invalid password")
	}

	return userDB.ID, nil
}

func (d *MongoDB) GetUser(ctx context.Context, userID string) (models.UserResponse, error) {
	col := d.DB.Collection("users")

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

func (d *MongoDB) GetUsers(ctx context.Context) ([]models.UserResponse, error) {
	col := d.DB.Collection("users")

	cursor, err := col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	users := make([]models.UserResponse, 0)
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user.AsResponse())
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (d *MongoDB) AddUser(ctx context.Context, user models.UserForm) (string, error) {
	col := d.DB.Collection("users")

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

func (d *MongoDB) DelUser(ctx context.Context, userID string) error {
	col := d.DB.Collection("users")

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
