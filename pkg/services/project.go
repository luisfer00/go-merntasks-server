package services

import (
	"context"
	"time"

	"github.com/luisfer00/go-merntasks-server/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetProjects(creador primitive.ObjectID) ([]models.Project, error) {
	var results []models.Project
	projectCollection := models.GetProjectCollection()

	cur, err := projectCollection.Find(context.Background(), bson.D{{"creador", creador}})
	if err != nil {
		return nil, err
	}

	err = cur.All(context.Background(), &results)
	if err != nil {
		return nil, err
	}
	if results == nil {
		results = make([]models.Project, 0)
	}
	return results, nil
}

func GetProject(ProjectID primitive.ObjectID) (*models.Project, error) {
	projectCollection := models.GetProjectCollection()
	var project models.Project

	err := projectCollection.FindOne(context.Background(), bson.D{{"_id", ProjectID}}).Decode(&project)

	return &project, err
}

func InsertProject(project models.Project) (*models.Project, error) {
	timeNow := time.Now()
	projectCollection := models.GetProjectCollection()

	project.Creado = &timeNow

	result, err := projectCollection.InsertOne(context.Background(), project)
	if err != nil {
		return nil, err
	}
	project.ID = result.InsertedID.(primitive.ObjectID)

	return &project, err
}

func UpdateProject(newProjectData models.Project) (*models.Project, error) {
	projectCollection := models.GetProjectCollection()

	updateData := bson.D{}
	if newProjectData.Nombre != "" {
		updateData = append(updateData, bson.E{"nombre", newProjectData.Nombre})
	}
	if newProjectData.Creador != nil && !newProjectData.Creador.IsZero() {
		updateData = append(updateData, bson.E{"creador", *newProjectData.Creador})
	}
	
	_, err := projectCollection.UpdateByID(context.Background(), newProjectData.ID, bson.D{{"$set", updateData}})
	if err != nil {
		return nil, err
	}

	return &newProjectData, nil
}

func DeleteProject(id primitive.ObjectID) error {
	projectCollection := models.GetProjectCollection()

	result, err := projectCollection.DeleteOne(context.Background(), bson.D{{"_id", id}})
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	
	return err
}