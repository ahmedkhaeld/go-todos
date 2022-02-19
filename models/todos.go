package models

type Todo struct {
	UserID    int         `json:"userId" bson:"userId"`
	ID        interface{} `json:"id,omitempty" bson:"_id,omitempty"`
	Title     string      `json:"title" bson:"title"`
	Completed bool        `json:"completed" bson:"completed"`
}

// TodoUpdate models contains the number of records modified and the updated record
type TodoUpdate struct {
	ModifiedCount int64 `json:"modifiedCount"`
	Result        Todo
}

// TodoDelete models contains the number of records deleted
type TodoDelete struct {
	DeletedCount int64 `json:"deletedCount"`
}
