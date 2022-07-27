package services

import (
	"context"
	"time"

	"github.com/luisfer00/go-merntasks-server/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetTasks(projectID primitive.ObjectID) ([]models.Task, error) {
	var results []models.Task
	taskCollection := models.GetTaskCollection()

	cur, err := taskCollection.Find(context.Background(), bson.D{{"proyecto", projectID}})
	if err != nil {
		return nil, err
	}
	err = cur.All(context.Background(), &results)
	if err != nil {
		return nil, err
	}
	if results == nil {
		results = make([]models.Task, 0)
	}
	return results, err
}

func GetTask(taskID primitive.ObjectID) (*models.Task, error) {
	taskCollection := models.GetTaskCollection()
	var task models.Task

	err := taskCollection.FindOne(context.Background(), bson.D{{"_id", taskID}}).Decode(&task)

	return &task, err
}

func InsertTask(task models.Task) (*models.Task, error) {
	taskCollection := models.GetTaskCollection()
	timeNow := time.Now()
	
	task.Fecha = &timeNow

	result, err := taskCollection.InsertOne(context.Background(), task)
	if err != nil {
		return nil, err
	}

	task.ID = result.InsertedID.(primitive.ObjectID)

	return &task, nil
}

func UpdateTask(task models.Task) (*models.Task, error) {
	taskCollection := models.GetTaskCollection()

	updateData := bson.D{}
	if task.Nombre != "" {
		updateData = append(updateData, bson.E{"nombre", task.Nombre})
	}
	if task.Estado != nil {
		updateData = append(updateData, bson.E{"estado", task.Estado})
	}

	_, err := taskCollection.UpdateByID(context.Background(), task.ID, bson.D{{"$set", updateData}})

	return &task, err
}

func DeleteTask(taskID primitive.ObjectID) error {
	taskCollection := models.GetTaskCollection()

	_, err := taskCollection.DeleteOne(context.Background(), bson.D{{"_id", taskID}})

	return err
}

func DeleteTasks(projectID primitive.ObjectID) error {
	taskCollection := models.GetTaskCollection()

	_, err := taskCollection.DeleteMany(context.Background(), bson.D{{"proyecto", projectID}})

	return err
}