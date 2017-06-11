package data

import (
	"github.com/graphql-go/graphql"
	"github.com/object88/relay"
	"golang.org/x/net/context"
)

var todoType *graphql.Object
var userType *graphql.Object

var todoConnection *relay.GraphQLConnectionDefinitions

var nodeDefinitions *relay.NodeDefinitions

// Schema is our published GraphQL representation of objects and mutations
var Schema graphql.Schema

func init() {
	nodeDefinitions = relay.NewNodeDefinitions(relay.NodeDefinitionsConfig{
		IDFetcher: func(id string, info graphql.ResolveInfo, ct context.Context) (interface{}, error) {
			resolvedID := relay.FromGlobalID(id)
			if resolvedID.Type == "Todo" {
				return GetTodo(resolvedID.ID), nil
			}
			if resolvedID.Type == "User" {
				return GetUser(resolvedID.ID), nil
			}
			return nil, nil
		},
		TypeResolve: func(p graphql.ResolveTypeParams) *graphql.Object {
			switch p.Value.(type) {
			case *Todo:
				return todoType
			case *User:
				return userType
			}
			return nil
		},
	})

	todoType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Todo",
		Description: "A todo task",
		Fields: graphql.Fields{
			"complete": &graphql.Field{
				Description: "Indicates the completeness of the todo",
				Type:        graphql.Boolean,
			},
			"id": relay.GlobalIDField("Todo", nil),
			"text": &graphql.Field{
				Description: "The text of todo",
				Type:        graphql.String,
			},
		},
		Interfaces: []*graphql.Interface{
			nodeDefinitions.NodeInterface,
		},
	})

	todoConnection = relay.ConnectionDefinitions(relay.ConnectionConfig{
		Name:     "TodoConnection",
		NodeType: todoType,
	})

	userToTodosCollectionArgs := graphql.FieldConfigArgument{
		"status": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
	}
	for k, v := range relay.ConnectionArgs {
		userToTodosCollectionArgs[k] = v
	}

	userType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "User",
		Description: "Me",
		Fields: graphql.Fields{
			"completedCount": &graphql.Field{
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					completed := true
					allTodos := GetTodos(&completed)
					return len(allTodos), nil
				},
				Type: graphql.Int,
			},
			"id": relay.GlobalIDField("user", nil),
			"todos": &graphql.Field{
				Args:        userToTodosCollectionArgs,
				Description: "The todos for this user",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					args := relay.NewConnectionArguments(p.Args)
					dataSlice := TodosToInterfaceSlice(GetTodos(nil)...)
					return relay.ConnectionFromArray(dataSlice, args), nil
				},
				Type: todoConnection.ConnectionType,
			},
			"totalCount": &graphql.Field{
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					allTodos := GetTodos(nil)
					return len(allTodos), nil
				},
				Type: graphql.Int,
			},
		},
		Interfaces: []*graphql.Interface{
			nodeDefinitions.NodeInterface,
		},
	})

	addTodoMutation := relay.MutationWithClientMutationID(relay.MutationConfig{
		Name: "AddTodoMutation",
		InputFields: graphql.InputObjectConfigFieldMap{
			"text": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		MutateAndGetPayload: func(inputMap map[string]interface{}, info graphql.ResolveInfo, ctx context.Context) (map[string]interface{}, error) {
			newID := AddTodo(inputMap["text"].(string), false)
			return map[string]interface{}{"localTodoID": newID}, nil
		},
		OutputFields: graphql.Fields{
			"todoEdge": {
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if payload, ok := p.Source.(map[string]interface{}); ok {
						todo := GetTodo(payload["localTodoID"].(string))
						cursor := relay.CursorForObjectInConnection(TodosToInterfaceSlice(GetTodos(nil)...), todo)
						return relay.EdgeType{
							Cursor: cursor,
							Node:   todo,
						}, nil
					}
					return nil, nil
				},
				Type: todoConnection.EdgeType,
			},
			"viewer": {
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return GetViewer(), nil
				},
				Type: userType,
			},
		},
	})

	changeTodoStatusMutation := relay.MutationWithClientMutationID(relay.MutationConfig{
		Name: "ChangeTodoStatusMutation",
		InputFields: graphql.InputObjectConfigFieldMap{
			"complete": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.Boolean),
			},
			"id": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		MutateAndGetPayload: func(inputMap map[string]interface{}, info graphql.ResolveInfo, ctx context.Context) (map[string]interface{}, error) {
			resolvedID := relay.FromGlobalID(inputMap["id"].(string))
			ChangeTodoComplete(resolvedID.ID, inputMap["complete"].(bool))
			return map[string]interface{}{"id": resolvedID.ID}, nil
		},
		OutputFields: graphql.Fields{
			"todo": {
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if payload, ok := p.Source.(map[string]interface{}); ok {
						return GetTodo(payload["id"].(string)), nil
					}
					return nil, nil
				},
				Type: todoType,
			},
			"viewer": {
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return GetViewer(), nil
				},
				Type: userType,
			},
		},
	})

	markAllTodosMutation := relay.MutationWithClientMutationID(relay.MutationConfig{
		Name: "MarkAllTodos",
		InputFields: graphql.InputObjectConfigFieldMap{
			"complete": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.Boolean),
			},
		},
		MutateAndGetPayload: func(inputMap map[string]interface{}, info graphql.ResolveInfo, ctx context.Context) (map[string]interface{}, error) {
			complete := inputMap["complete"].(bool)
			changedIDs := MarkAllTodos(complete)
			return map[string]interface{}{"changedIDs": changedIDs}, nil
		},
		OutputFields: graphql.Fields{
			"changedTodos": {
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if payload, ok := p.Source.(map[string]interface{}); ok {
						changedIDs := payload["changedIDs"].([]string)
						changedTodos := make([]*Todo, len(changedIDs))
						for k, v := range changedIDs {
							changedTodos[k] = GetTodo(v)
						}
						return changedTodos, nil
					}
					return nil, nil
				},
				Type: graphql.NewList(todoType),
			},
			"viewer": {
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return GetViewer(), nil
				},
				Type: userType,
			},
		},
	})

	removeCompletedTodosMutation := relay.MutationWithClientMutationID(relay.MutationConfig{
		Name: "RemoveCompletedTodos",
		MutateAndGetPayload: func(inputMap map[string]interface{}, info graphql.ResolveInfo, ctx context.Context) (map[string]interface{}, error) {
			completedIDs := RemoveCompletedTodos()
			return map[string]interface{}{"completedIDs": completedIDs}, nil
		},
		OutputFields: graphql.Fields{
			"deletedTodoIds": {
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if payload, ok := p.Source.(map[string]interface{}); ok {
						completedIDs := payload["completedIDs"].([]string)
						return completedIDs, nil
					}
					return nil, nil
				},
				Type: graphql.NewList(graphql.String),
			},
			"viewer": {
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return GetViewer(), nil
				},
				Type: userType,
			},
		},
	})

	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"addTodo":              addTodoMutation,
			"changeTodoStatus":     changeTodoStatusMutation,
			"markAllTodos":         markAllTodosMutation,
			"removeCompletedTodos": removeCompletedTodosMutation,
		},
	})

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"node": nodeDefinitions.NodeField,
			"viewer": &graphql.Field{
				Type: userType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Use context?
					return GetViewer(), nil
				},
			},
		},
	})

	var err error
	Schema, err = graphql.NewSchema(graphql.SchemaConfig{
		Mutation: mutationType,
		Query:    queryType,
		Types:    []graphql.Type{queryType, userType},
	})
	if err != nil {
		panic(err)
	}
}

// TodosToInterfaceSlice gets an interface slice.
// See https://github.com/golang/go/wiki/InterfaceSlice
func TodosToInterfaceSlice(todos ...*Todo) []interface{} {
	var interfaceSlice = make([]interface{}, len(todos))
	for i, d := range todos {
		interfaceSlice[i] = d
	}
	return interfaceSlice
}
