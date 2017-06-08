package data

const (
	me string = "Me"
)

// Todo is a task for a user to complete
type Todo struct {
	complete bool
	id       string
	text     string
}

// User represents a system user
type User struct {
}

var nextTodoID = 0
var todosByID = map[string]*Todo{}
var todosByUserID = map[string][]*Todo{
	me: []*Todo{},
}
var usersByID = map[string]User{}

// AddTodo creates a new Todo struct based on the given values, adds it to
// the collections, and returns its ID
func AddTodo(text string, complete bool) string {
	id := string(nextTodoID)
	nextTodoID++
	todo := &Todo{
		complete,
		id,
		text,
	}
	todosByID[id] = todo
	newList := append(todosByUserID[me], todo)
	todosByUserID[me] = newList
	return id
}

// ChangeTodoComplete sets the complete propery for the given Todo ID
func ChangeTodoComplete(id string, complete bool) {
	todo := GetTodo(id)
	todo.complete = complete
}

// GetTodo returns the Todo struct matching the given ID
func GetTodo(id string) Todo {
	return *todosByID[id]
}

// GetTodos returns all the todos for "me".  If the optional `complete`
// parameter is provided, filter the results to match that value
func GetTodos(complete *bool) []*Todo {
	todos := todosByUserID[me]
	if complete == nil {
		return todos
	}

	results := []*Todo{}

	for _, v := range todos {
		if v.complete == *complete {
			results = append(results, v)
		}
	}

	return results
}

// GetUser returns the User struct matching the given ID
func GetUser(id string) User {
	return usersByID[id]
}

// MarkAllTodos sets the `complete` value of all Todos to the provided value
func MarkAllTodos(complete bool) []*Todo {
	allTodos := GetTodos(nil)
	changedTodos := []*Todo{}

	for _, v := range allTodos {
		if v.complete != complete {
			v.complete = complete
			changedTodos = append(changedTodos, v)
		}
	}

	return changedTodos
}

// RemoveTodo removes the todo with the given id
func RemoveTodo(id string) {
	todos := GetTodos(nil)

	index := -1
	for k, v := range todos {
		if v.id == id {
			index = k
			break
		}
	}

	remainderTodos := append(todos[:index], todos[index+1:]...)
	todosByUserID[me] = remainderTodos
}

// RemoveCompletedTodos removes the todos which are complete, and returns the
// list of thier IDs
func RemoveCompletedTodos() []string {
	todos := GetTodos(nil)

	removedIDs := []string{}
	remainderTodos := []*Todo{}
	for _, v := range todos {
		if v.complete {
			removedIDs = append(removedIDs, v.id)
		} else {
			remainderTodos = append(remainderTodos, v)
		}
	}

	todosByUserID[me] = remainderTodos
	return removedIDs
}

// RenameTodo changes the todo text
func RenameTodo(id, text string) {
	todo := todosByID[id]
	todo.text = text
}
