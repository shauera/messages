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
	// Message endpoint handlers
	router.HandleFunc("/messages", messageController.CreateMessage).Methods("POST")
	router.HandleFunc("/messages", messageController.ListMessages).Methods("GET")
	router.HandleFunc("/messages/{id}", messageController.GetMessageByID).Methods("GET")
	router.HandleFunc("/messages/{id}", messageController.UpdateMessageByID).Methods("PUT")
	router.HandleFunc("/messages/{id}", messageController.DeleteMessageByID).Methods("DELETE")
	// Swagger route
	stripPrefixHandler := http.StripPrefix("/swaggerui/", http.FileServer(http.Dir("./dist/")))
	router.PathPrefix("/swaggerui/").Handler(stripPrefixHandler)
	// TODO - add middleware for JWT authorization
	//router.Use(call JWT validation)
	//https://medium.com/@adigunhammedolalekan/build-and-deploy-a-secure-rest-api-with-go-postgresql-jwt-and-gorm-6fadf3da505b
	//https://www.sohamkamani.com/blog/golang/2019-01-01-jwt-authentication/
	return router
}

// TODO - add metrics https://opencensus.io/stats/
// TODO - add health

// StartHTTPServer - start service messages
func StartHTTPServer(ctx context.Context) {
	var messageRepository MessageRepository
	var err error

	databaseType := config.GetString("database.type")
	switch databaseType {
	case "memory":
		messageRepository, err = persistence.NewMemoryRepository()
	case "mongo":
		messageRepository, err = persistence.NewMongoRepository(ctx)
	default:
		log.WithField("databaseType", databaseType).Fatal("Non supported database type")
	}

	//TODO - Need to make this more tollerant for the ocassion that the DB is not running yet
	if err != nil {
		log.WithError(err).WithField("databaseType", databaseType).Fatal("Could not initialize database")
	}

	messageController := MessageController{repository: messageRepository}

	router := setupMux(messageController)

	bindPort := ":" + config.GetString("service.port")
	log.Fatal(http.ListenAndServe(bindPort, router))
}
