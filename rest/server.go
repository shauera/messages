package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/shauera/messages/persistence"

	log "github.com/sirupsen/logrus"
	config "github.com/spf13/viper"
)

func setupMux(personController PersonController) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/person", personController.CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people", personController.GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/person/{id}", personController.GetPersonEndpoint).Methods("GET")
	return router
}

// StartHTTPServer - start service messages
func StartHTTPServer(ctx context.Context, ) {
	var personRepository PersonRepository

	databaseType := config.GetString("database.type")
	switch databaseType {
		case "memory":
			personRepository = persistence.NewMemoryRepository()
		case "mongo":
			personRepository = persistence.NewMongoRepository(ctx)
		default:
			log.WithField("databaseType", databaseType).Fatal("Non supported database type")
	}

	repositoryContext, cancel := context.WithTimeout(ctx, config.GetDuration("database.timeout") * time.Second)
	defer cancel()

	personController := PersonController{repository: personRepository, ctx: repositoryContext}

	router := setupMux(personController)

	bindPort := ":" + config.GetString("service.port")
	log.Fatal(http.ListenAndServe(bindPort, router))
}