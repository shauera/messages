// Messages Manager
//
// The purpose of this application is to provide message persistence, analysis and easy retrieval.
//
//     Schemes: http
//     Version: 0.0.1
//     Contact: shalom<shauera@gmail.com>
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/shauera/messages/application"
	rest "github.com/shauera/messages/rest"

	log "github.com/sirupsen/logrus"
	config "github.com/spf13/viper"
)

func init() {
	application.InitConfig()
	application.InitLogging()
}

func main() {
	log.Println("Starting the Messages Manager service...")

	cancellableContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	application.InitHealth(cancellableContext)

	// Start internal HTTP - serving Health, Monitoring and internal API
	go rest.StartInternalHTTPServer(cancellableContext)

	// Start HTTP server - serving external endpoints
	go rest.StartExternalHTTPServer(cancellableContext)

	// -- Wait for a SIGINT. Run cleanup when signal is received ---
	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan struct{})
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		log.Println("Received an interrupt, stopping service...")
		// Do some cleanup...

		// Stop all running routines
		cancel()

		// Wait for routines to stop
		shutdownGraceDuration := config.GetDuration("service.shutdownGraceDuration")
		log.WithField("gracePeriod", shutdownGraceDuration).Info("Waiting for routines to stop")
		time.Sleep(shutdownGraceDuration)

		// Finished cleanup...
		close(cleanupDone)
	}()
	<-cleanupDone
	// -------------------------------------------------------------
}
