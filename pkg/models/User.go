package models

import (
	"time"

	"github.com/luisfer00/go-merntasks-server/internal/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Nombre string `bson:"nombre" json:"nombre"`
	Email string `bson:"email" json:"email"`
	Password string `bson:"password,omitempty" json:"password,omitempty"`
	Registro time.Time `bson:"registro,omitempty" json:"registro,omitempty"`
}

func GetUserCollection() *mongo.Collection {
	db:= config.GetDB()
	col := db.Collection("user")

	return col
}