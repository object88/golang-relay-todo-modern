package data

import (
	"strconv"
)

const (
	me string = "Me"
)

// Todo is a task for a user to complete
type Todo struct {
	Complete bool   `json:"complete"`
	ID       string `json:"id"`
	Text     string `json:"text"`
}

// User represents a system user
type User struct {
	ID string `json:"id"`
}

var nextTodoID = 2
var todosByID = map[string]*Todo{
	"0": &Todo{
		Complete: true,
		ID:       "0",
		Text:     "Taste JavaScript",
	},
	"1": &Todo{
		Complete: false,
		ID:       "1",
		Text:     "Buy a unicorn",
	},
}
var todoIDsByUser = map[string][]string{
	me: []string{
		"0",
		"1",
	},
}
var usersByID = map[string]*User{
	me: &User{
		ID: me,
	},
}

// AddTodo creates a new Todo struct based on the given values, adds it to
// the collections, and returns its ID
func AddTodo(text string, complete bool) string {
	id := strconv.FormatInt(int64(nextTodoID), 10)
	nextTodoID++
	todo := &Todo{
		complete,
		id,
		text,
	}
	todosByID[id] = todo
	newList := append(todoIDsByUser[me], id)
	todoIDsByUser[me] = newList
	return id
}

// ChangeTodoComplete sets the complete propery for the given Todo ID
func ChangeTodoComplete(id string, complete bool) {
	todo := GetTodo(id)
	todo.Complete = complete
}

// GetTodo returns the Todo struct matching the given ID
func GetTodo(id string) *Todo {
	return todosByID[id]
}

// GetTodos returns all the todos for "me".  If the optional `complete`
// parameter is provided, filter the results to match that value
func GetTodos(complete *bool) []*Todo {
	todos := todoIDsByUser[me]
	if complete == nil {
		results := make([]*Todo, len(todos))
		for k, v := range todos {
			results[k] = todosByID[v]
		}
		return results
	}

	results := []*Todo{}

	for _, v := range todos {
		todo := todosByID[v]
		if todo.Complete == *complete {
			results = append(results, todo)
		}
	}

	return results
}

// GetUser returns the User struct matching the given ID
func GetUser(id string) *User {
	return usersByID[id]
}

// GetViewer returns the user
func GetViewer() *User {
	return usersByID[me]
}

// MarkAllTodos sets the `complete` value of all Todos to the provided value
func MarkAllTodos(complete bool) []*Todo {
	allTodos := GetTodos(nil)
	changedTodos := []*Todo{}

	for _, v := range allTodos {
		if v.Complete != complete {
			v.Complete = complete
			changedTodos = append(changedTodos, v)
		}
	}

	return changedTodos
}

// RemoveTodo removes the todo with the given id
func RemoveTodo(id string) {
	delete(todosByID, id)

	index := -1
	for k, v := range todoIDsByUser[me] {
		if v == id {
			index = k
			break
		}
	}

	remainderTodos := append(todoIDsByUser[me][:index], todoIDsByUser[me][index+1:]...)
	todoIDsByUser[me] = remainderTodos
}

// RemoveCompletedTodos removes the todos which are complete, and returns the
// list of thier IDs
func RemoveCompletedTodos() []string {
	todos := GetTodos(nil)

	removedIDs := []string{}
	remainderTodos := []string{}
	for _, v := range todos {
		if v.Complete {
			removedIDs = append(removedIDs, v.ID)
			delete(todosByID, v.ID)
		} else {
			remainderTodos = append(remainderTodos, v.ID)
		}
	}

	todoIDsByUser[me] = remainderTodos
	return removedIDs
}

// RenameTodo changes the todo text
func RenameTodo(id, text string) {
	todo := todosByID[id]
	todo.Text = text
}
