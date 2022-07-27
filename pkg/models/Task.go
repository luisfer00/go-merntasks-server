package models

import (
	"time"

	"github.com/luisfer00/go-merntasks-server/internal/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Task struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Nombre string `bson:"nombre" json:"nombre"`
	Estado *bool `bson:"estado" json:"estado"`
	Fecha *time.Time `bson:"fecha,omitempty" json:"fecha,omitempty"`
	Proyecto *primitive.ObjectID `bson:"proyecto,omitempty" json:"proyecto,omitempty"`
}

func GetTaskCollection() *mongo.Collection {
	db := config.GetDB()
	col := db.Collection("task")

	return col
}
