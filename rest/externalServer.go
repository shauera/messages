package rest

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/shauera/messages/persistence"

	log "github.com/sirupsen/logrus"
	config "github.com/spf13/viper"
)

// ServiceController - abstraction to be implemented by controllers serving endpoints
type ServiceController interface {
	PublishEndpoints(*mux.Router)
}

func setupMux(serviceControllers []ServiceController) *mux.Router {
	router := mux.NewRouter()

	// publish all endpoint handlers
	for _, serviceController := range serviceControllers {
		serviceController.PublishEndpoints(router)
	}

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

// StartExternalHTTPServer - start serving messages endpoints
func StartExternalHTTPServer(ctx context.Context) {
	log.Debug("Starting external HTTP server")

	var messageRepository MessageRepository
	var err error

	databaseType := config.GetString("database.type")
	switch databaseType {
	case "memory":
		messageRepository, err = persistence.NewMemoryRepository()
	case "mongo":		
		messageRepository, err = persistence.NewMongoRepository(ctx, "MongoDB-MessagesRepository")
	default:
		log.WithField("databaseType", databaseType).Fatal("Non supported database type")
	}

	//TODO - Need to make this more tollerant for the occasion that the DB is not running yet
	if err != nil {
		log.WithError(err).WithField("databaseType", databaseType).Fatal("Could not initialize database connection")
	}

	var serviceControllers []ServiceController
	serviceControllers = append(serviceControllers, NewMessageController(messageRepository))

	router := setupMux(serviceControllers)

	bindPort := ":" + config.GetString("service.port")
	log.Fatal(http.ListenAndServe(bindPort, router))
}
