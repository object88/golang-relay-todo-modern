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
			"id": relay.GlobalIDField("Todo", nil),
			"complete": &graphql.Field{
				Description: "Indicates the completeness of the todo",
				Type:        graphql.Boolean,
			},
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
				Args:        relay.ConnectionArgs,
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
		OutputFields: graphql.Fields{
			"todoEdge": {
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if payload, ok := p.Source.(map[string]interface{}); ok {
						return GetTodo(payload["id"].(string)), nil
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

	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"addTodo": addTodoMutation,
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
