package database

import (
	"context"
	"encoding/json"
	"go-todos/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TodoInterface contains all functions to communicate with the database
type TodoInterface interface {
	Insert(models.Todo) (models.Todo, error)
	Update(string, interface{}) (models.TodoUpdate, error)
	Delete(string) (models.TodoDelete, error)
	Get(string) (models.Todo, error)
	Search(interface{}) ([]models.Todo, error)
}

// TodoClient contains context interface, and a pointer to mongo Collection this will give access
// to all functions of mongo Collection
type TodoClient struct {
	Ctx context.Context
	Col *mongo.Collection
}

// Insert accepts to-do object and return to-do object with inserted id otherwise, return error
func (c *TodoClient) Insert(docs models.Todo) (models.Todo, error) {
	todo := models.Todo{}

	res, err := c.Col.InsertOne(c.Ctx, docs)
	if err != nil {
		return todo, err
	}
	id := res.InsertedID.(primitive.ObjectID).Hex()
	return c.Get(id)
}

// Update expects string id, and update as interface, return TodoUpdate model or error
func (c *TodoClient) Update(id string, update interface{}) (models.TodoUpdate, error) {
	result := models.TodoUpdate{
		ModifiedCount: 0,
	}
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	todo, err := c.Get(id)
	if err != nil {
		return result, err
	}
	var exist map[string]interface{}
	b, err := json.Marshal(todo)
	if err != nil {
		return result, err
	}
	json.Unmarshal(b, &exist)

	change := update.(map[string]interface{})
	for k := range change {
		if change[k] == exist[k] {
			delete(change, k)
		}
	}

	if len(change) == 0 {
		return result, nil
	}

	res, err := c.Col.UpdateOne(c.Ctx, bson.M{"_id": _id}, bson.M{"$set": change})
	if err != nil {
		return result, err
	}

	newTodo, err := c.Get(id)
	if err != nil {
		return result, err
	}

	result.ModifiedCount = res.ModifiedCount
	result.Result = newTodo
	return result, nil
}

// Delete expects string id, and it returns TodDelete models or error
func (c *TodoClient) Delete(id string) (models.TodoDelete, error) {
	result := models.TodoDelete{
		DeletedCount: 0,
	}
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	res, err := c.Col.DeleteOne(c.Ctx, bson.M{"_id": _id})
	if err != nil {
		return result, err
	}
	result.DeletedCount = res.DeletedCount
	return result, nil
}

// Get expects string id, and return To-do models or error
func (c *TodoClient) Get(id string) (models.Todo, error) {
	todo := models.Todo{}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return todo, err
	}

	err = c.Col.FindOne(c.Ctx, bson.M{"_id": _id}).Decode(&todo)
	if err != nil {
		return todo, err
	}

	return todo, nil
}

// Search accepts filter as interface, and return slice of to-do models or error
func (c *TodoClient) Search(filter interface{}) ([]models.Todo, error) {
	todos := []models.Todo{}
	// check if filter nil then set to an empty bson.M this will return entire collection
	if filter == nil {
		filter = bson.M{}
	}

	cursor, err := c.Col.Find(c.Ctx, filter)
	if err != nil {
		return todos, err
	}

	for cursor.Next(c.Ctx) {
		row := models.Todo{}
		cursor.Decode(&row)
		todos = append(todos, row)
	}

	return todos, nil
}
