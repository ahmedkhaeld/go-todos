package main

import (
	"context"
	"github.com/gorilla/mux"
	"go-todos/config"
	"go-todos/database"
	"go-todos/handlers"
	"net/http"
)

func main() {
	conf := config.GetConfig()
	// var context which is needed for mongo operation
	ctx := context.TODO()

	db := database.ConnectDB(ctx, conf.Mongo)
	collection := db.Collection(conf.Mongo.Collection)

	client := &database.TodoClient{
		Col: collection,
		Ctx: ctx,
	}

	r := mux.NewRouter()

	r.HandleFunc("/todos", handlers.SearchTodos(client)).Methods("GET")
	r.HandleFunc("/todos/{id}", handlers.GetTodo(client)).Methods("GET")
	r.HandleFunc("/todos", handlers.InsertTodo(client)).Methods("POST")
	r.HandleFunc("/todos/{id}", handlers.UpdateTodo(client)).Methods("PATCH")
	r.HandleFunc("/todos/{id}", handlers.DeleteTodo(client)).Methods("DELETE")

	http.ListenAndServe(":8080", r)
}
