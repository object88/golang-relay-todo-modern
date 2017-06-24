package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/graphql-go/handler"
	"github.com/object88/golang-relay-todo-modern/data"
)

func main() {
	// Spin up a node instance.
	ctx, cancelFn := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, "babel-node", "./server.js")
	cmd.Dir = "."
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	err := cmd.Start()
	if err != nil {
		// Failed to start the node client; shutting down.
		log.Fatalf("The node process failed to start with error '%s'; shutting down", err.Error())
		return
	}

	// simplest relay-compliant graphql server HTTP handler
	h := handler.New(&handler.Config{
		Schema: &data.Schema,
		Pretty: true,
	})

	// create graphql endpoint
	http.Handle("/graphql", h)

	// serve!
	port := ":8080"
	log.Printf(`GraphQL server starting up on http://localhost%s`, port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("ListenAndServe failed, %v", err)
	}

	cancelFn()
}
