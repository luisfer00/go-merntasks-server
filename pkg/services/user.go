package services

import (
	"context"
	"errors"
	"time"

	"github.com/luisfer00/go-merntasks-server/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrUserExist = errors.New("this user already exist in db")
)

func GetUser(email string) (*models.User, error) {
	result := models.User{}
	userCollection := models.GetUserCollection()

	err := userCollection.FindOne(context.Background(), bson.D{{"email", email}}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, err
}

func GetUserByID(ID primitive.ObjectID) (*models.User, error) {
	result := models.User{}
	userCollection := models.GetUserCollection()

	err := userCollection.FindOne(context.Background(), bson.D{{"_id", ID}}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, err
}

func InsertUser(user models.User) (*models.User, error) {
	userCollection := models.GetUserCollection()
	if user.Email == "" {
		return nil, errors.New("email is empty")
	}

	existingUser, err := GetUser(user.Email)
	if err != mongo.ErrNoDocuments && err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, ErrUserExist
	}

	user.Registro = time.Now()

	result, err := userCollection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	return &user, err
}