# Messages Manager
[![Build Status](https://travis-ci.org/shauera/messages.png?branch=master)](https://travis-ci.org/shauera/messages)

Demo CRUD operations using REST, Swagger, Go and MongoDB. Use this repo to draw some examples of how to put together those technologies.

## Getting started
This service is written in Go. If Go is new to you, [the documentation](https://golang.org/doc/) includes information on getting set up and learning the language.

You will need a text editor. [Visual Studio Code](https://code.visualstudio.com) is useful and has reasonable support for Go.

If you want to build this project you will need a UNIX-ish development environment including the `make`, `git` and possible `docker` for building docker images.

### Retrieving third-party dependencies

This project uses [`dep`](https://golang.github.io/dep/) to manage third-party dependencies. If you are unfamiliar with `dep`, see the [introduction](https://golang.github.io/dep/docs/introduction.html) and [installation instructions](https://golang.github.io/dep/docs/installation.html).

To download all of the third-party dependencies, run `dep ensure`.

### Running the service

To run the service locally against a MongoDB container, you can use [`docker-compose`](https://docs.docker.com/compose/) like so:
```sh
docker-compose -f message-stack.yml up -d
```
Alternatively, if your docker node is part of swarm you could use:
```sh
docker stack deploy -c message-stack.yml mongo
```
Running the executable directly can be done like so:
```sh
MESSAGES_DATABASE_TYPE=memory MESSAGES_SERVICE_PORT=8092 ./messages
```
Make sure to set the correct environment variables for the services to start. In the above example the service is set to listen to port 8092 and run against an in memory database.

#### Exposed ports
##### 8090 - Messages Manager API
The external API is intended for customer use; it includes endpoints for managing messages.
##### 8081 - Mongo Express web UI
A web UI served by the mongo-express allowing direct access to MongoDB. This can be used when experimenting and for development but is not part of the Messages Manager service.

## Configuration
A configuration file (even an empty one) must be present for the service to start. All configuration settings can be written into `config.yml`. Specific configurations values can be overridden with environment variables. Look into the `config.yml` file for specific examples.
| Variable                               | Description|
| -------------------------------------- | ------------- |
| MESSAGES_SERVICE_PORT                  | TCP port that the service will listen on
| MESSAGES_SERVICE_SHUTDOWNGRACEDURATION | Duration in which the service will wait for cleaning up e.g closing db connections
|                                        |
| MESSAGES_DATABASE_TYPE                 | Use `mongo` to work against MongoDB or `memory` to simulate a database with an in memeory structure
| MESSAGES_DATABASE_SERVER               | MongoDB - the server's socket `<host>:<ip>`
| MESSAGES_DATABASE_DBNAME               | MongoDB - the collection to work against
| MESSAGES_DATABASE_USERNAME             | MongoDB - user name
| MESSAGES_DATABASE_PASSWORD             | MongoDB - password
| MESSAGES_DATABASE_TIMEOUT              | MongoDB - timeout duration for all database operations
| 
| MESSAGES_LOGGING_LEVEL                 | Logging level: `debug`, `info`, `warning`, `error`, `fatal`


## API specification
The API specification is captured in the `dist/swagger.json` file.

When the service is running the API specification can be visually interacted with using the swagger ui available at http://<hostname>:8090/swaggerui/
<img src="/images/ReadmeSwagger.png" width="150">

## Build
Run the tests with:
```sh
make tests
```
Build the messages executable:
```sh
make build
```
Build a deployable tar file:
```sh
make publish
```
