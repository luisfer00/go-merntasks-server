package models

import (
	"time"

	"github.com/luisfer00/go-merntasks-server/internal/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Project struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Nombre string `bson:"nombre" json:"nombre"`
	Creador *primitive.ObjectID `bson:"creador,omitempty" json:"creador,omitempty"`
	Creado *time.Time `bson:"creado,omitempty" json:"creado,omitempty"`
}

func GetProjectCollection() *mongo.Collection {
	db := config.GetDB()
	col := db.Collection("project")

	return col
}