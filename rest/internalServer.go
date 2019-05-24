package rest

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/shauera/messages/application"

	log "github.com/sirupsen/logrus"
	config "github.com/spf13/viper"
)

// TODO - add metrics https://opencensus.io/stats/

// StartInternalHTTPServer - start serving health, monitoring and internal API endpoints
func StartInternalHTTPServer(ctx context.Context) {
	log.Debug("Starting internal HTTP server")

	router := mux.NewRouter()

	router.HandleFunc("/health", application.GetHealth).Methods("GET")

	bindPort := ":" + config.GetString("service.healthPort")
	log.Fatal(http.ListenAndServe(bindPort, router))
}
