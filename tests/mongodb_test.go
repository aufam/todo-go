package main

import (
	"context"
	"testing"
	"time"
	"todo-go/core"
	"todo-go/db"
	"todo-go/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

var (
	sucipto           = models.UserForm{Username: "Sucipto", Password: "qwerty"}
	marwoto           = models.UserForm{Username: "Marwoto", Password: "12345"}
	suciptoExpectedID string
	marwotoExpectedID string
	err               error
)

func testAddUser(mt *mtest.T, ctx context.Context, database db.Database) {
	mt.AddMockResponses(mtest.CreateSuccessResponse())
	suciptoExpectedID, err = database.AddUser(ctx, sucipto)
	assert.NoError(mt, err)

	mt.AddMockResponses(mtest.CreateSuccessResponse())
	marwotoExpectedID, err = database.AddUser(ctx, marwoto)
	assert.NoError(mt, err)

	mt.Logf("suciptoExpectedID: %v\n", suciptoExpectedID)
	mt.Logf("marwotoExpectedID: %v\n", marwotoExpectedID)
	assert.NotEqual(mt, suciptoExpectedID, marwotoExpectedID)
}

func testGetUserID(mt *mtest.T, ctx context.Context, database db.Database) {
	suciptoHashedPassword, err := core.HashPassword(sucipto.Password)
	assert.NoError(mt, err)

	suciptoBSON, err := bson.Marshal(models.User{
		ID:             suciptoExpectedID,
		Username:       sucipto.Username,
		HashedPassword: suciptoHashedPassword,
		CreatedAt:      time.Now(),
	})
	assert.NoError(mt, err)

	var suciptoDB primitive.D
	err = bson.Unmarshal(suciptoBSON, &suciptoDB)
	assert.NoError(mt, err)

	mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.users", mtest.FirstBatch, suciptoDB))

	suciptoID, err := database.GetUserID(ctx, sucipto)
	assert.NoError(mt, err)
	assert.Equal(mt, suciptoID, suciptoExpectedID)
}

func TestMongoDB(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("MongoDB Mock Test", func(mt *mtest.T) {
		ctx := context.Background()
		mockDB := mt.Client.Database("test")
		database := db.MongoDB{
			DB:     mockDB,
			Client: mockDB.Client(),
		}

		testAddUser(mt, ctx, &database)
		testGetUserID(mt, ctx, &database)
	})
}
