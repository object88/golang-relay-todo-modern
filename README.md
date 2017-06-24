# Golang Relay TodoMVC

This is a experiment to merge the front-end part of the [`todo-modern`](https://github.com/relayjs/relay-examples/tree/master/todo-modern) relay example with a Golang [GraphQL / Relay implementation](https://github.com/graphql-go/relay).  I wasn't sure whether the graphql-go implementation would work with relay-modern, and it does.

This project would not happen without extensive pre-existing code from [Facebook](https://facebook.github.io/relay/) and the above projects.  When I figure out what the proper attribution should be, I will add it.

# Setup

This project is set up to use [yarn](https://yarnpkg.com/en/), but NPM should work as well.

There are two setup steps: building the schema file, and building the go server:
* `yarn run update-schema`
* `yarn run build`

Once these are complete, you can start the servers.

# Running

The Go Relay server is set up to start a node process running the front end web server, so you only have to start one process:
* `yarn run start`

Once this is started, you can navigate to `http://localhost:3000`.

All debug messages from both services should log to the console.

Ending the Go process will stop the node process as well.  Note that if the node process is stopped independently, the go process will not notice, and will not restart it.  The node service must be running, because it proxies all GraphQL requests to the GraphQL service.
