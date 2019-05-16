package rest

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/shauera/messages/persistence"

	log "github.com/sirupsen/logrus"
	config "github.com/spf13/viper"
)

func setupMux(messageController MessageController) *mux.Router {
	router := mux.NewRouter()
	// Messages endpoints handlers
	router.HandleFunc("/messages", messageController.CreateMessage).Methods("POST")
	router.HandleFunc("/messages", messageController.ListMessages).Methods("GET")
	router.HandleFunc("/messages/{id}", messageController.GetMessageByID).Methods("GET")
	router.HandleFunc("/messages/{id}", messageController.UpdateMessageByID).Methods("PUT")
	router.HandleFunc("/messages/{id}", messageController.DeleteMessageByID).Methods("DELETE")
	// Swagger support
	stripPrefixHandler := http.StripPrefix("/swaggerui/", http.FileServer(http.Dir("./dist/")))
	router.PathPrefix("/swaggerui/").Handler(stripPrefixHandler)
	// TODO - add middleware for JWT authorization
	return router
}

// TODO - add metrics https://opencensus.io/stats/

// StartHTTPServer - start service messages
func StartHTTPServer(ctx context.Context) {
	var messageRepository MessageRepository

	databaseType := config.GetString("database.type")
	switch databaseType {
	case "memory":
		messageRepository = persistence.NewMemoryRepository()
	case "mongo":
		messageRepository = persistence.NewMongoRepository(ctx)
	default:
		log.WithField("databaseType", databaseType).Fatal("Non supported database type")
	}

	repositoryContext, cancel := context.WithTimeout(ctx, config.GetDuration("database.timeout"))
	defer cancel()

	messageController := MessageController{repository: messageRepository, ctx: repositoryContext}

	router := setupMux(messageController)

	bindPort := ":" + config.GetString("service.port")
	log.Fatal(http.ListenAndServe(bindPort, router))
}
