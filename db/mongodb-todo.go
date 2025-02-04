package db

import (
	"context"
	"fmt"
	"todo-go/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (d *MongoDB) GetTodo(ctx context.Context, userID string, todoID string) (models.TodoResponse, error) {
	col := d.DB.Collection("todos")

	id, err := primitive.ObjectIDFromHex(todoID)
	if err != nil {
		return models.TodoResponse{}, err
	}

	var todo models.Todo
	err = col.FindOne(ctx, bson.M{"_id": id, "userId": userID}).Decode(&todo)
	if err != nil {
		return models.TodoResponse{}, err
	}

	return todo.AsResponse(), nil
}

func (d *MongoDB) GetTodos(ctx context.Context, userID string) ([]models.TodoResponse, error) {
	col := d.DB.Collection("todos")

	cursor, err := col.Find(ctx, bson.M{"userId": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var todos []models.TodoResponse
	for cursor.Next(ctx) {
		var todo models.Todo
		if err := cursor.Decode(&todo); err != nil {
			return todos, err
		}

		todos = append(todos, todo.AsResponse())
	}

	return todos, nil
}

func (d *MongoDB) AddTodo(ctx context.Context, userID string, todo models.TodoForm) (todoID string, err error) {
	col := d.DB.Collection("todos")

	newTodo, err := todo.CreateModel(userID)
	if err != nil {
		return "", err
	}

	result, err := col.InsertOne(ctx, newTodo)
	if err != nil {
		return "", err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("Failed to parse InsertedID")
	}

	return insertedID.Hex(), nil
}

func (d *MongoDB) ModTodo(ctx context.Context, userID string, todoID string, todo models.TodoForm) error {
	col := d.DB.Collection("todos")

	id, err := primitive.ObjectIDFromHex(todoID)
	if err != nil {
		return err
	}

	_, err = col.UpdateOne(ctx, bson.M{"_id": id, "userId": userID}, bson.M{"$set": todo})
	return err
}

func (d *MongoDB) DelTodo(ctx context.Context, userID string, todoID string) error {
	col := d.DB.Collection("todos")

	id, err := primitive.ObjectIDFromHex(todoID)
	if err != nil {
		return err
	}

	result, err := col.DeleteOne(ctx, bson.M{"_id": id, "userId": userID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("No User with id `%s` found", userID)
	}
	return nil
}
