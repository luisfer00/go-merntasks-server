package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func GetClient() *mongo.Client {
	if client == nil {
		DB_USER := os.Getenv("DB_USER")
		DB_PASSWD := os.Getenv("DB_PASSWD")
		DB_HOST := os.Getenv("DB_HOST")

		if DB_USER    == "" ||
			 DB_PASSWD  == "" ||
			 DB_HOST    == "" {
			log.Fatalln("empty parameters")
		}
		uri := fmt.Sprintf("mongodb+srv://%v:%v@%v", DB_USER, DB_PASSWD, DB_HOST)
		newClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))

		if err != nil {
			log.Fatalln(err)
		}
		
		log.Println("Successfully connected to DB")
		client = newClient
	}

	return client
}

func GetDB() (*mongo.Database) {
	DB_NAME := os.Getenv("DB_NAME")
	if DB_NAME == "" {
		log.Fatalln("empty parameters")
	}
	client := GetClient()

	return client.Database(DB_NAME)
}