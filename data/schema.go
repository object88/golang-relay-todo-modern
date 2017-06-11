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
			"id": relay.GlobalIDField("user", nil),
		},
		Interfaces: []*graphql.Interface{
			nodeDefinitions.NodeInterface,
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
		Query: queryType,
		Types: []graphql.Type{queryType, userType},
	})
	if err != nil {
		panic(err)
	}
}
